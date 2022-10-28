package connections

import (
	"context"
	"encoding/json"
	"github.com/sourcegraph/jsonrpc2"
	"sync"
)

// entry to track an active subscription on connection: channel to send updates on and reference to cancel the subscription
type subscriptionEntry struct {
	active    bool
	ch        chan json.RawMessage
	onceClose sync.Once
	cancel    context.CancelFunc
}

func (s *subscriptionEntry) close() {
	s.onceClose.Do(func() {
		close(s.ch)
	})

	s.cancel()
}

type responseUpdate struct {
	v        jsonrpc2.Response
	lockHeld bool
}

type requestTracker struct {
	ch chan responseUpdate
	// can be set to hold message processing lock to ensure processing completes before next message (particularly useful for registering subscription before processing any potential updates on the connection)
	lockRequired bool
}
