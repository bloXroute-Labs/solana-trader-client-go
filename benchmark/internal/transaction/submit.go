package transaction

import (
	"context"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/logger"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/utils"
	"github.com/bloXroute-Labs/solana-trader-client-go/provider"
	"github.com/bloXroute-Labs/solana-trader-client-go/transaction"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"github.com/gagliardetto/solana-go"
	solanarpc "github.com/gagliardetto/solana-go/rpc"
	"strconv"
	"sync"
	"time"
)

const (
	defaultSubmissionInterval = 2 * time.Second
	defaultSkipPreflight      = true
)

type Builder func() (string, error)

type SubmitterOpts struct {
	SubmissionInterval time.Duration
	SkipPreflight      bool
}

var defaultSubmitterOpts = SubmitterOpts{
	SubmissionInterval: defaultSubmissionInterval,
	SkipPreflight:      defaultSkipPreflight,
}

type Submitter struct {
	clients   []*solanarpc.Client
	txBuilder Builder
	opts      SubmitterOpts
}

func NewSubmitter(endpoints []string, txBuilder Builder) *Submitter {
	return NewSubmitterWithOpts(endpoints, txBuilder, defaultSubmitterOpts)
}

func NewSubmitterWithOpts(endpoints []string, txBuilder Builder, opts SubmitterOpts) *Submitter {
	clients := make([]*solanarpc.Client, 0, len(endpoints))
	for _, endpoint := range endpoints {
		clients = append(clients, solanarpc.New(endpoint))
	}

	ts := &Submitter{
		clients:   clients,
		txBuilder: txBuilder,
		opts:      opts,
	}
	return ts
}

// SubmitIterations submits n iterations of transactions created by the builder to each of the endpoints and returns all signatures and creation times
func (ts Submitter) SubmitIterations(ctx context.Context, iterations int) ([][]solana.Signature, []time.Time, error) {
	signatures := make([][]solana.Signature, 0, iterations)
	creationTimes := make([]time.Time, 0, iterations)
	for i := 0; i < iterations; i++ {
		iterationSignatures, creationTime, err := ts.SubmitIteration(ctx)
		if err != nil {
			return nil, nil, err
		}

		creationTimes = append(creationTimes, creationTime)
		signatures = append(signatures, iterationSignatures)
		logger.Log().Debugw("submitted iteration of transactions", "iteration", i, "count", len(iterationSignatures))

		time.Sleep(ts.opts.SubmissionInterval)
	}

	return signatures, creationTimes, nil
}

// SubmitIteration uses the builder function to construct transactions for each endpoint, then sends all transactions concurrently (to be as fair as possible)
func (ts Submitter) SubmitIteration(ctx context.Context) ([]solana.Signature, time.Time, error) {
	// assume that in order transaction building is ok
	txs := make([]string, 0, len(ts.clients))
	for range ts.clients {
		tx, err := ts.txBuilder()
		if err != nil {
			return nil, time.Time{}, err
		}
		txs = append(txs, tx)
	}
	creationTime := time.Now()

	results, err := utils.AsyncGather(ctx, txs, func(i int, ctx context.Context, tx string) (solana.Signature, error) {
		return ts.submit(ctx, tx, i)
	})
	if err != nil {
		return nil, creationTime, err
	}

	for _, result := range results {
		logger.Log().Debugw("submitted transaction", "signature", result)
	}
	return results, creationTime, nil
}

func (ts Submitter) submit(ctx context.Context, txBase64 string, index int) (solana.Signature, error) {
	txBytes, err := solanarpc.DataBytesOrJSONFromBase64(txBase64)
	if err != nil {
		return solana.Signature{}, err
	}

	twm := solanarpc.TransactionWithMeta{
		Transaction: txBytes,
	}
	tx, err := twm.GetTransaction()
	if err != nil {
		return solana.Signature{}, err
	}
	opts := solanarpc.TransactionOpts{
		SkipPreflight:       ts.opts.SkipPreflight,
		PreflightCommitment: "",
	}
	signature, err := ts.clients[index].SendTransactionWithOpts(ctx, tx, opts)
	if err != nil {
		return solana.Signature{}, err
	}

	return signature, nil
}

const (
	market = "9wFFyRfZBsuAha4YcuxcXLKwMxJR43S7fPfQLusDBzvT"
)

var (
	orderID  = 1
	orderIDM = sync.Mutex{}
)

// SerumBuilder builds a transaction that's expected to fail (canceling a not found order from Serum). Transactions are submitted with `skipPreflight` however, so it should still be "executed."
func SerumBuilder(ctx context.Context, g *provider.GRPCClient, publicKey solana.PublicKey, ooAddress solana.PublicKey, privateKey solana.PrivateKey) Builder {
	return func() (string, error) {
		orderIDM.Lock()
		defer orderIDM.Unlock()

		response, err := g.PostCancelOrder(ctx, strconv.Itoa(orderID), pb.Side_S_ASK, publicKey.String(), market, ooAddress.String(), pb.Project_P_SERUM)
		if err != nil {
			return "", err
		}

		orderID++

		signedTx, err := transaction.SignTxWithPrivateKey(response.Transaction.Content, privateKey)
		if err != nil {
			return "", err
		}

		return signedTx, nil
	}
}

var (
	memoID  = 0
	memoIDM = sync.Mutex{}
)

// MemoBuilder builds a transaction with a simple memo
func MemoBuilder(privateKey solana.PrivateKey, recentBlockHashFn func() (solana.Hash, error)) Builder {
	return func() (string, error) {
		memoIDM.Lock()
		memoID++
		memoIDM.Unlock()

		publicKey := privateKey.PublicKey()

		builder := solana.NewTransactionBuilder()
		am := []*solana.AccountMeta{
			solana.Meta(publicKey).WRITE().SIGNER(),
		}

		instruction := &solana.GenericInstruction{
			AccountValues: am,
			ProgID:        solana.MemoProgramID,
			DataBytes:     []byte(strconv.Itoa(memoID)),
		}

		builder.AddInstruction(instruction)
		builder.SetFeePayer(publicKey)

		recentBlockHash, err := recentBlockHashFn()
		if err != nil {
			return "", err
		}
		builder.SetRecentBlockHash(recentBlockHash)

		tx, err := builder.Build()
		if err != nil {
			return "", err
		}

		_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
			if key == publicKey {
				return &privateKey
			}
			return nil
		})

		if err != nil {
			return "", nil
		}

		return tx.ToBase64()
	}
}
