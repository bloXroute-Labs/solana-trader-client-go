package bxassert

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const defaultChanTimeout = time.Millisecond

// ChanEqual asserts that the channel contains the expected value as the first value. ChanEqual uses a timeout to tests that end up waiting forever.
func ChanEqual[T any](t *testing.T, expected interface{}, ch chan T) {
	v := ReadChan(t, ch, defaultChanTimeout)
	assert.Equal(t, expected, v)
}

// ReadChan reads from a channel with a timeout. Timeout can be specified to avoid waiting too long on channels.
func ReadChan[T any](t *testing.T, ch chan T, timeout time.Duration) T {
	select {
	case v := <-ch:
		return v
	case <-time.After(timeout):
		assert.Fail(t, "no messages on channel")
	}
	return *new(T)
}
