package arrival

import (
	"context"
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/logger"
	pb "github.com/bloXroute-Labs/solana-trader-proto/proto/api"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	gserum "github.com/gagliardetto/solana-go/programs/serum"
	solanarpc "github.com/gagliardetto/solana-go/rpc"
	solanaws "github.com/gagliardetto/solana-go/rpc/ws"
	"go.uber.org/zap"
)

type solanaOrderbookStream struct {
	rpcClient *solanarpc.Client
	wsClient  *solanaws.Client

	wsAddress string
	marketPk  solana.PublicKey
	market    *gserum.MarketV2
	askPk     solana.PublicKey
	bidPk     solana.PublicKey
}

type SolanaRawUpdate struct {
	Data *solanaws.AccountResult
	Side gserum.Side
}

type SolanaUpdate struct {
	Side     gserum.Side
	Orders   []*pb.OrderbookItem
	previous *SolanaUpdate
}

func (s SolanaUpdate) IsRedundant() bool {
	if s.previous == nil {
		return false
	}
	return orderbookEqual(s.Orders, s.previous.Orders)
}

func orderbookEqual(o1, o2 []*pb.OrderbookItem) bool {
	if len(o1) != len(o2) {
		return false
	}

	for i, o := range o1 {
		if o.Size != o2[i].Size || o.Price != o2[i].Price {
			return false
		}
	}
	return true
}

func NewSolanaOrderbookStream(ctx context.Context, rpcAddress string, wsAddress, marketAddr string) (Source[SolanaRawUpdate, SolanaUpdate], error) {
	marketPk, err := solana.PublicKeyFromBase58(marketAddr)
	if err != nil {
		return nil, nil
	}

	s := &solanaOrderbookStream{
		rpcClient: solanarpc.New(rpcAddress),
		wsAddress: wsAddress,
		marketPk:  marketPk,
	}

	s.market, err = s.fetchMarket(ctx, marketPk)
	if err != nil {
		return nil, err
	}

	s.askPk = s.market.Asks
	s.bidPk = s.market.Bids

	s.wsClient, err = solanaws.Connect(ctx, s.wsAddress)
	if err != nil {
		return nil, err
	}

	s.log().Debugw("connection established")
	return s, nil
}

func (s solanaOrderbookStream) log() *zap.SugaredLogger {
	return logger.Log().With("source", "solanaws", "address", s.wsAddress, "market", s.marketPk.String())
}

func (s solanaOrderbookStream) Name() string {
	return fmt.Sprintf("solanaws[%v]", s.wsAddress)
}

// Run stops when parent ctx is canceled
func (s solanaOrderbookStream) Run(parent context.Context) ([]StreamUpdate[SolanaRawUpdate], error) {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	asksSub, err := s.wsClient.AccountSubscribe(s.askPk, solanarpc.CommitmentProcessed)
	if err != nil {
		return nil, err
	}

	bidsSub, err := s.wsClient.AccountSubscribe(s.bidPk, solanarpc.CommitmentProcessed)
	if err != nil {
		return nil, err
	}

	s.log().Debug("subscription created")

	messageCh := make(chan StreamUpdate[SolanaRawUpdate], 200)

	// dispatch ask/bid subs
	go func() {
		for {
			if ctx.Err() != nil {
				return
			}

			ar, err := asksSub.Recv()
			if err != nil {
				s.log().Debugw("closing Asks subscription", "err", err)
				cancel()
				return
			}

			messageCh <- NewStreamUpdate(SolanaRawUpdate{
				Data: ar,
				Side: gserum.SideAsk,
			})
		}
	}()
	go func() {
		for {
			if ctx.Err() != nil {
				return
			}

			ar, err := bidsSub.Recv()
			if err != nil {
				s.log().Debugw("closing Bids subscription", "err", err)
				cancel()
				return
			}

			messageCh <- NewStreamUpdate(SolanaRawUpdate{
				Data: ar,
				Side: gserum.SideBid,
			})
		}
	}()

	messages := make([]StreamUpdate[SolanaRawUpdate], 0)
	for {
		select {
		case msg := <-messageCh:
			messages = append(messages, msg)
		case <-ctx.Done():
			s.wsClient.Close()
			return messages, nil
		}
	}
}

func (s solanaOrderbookStream) Process(updates []StreamUpdate[SolanaRawUpdate], removeDuplicates bool) (map[int][]ProcessedUpdate[SolanaUpdate], map[int][]ProcessedUpdate[SolanaUpdate], error) {
	results := make(map[int][]ProcessedUpdate[SolanaUpdate])
	duplicates := make(map[int][]ProcessedUpdate[SolanaUpdate])

	previous := make(map[gserum.Side]*SolanaUpdate)
	for _, update := range updates {
		var orderbook gserum.Orderbook
		err := bin.NewBinDecoder(update.Data.Data.Value.Data.GetBinary()).Decode(&orderbook)
		if err != nil {
			return nil, nil, err
		}

		slot := int(update.Data.Data.Context.Slot)
		orders := make([]*pb.OrderbookItem, 0)
		err = orderbook.Items(false, func(node *gserum.SlabLeafNode) error {
			// note: price/size are not properly converted into lot sizes
			orders = append(orders, &pb.OrderbookItem{
				Price: float64(node.GetPrice().Int64()),
				Size:  float64(node.Quantity),
			})
			return nil
		})
		if err != nil {
			return nil, nil, err
		}

		side := update.Data.Side
		su := SolanaUpdate{
			Side:     side,
			Orders:   orders,
			previous: previous[side],
		}
		pu := ProcessedUpdate[SolanaUpdate]{
			Timestamp: update.Timestamp,
			Slot:      slot,
			Data:      su,
		}

		redundant := su.IsRedundant()
		if redundant {
			duplicates[slot] = append(duplicates[slot], pu)
		} else {
			previous[side] = &su
		}

		if !(removeDuplicates && redundant) {
			results[slot] = append(results[slot], pu)
			_, ok := results[slot]
			if !ok {
				results[slot] = make([]ProcessedUpdate[SolanaUpdate], 0)
			}
		}
	}

	return results, duplicates, nil
}

func (s solanaOrderbookStream) fetchMarket(ctx context.Context, marketPk solana.PublicKey) (*gserum.MarketV2, error) {
	accountInfo, err := s.rpcClient.GetAccountInfo(ctx, marketPk)
	if err != nil {
		return nil, err
	}

	var market gserum.MarketV2
	err = bin.NewBinDecoder(accountInfo.Value.Data.GetBinary()).Decode(&market)
	return &market, err
}
