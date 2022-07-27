package throughput

import (
	"context"
	"github.com/bloXroute-Labs/serum-client-go/benchmark/internal/logger"
	"time"
)

const (
	defaultTickInterval  = 10 * time.Second
	defaultChannelBuffer = 100
)

type StreamConfig struct {
	Name          string
	TickInterval  time.Duration
	ChannelBuffer int
}

func DefaultStreamConfig(name string) StreamConfig {
	return StreamConfig{
		Name:          name,
		TickInterval:  defaultTickInterval,
		ChannelBuffer: defaultChannelBuffer,
	}
}

type StreamListener[T any] interface {
	// Connect indicates how the stream listener should connect and setup its primary stream
	Connect(context.Context) error

	// Produce indicates how to receive a series of messages on the connection
	Produce() ([]T, error)

	// Filter applies to produced items to determine whether they should be counted
	Filter(T) bool

	// Size applies to produced items to indicate its byte size
	Size(T) int

	// OnUpdate is a flexible trigger that can fire on each produced item
	OnUpdate(T)
}

func Listen[T any](parent context.Context, sl StreamListener[T], sc StreamConfig) error {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	err := sl.Connect(ctx)
	if err != nil {
		return err
	}

	logger.Log().Infow("connection made", "stream", sc.Name)

	ch := make(chan T, sc.ChannelBuffer)
	go func() {
		for {
			if ctx.Err() != nil {
				logger.Log().Infow("connection completed", "stream", sc.Name)
				return
			}

			updates, err := sl.Produce()
			if err != nil {
				cancel()
				logger.Log().Errorw("connection broken", "stream", sc.Name, "err", err)
				return
			}

			for _, update := range updates {
				if sl.Filter(update) {
					ch <- update
				}
			}
		}
	}()

	ticker := time.NewTicker(sc.TickInterval)

	count := 0
	size := 0
	startTime := time.Now()

	for {
		select {
		case msg := <-ch:
			count++
			size += sl.Size(msg)
			sl.OnUpdate(msg)
		case <-ticker.C:
			elapsedSeconds := int(time.Since(startTime).Seconds())
			throughput := size / elapsedSeconds

			logger.Log().Infow("ticker update", "stream", sc.Name, "count", count, "cps", count/elapsedSeconds, "total throughput", FormatSize(size), "throughput (/s)", FormatSize(throughput))
		}
	}
}
