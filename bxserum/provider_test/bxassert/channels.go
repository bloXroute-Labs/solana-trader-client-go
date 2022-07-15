package bxassert

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const defaultChanTimeout = time.Millisecond * 10

// ChanEqual asserts that the channel contains the expected value as the first value. ChanEqual uses a timeout to tests that end up waiting forever.
func ChanEqual[T any](t *testing.T, expected interface{}, ch chan T) {
	v := ReadChanWithTimeout(t, ch, defaultChanTimeout)
	assert.Equal(t, expected, v)
}

func ReadChan[T any](t *testing.T, ch chan T) T {
	return ReadChanWithTimeout(t, ch, defaultChanTimeout)
}

// ReadChanWithTimeout reads from a channel with a timeout. Timeout can be specified to avoid waiting too long on channels.
func ReadChanWithTimeout[T any](t *testing.T, ch chan T, timeout time.Duration) T {
	select {
	case v := <-ch:
		return v
	case <-time.After(timeout):
		assert.Fail(t, "no messages on channel")
	}
	return *new(T)
}

// ReadChanCountEquals reads from a channel and counts.When the time expires it returns the count
func ReadChanCountEquals[T any](t *testing.T, ch chan T, expected int, timeout time.Duration) {
	count := 0

Loop:
	for {
		select {
		case <-ch:
			count++
		case <-time.After(timeout):
			break Loop
		}
	}
	assert.Equal(t, expected, count)
}

func ChanEmpty[T any](t *testing.T, ch chan T) {
	select {
	case <-ch:
		assert.Fail(t, "unexpected message on channel")
	case <-time.After(defaultChanTimeout):
		return
	}
}
