package provider

import (
	"context"
	"github.com/bloXroute-Labs/serum-client-go/bxserum/connections"
	pb "github.com/bloXroute-Labs/serum-client-go/proto"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type blockHashProvider func(ctx context.Context) (*pb.GetRecentBlockHashResponse, error)
type blockHashStreamProvider func(ctx context.Context) (connections.Streamer[*pb.GetRecentBlockHashResponse], error)

type recentBlockHashStore struct {
	mutex              sync.RWMutex
	hashProvider       blockHashProvider
	hashStreamProvider blockHashStreamProvider
	hash               string
	hashTime           time.Time
	hashExpiryDuration time.Duration
}

func newRecentBlockHashStore(
	hashProvider blockHashProvider,
	streamProvider blockHashStreamProvider,
	opts RPCOpts,
) *recentBlockHashStore {
	return &recentBlockHashStore{
		mutex:              sync.RWMutex{},
		hashProvider:       hashProvider,
		hashStreamProvider: streamProvider,
		hash:               "",
		hashTime:           time.Time{},
		hashExpiryDuration: opts.CacheBlockHash,
	}
}

func (s *recentBlockHashStore) run() {
	ctx := context.Background()
	stream, err := s.hashStreamProvider(ctx)
	if err != nil {
		log.Error("can't open recent block hash stream")
		return
	}
	ch := stream.Channel(1)
	for {
		select {
		case hash := <-ch:
			s.updateBlockHash(hash)
		case <-ctx.Done():
			return
		}
	}
}

func (s *recentBlockHashStore) updateBlockHash(hash *pb.GetRecentBlockHashResponse) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.hash = hash.BlockHash
	now := time.Now()
	s.hashTime = now
}

func (s *recentBlockHashStore) recentBlockHash(ctx context.Context) (*pb.GetRecentBlockHashResponse, error) {
	response := s.getCachedBlockHash()
	if response != nil {
		return response, nil
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now()
	if s.hash != "" && s.hashTime.Before(now.Add(s.hashExpiryDuration)) {
		hash, err := s.hashProvider(ctx)
		if err != nil {
			return nil, err
		}
		s.hash = hash.BlockHash
		s.hashTime = now
	}
	return &pb.GetRecentBlockHashResponse{
		BlockHash: s.hash,
	}, nil
}

func (s *recentBlockHashStore) getCachedBlockHash() *pb.GetRecentBlockHashResponse {
	now := time.Now()
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if s.hash != "" && s.hashTime.Before(now.Add(s.hashExpiryDuration)) {
		return &pb.GetRecentBlockHashResponse{
			BlockHash: s.hash,
		}
	}
	return nil
}
