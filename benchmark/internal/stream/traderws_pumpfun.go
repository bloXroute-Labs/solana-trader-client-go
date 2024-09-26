package stream

import (
	"context"
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/logger"
	"github.com/bloXroute-Labs/solana-trader-client-go/provider"
	"github.com/bloXroute-Labs/solana-trader-client-go/utils"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"github.com/gagliardetto/solana-go/rpc"
	"strings"

	"time"
)

type traderWSPPumpFunNewToken struct {
	w           *provider.WSClient
	pumpTxMap   *utils.LockedMap[string, benchmark.PumpTxInfo]
	messageChan chan *benchmark.NewTokenResult
	authHeader  string
	address     string
	rpcHost     string
}

func NewTraderWSPPumpFunNewToken(messageChan chan *benchmark.NewTokenResult, pumpTxMap *utils.LockedMap[string, benchmark.PumpTxInfo],
	address, authHeader, rpcHost string) (Source[*benchmark.NewTokenResult, benchmark.NewTokenResult], error) {

	s := &traderWSPPumpFunNewToken{
		pumpTxMap:   pumpTxMap,
		messageChan: messageChan,
		rpcHost:     rpcHost,
	}

	if s.w == nil {
		w, err := provider.NewWSClientWithOpts(provider.RPCOpts{
			Endpoint:   address,
			AuthHeader: authHeader,
		})
		s.address = address
		s.authHeader = authHeader
		if err != nil {
			logger.Log().Errorw("failed to connect to trader api", "address", address, "err", err)
			return nil, err
		}
		s.w = w
	}

	return s, nil
}

func (s traderWSPPumpFunNewToken) Name() string {
	return fmt.Sprintf("traderapi")
}

// Run stops when parent ctx is canceled
func (s traderWSPPumpFunNewToken) Run(parent context.Context) ([]RawUpdate[*benchmark.NewTokenResult], error) {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	solanaRpc := rpc.New(fmt.Sprintf("https://%s", s.rpcHost))

	stream, err := s.w.GetPumpFunNewTokensStream(ctx, &pb.GetPumpFunNewTokensStreamRequest{})
	if err != nil {
		return nil, err
	}

	ch := make(chan *pb.GetPumpFunNewTokensStreamResponse, 10)
	go func() {
		for {
			v, err := stream()
			if err != nil {
				if strings.Contains(err.Error(), "shutdown requested") ||
					strings.Contains(err.Error(), "stream context has been closed") {
					return
				}
				time.Sleep(time.Second)
				logger.Log().Errorf("resetting the stream, because of error %v \n", err)
				if s.w == nil {
					w, err := provider.NewWSClientWithOpts(provider.RPCOpts{
						Endpoint:   s.address,
						AuthHeader: s.authHeader,
					})
					if err != nil {
						logger.Log().Errorw("err again", "err", err)
						continue
					} else {
						s.w = w
						stream, err = s.w.GetPumpFunNewTokensStream(ctx, &pb.GetPumpFunNewTokensStreamRequest{})
						if err != nil {
							logger.Log().Errorw("err again", "err", err)
							continue
						}
					}
				}
			} else {
				ch <- v
			}
		}
	}()
	for {
		select {
		case msg := <-ch:
			if msg == nil {
				logger.Log().Infow("receiving nil in chann")
				continue
			}

			go func() {
				tryCount := 0
				slotInfoCalled := false
				slotInfo := &rpc.GetBlockResult{}

				for {
					time.Sleep(10 * time.Second)
					tryCount++
					if tryCount >= 6 {
						logger.Log().Infow("failed to find info for tx",
							"sig", msg.TxnHash,
							"s.pumpTxMap.Len()", s.pumpTxMap.Len())
						break
					}

					if !slotInfoCalled {
						mstv := uint64(0)
						slotInfo, err = solanaRpc.GetBlockWithOpts(
							context.Background(),
							uint64(msg.Slot),
							&rpc.GetBlockOpts{
								TransactionDetails:             rpc.TransactionDetailsNone,
								Commitment:                     rpc.CommitmentConfirmed,
								MaxSupportedTransactionVersion: &mstv,
							},
						)
						if err != nil {
							logger.Log().Errorw("error occurred when getting slot info", "tryCount", tryCount)
						} else {
							slotInfoCalled = true
						}
					}

					if slotInfo != nil {
						if v, ok := s.pumpTxMap.Get(msg.TxnHash); ok {
							msg.Timestamp.AsTime().Sub(v.TimeSeen)
							logger.Log().Infow("diff", "traderAPIEventTime - rpcNodePumpTxTime = ",
								msg.Timestamp.AsTime().Sub(v.TimeSeen),
								"traderAPIEventTime - BlockTime = ",
								msg.Timestamp.AsTime().Sub(slotInfo.BlockTime.Time()),
								"traderAPIEventTime", msg.Timestamp.AsTime().UnixMilli(),
								"v.TimeSeen", v.TimeSeen.UnixMilli(),
								"BlockTime.Time()", slotInfo.BlockTime.Time())

							s.messageChan <- &benchmark.NewTokenResult{
								TraderAPIEventTime:  msg.Timestamp.AsTime(),
								ThirdPartyEventTime: v.TimeSeen,
								BlockTime:           slotInfo.BlockTime.Time(),
								Diff:                msg.Timestamp.AsTime().Sub(v.TimeSeen),
							}
							return
						}
					}
				}
			}()

		case <-ctx.Done():
			err = s.w.Close()
			if err != nil {
				logger.Log().Errorw("could not close connection", "err", err)
			}
			close(s.messageChan)
			logger.Log().Infow("end of ws")
			return nil, err
		}
	}

}

func (s traderWSPPumpFunNewToken) Process(_ []RawUpdate[*benchmark.NewTokenResult], _ bool) (results map[int][]ProcessedUpdate[benchmark.NewTokenResult], duplicates map[int][]ProcessedUpdate[benchmark.NewTokenResult], err error) {
	return
}
