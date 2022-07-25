package main

import (
	"context"
	"github.com/bloXroute-Labs/serum-client-go/examples/benchmark/internal/arrival"
	"github.com/bloXroute-Labs/serum-client-go/examples/benchmark/internal/logger"
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
	side   gserum.Side
	orders []*pb.OrderbookItem
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

func (s solanaOrderbookStream) Process(updates []arrival.StreamUpdate[solanaRawUpdate]) (map[int][]arrival.ProcessedUpdate[solanaUpdate], error) {
	results := make(map[int][]arrival.ProcessedUpdate[solanaUpdate])

	for _, update := range updates {

		var orderbook gserum.Orderbook
		err := bin.NewBinDecoder(update.Data.data.Value.Data.GetBinary()).Decode(&orderbook)
		if err != nil {
			return nil, err
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
			return nil, err
		}

		su := solanaUpdate{
			side:   update.Data.side,
			orders: orders,
		}

		results[slot] = append(results[slot], arrival.ProcessedUpdate[solanaUpdate]{
			Timestamp: update.Timestamp,
			Slot:      slot,
			Data:      su,
		})
	}

	return results, nil
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
