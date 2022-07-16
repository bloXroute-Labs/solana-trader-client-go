package connections

import (
	"encoding/json"
	"fmt"
)

const (
	subscribeMethod   = "subscribe"
	unsubscribeMethod = "unsubscribe"
)

// FeedUpdate wraps the result from any particular stream with the subscription ID it's associated with
type FeedUpdate struct {
	SubscriptionID string          `json:"subscription"`
	Result         json.RawMessage `json:"result"`
}

// SubscribeParams exist because subscribe arguments usually look like ["streamName", {"some": "opts"}], which doesn't map elegantly to Go structs
type SubscribeParams struct {
	StreamName string
	StreamOpts json.RawMessage
}

func (s SubscribeParams) MarshalJSON() ([]byte, error) {
	nameB, err := json.Marshal(s.StreamName)
	if err != nil {
		return nil, err
	}

	params := []json.RawMessage{nameB, s.StreamOpts}
	return json.Marshal(params)
}

func (s *SubscribeParams) UnmarshalJSON(b []byte) error {
	var result []json.RawMessage
	err := json.Unmarshal(b, &result)
	if err != nil {
		return err
	}

	if len(result) != 2 {
		return fmt.Errorf("invalid argument count: expected 2, got %v", len(result))
	}

	err = json.Unmarshal(result[0], &s.StreamName)
	if err != nil {
		return err
	}

	s.StreamOpts = result[1]
	return nil
}

// UnsubscribeParams exist because unsubscribe arguments usually look like ["subscriptionID"], which doesn't map elegantly to Go structs
type UnsubscribeParams struct {
	SubscriptionID string
}

func (s UnsubscribeParams) MarshalJSON() ([]byte, error) {
	params := []string{s.SubscriptionID}
	return json.Marshal(params)
}

func (s *UnsubscribeParams) UnmarshalJSON(b []byte) error {
	var result []json.RawMessage
	err := json.Unmarshal(b, &result)
	if err != nil {
		return err
	}

	if len(result) != 1 {
		return fmt.Errorf("invalid argument count: expected 1, got %v", len(result))
	}

	err = json.Unmarshal(result[0], &s.SubscriptionID)
	if err != nil {
		return err
	}
	return nil
}
