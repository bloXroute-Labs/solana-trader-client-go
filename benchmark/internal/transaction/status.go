package transaction

import (
	"context"
	"errors"
	"github.com/bloXroute-Labs/serum-client-go/benchmark/internal/logger"
	"github.com/bloXroute-Labs/serum-client-go/benchmark/internal/utils"
	"github.com/gagliardetto/solana-go"
	solanarpc "github.com/gagliardetto/solana-go/rpc"
	"time"
)

const (
	defaultFetchAttempts = 10
	defaultFetchInterval = 10 * time.Second
)

type StatusQuerierOpts struct {
	FetchAttempts int
	FetchInterval time.Duration
}

var defaultStatusQuerierOpts = StatusQuerierOpts{
	FetchAttempts: defaultFetchAttempts,
	FetchInterval: defaultFetchInterval,
}

type StatusQuerier struct {
	client *solanarpc.Client
	blocks map[uint64]*solanarpc.GetBlockResult
	opts   StatusQuerierOpts
}

func NewStatusQuerier(endpoint string) *StatusQuerier {
	return NewStatusQuerierWithOpts(endpoint, defaultStatusQuerierOpts)
}

func NewStatusQuerierWithOpts(endpoint string, opts StatusQuerierOpts) *StatusQuerier {
	client := solanarpc.New(endpoint)
	tsq := &StatusQuerier{
		client: client,
		blocks: make(map[uint64]*solanarpc.GetBlockResult),
		opts:   opts,
	}
	return tsq
}

func (q *StatusQuerier) FetchBatch(ctx context.Context, signatures []solana.Signature) (BatchSummary, []BlockStatus, error) {
	statuses, err := utils.AsyncGather(ctx, signatures, func(i int, ctx context.Context, signature solana.Signature) (BlockStatus, error) {
		return q.Fetch(ctx, signature)
	})

	if err != nil {
		return BatchSummary{}, nil, err
	}

	bestSlot := -1
	bestPosition := -1
	bestIndex := -1
	lostTransactions := make([]int, 0)

	for i, status := range statuses[1:] {
		if !status.Found {
			lostTransactions = append(lostTransactions, i)
			continue
		}

		replace := func() {
			bestSlot = int(status.Slot)
			bestPosition = status.Position
			bestIndex = i
		}

		// first found transaction: always best
		if bestSlot == -1 {
			replace()
			continue
		}

		// better slot: replace
		if int(status.Slot) < bestSlot {
			replace()
			continue
		}

		// same slot but better position: replace
		if int(status.Slot) == bestSlot && status.Position < bestPosition {
			replace()
			continue
		}
	}

	summary := BatchSummary{
		Best:            bestIndex,
		LostTransaction: lostTransactions,
	}
	return summary, statuses, err
}

// Fetch retrieve a transaction's slot in a block and its position within the block. This call blocks until timeout or success.
func (q *StatusQuerier) Fetch(ctx context.Context, signature solana.Signature) (BlockStatus, error) {
	ts := BlockStatus{Position: -1}

	var (
		tx  *solanarpc.GetTransactionResult
		err error
	)
	for i := 0; i < q.opts.FetchAttempts; i++ {
		tx, err = q.client.GetTransaction(ctx, signature, nil)
		if err == solanarpc.ErrNotFound {
			time.Sleep(q.opts.FetchInterval)
			continue
		}
		if err != nil {
			return ts, err
		}

		break
	}

	if tx == nil {
		logger.Log().Debugw("transaction failed execution", "signature", signature)
		return ts, nil
	}

	ts.Slot = tx.Slot
	ts.Found = true
	var (
		ok    bool
		block *solanarpc.GetBlockResult
	)
	if block, ok = q.blocks[ts.Slot]; !ok {
		opts := &solanarpc.GetBlockOpts{TransactionDetails: solanarpc.TransactionDetailsSignatures}
		block, err = q.client.GetBlockWithOpts(ctx, tx.Slot, opts)
		if err != nil {
			return ts, nil
		}

		q.blocks[ts.Slot] = block
	}

	for i, blockSignature := range block.Signatures {
		if signature == blockSignature {
			ts.Position = i
			break
		}
	}

	if ts.Position == -1 {
		return ts, errors.New("transaction signature was not found in expected block")
	}

	ts.ExecutionTime = block.BlockTime.Time()
	return ts, nil
}
