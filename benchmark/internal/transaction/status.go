package transaction

import (
	"context"
	"errors"
	"github.com/bloXroute-Labs/serum-client-go/examples/benchmark/internal/utils"
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

	return BatchSummary{}, statuses, err
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
		return ts, solanarpc.ErrNotFound
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

	return ts, nil
}
