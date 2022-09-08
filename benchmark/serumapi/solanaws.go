package main

import (
	"context"
	"github.com/bloXroute-Labs/serum-client-go/benchmark/internal/arrival"
	"github.com/bloXroute-Labs/serum-client-go/benchmark/internal/logger"
	pb "github.com/bloXroute-Labs/serum-client-go/proto"
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

type solanaRawUpdate struct {
	data *solanaws.AccountResult
	side gserum.Side
}

type solanaUpdate struct {
	side     gserum.Side
	orders   []*pb.OrderbookItem
	previous *solanaUpdate
}

func (s solanaUpdate) isRedundant() bool {
	if s.previous == nil {
		return false
	}
	return orderbookEqual(s.orders, s.previous.orders)
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

func newSolanaOrderbookStream(ctx context.Context, rpcAddress string, wsAddress, marketAddr string) (arrival.Source[solanaRawUpdate, solanaUpdate], error) {
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

// Run stops when parent ctx is canceled
func (s solanaOrderbookStream) Run(parent context.Context) ([]arrival.StreamUpdate[solanaRawUpdate], error) {
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

	messageCh := make(chan arrival.StreamUpdate[solanaRawUpdate], 200)

	// dispatch ask/bid subs
	go func() {
		for {
			if ctx.Err() != nil {
				return
			}

			ar, err := asksSub.Recv()
			if err != nil {
				s.log().Debugw("closing asks subscription", "err", err)
				cancel()
				return
			}

			messageCh <- arrival.NewStreamUpdate(solanaRawUpdate{
				data: ar,
				side: gserum.SideAsk,
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
				s.log().Debugw("closing bids subscription", "err", err)
				cancel()
				return
			}

			messageCh <- arrival.NewStreamUpdate(solanaRawUpdate{
				data: ar,
				side: gserum.SideBid,
			})
		}
	}()

	messages := make([]arrival.StreamUpdate[solanaRawUpdate], 0)
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

func (s solanaOrderbookStream) Process(updates []arrival.StreamUpdate[solanaRawUpdate], removeDuplicates bool) (map[int][]arrival.ProcessedUpdate[solanaUpdate], map[int][]arrival.ProcessedUpdate[solanaUpdate], error) {
	results := make(map[int][]arrival.ProcessedUpdate[solanaUpdate])
	duplicates := make(map[int][]arrival.ProcessedUpdate[solanaUpdate])

	previous := make(map[gserum.Side]*solanaUpdate)
	for _, update := range updates {
		var orderbook gserum.Orderbook
		err := bin.NewBinDecoder(update.Data.data.Value.Data.GetBinary()).Decode(&orderbook)
		if err != nil {
			return nil, nil, err
		}

		slot := int(update.Data.data.Context.Slot)
		_, ok := results[slot]
		if !ok {
			results[slot] = make([]arrival.ProcessedUpdate[solanaUpdate], 0)
		}

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

		side := update.Data.side
		su := solanaUpdate{
			side:     side,
			orders:   orders,
			previous: previous[side],
		}
		pu := arrival.ProcessedUpdate[solanaUpdate]{
			Timestamp: update.Timestamp,
			Slot:      slot,
			Data:      su,
		}

		if su.isRedundant() {
			duplicates[slot] = append(results[slot], pu)
		} else {
			previous[side] = &su
		}

		if !removeDuplicates || !su.isRedundant() {
			results[slot] = append(results[slot], pu)
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
