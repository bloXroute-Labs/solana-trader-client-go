package stream

import (
	"context"
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/logger"
	"github.com/bloXroute-Labs/solana-trader-client-go/provider"
	"github.com/bloXroute-Labs/solana-trader-client-go/utils"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"strings"

	"time"
)

type traderWSPPumpFunNewToken struct {
	w            *provider.WSClient
	pumpTxMap    *utils.LockedMap[string, benchmark.PumpTxInfo]
	messageChan  chan *benchmark.NewTokenResult
	authHeader   string
	address      string
	rpcHost      string
	isFirstParty bool
}

func NewTraderWSPPumpFunNewToken(isFirstParty bool, messageChan chan *benchmark.NewTokenResult, pumpTxMap *utils.LockedMap[string, benchmark.PumpTxInfo],
	address, authHeader string) (Source[*benchmark.NewTokenResult, benchmark.NewTokenResult], error) {

	s := &traderWSPPumpFunNewToken{
		pumpTxMap:    pumpTxMap,
		messageChan:  messageChan,
		isFirstParty: isFirstParty,
	}

	if s.w == nil {
		w, err := provider.NewWSClientWithOpts(provider.RPCOpts{
			Endpoint:   address,
			AuthHeader: authHeader,
		})
		s.address = address
		s.authHeader = authHeader
		if err != nil {
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

				s.pumpTxMap.Update(msg.TxnHash, func(v benchmark.PumpTxInfo, exists bool) benchmark.PumpTxInfo {
					if exists {
						var firstPartyEventTime time.Time
						var thirdPartyEventTime time.Time
						if s.isFirstParty {
							// first party is late
							firstPartyEventTime = time.Now()
							thirdPartyEventTime = v.TimeSeen

						} else {
							// first party is first
							firstPartyEventTime = v.TimeSeen
							thirdPartyEventTime = time.Now()
						}

						res := &benchmark.NewTokenResult{
							TraderAPIEventTime:  firstPartyEventTime,
							ThirdPartyEventTime: thirdPartyEventTime,
							TxHash:              msg.TxnHash,
							Slot:                msg.Slot,
							Diff:                firstPartyEventTime.Sub(thirdPartyEventTime),
						}
						logger.Log().Infow("diff", "firstParty millis",
							res.Diff.Milliseconds(), "msg.TxnHash", msg.TxnHash)

						s.messageChan <- res

					} else {
						v = benchmark.PumpTxInfo{
							TimeSeen: time.Now(),
						}
					}
					return v
				})

			}()

		case <-ctx.Done():
			err = s.w.Close()
			if err != nil {
				logger.Log().Errorw("could not close connection", "err", err)
			}
			//close(s.messageChan)

			logger.Log().Infow("end of ws")
			return nil, err
		}
	}

}

func (s traderWSPPumpFunNewToken) Process(_ []RawUpdate[*benchmark.NewTokenResult], _ bool) (results map[int][]ProcessedUpdate[benchmark.NewTokenResult], duplicates map[int][]ProcessedUpdate[benchmark.NewTokenResult], err error) {
	return
}
