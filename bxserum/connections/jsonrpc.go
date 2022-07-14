package connections

import (
	"encoding/json"
	"fmt"
)

const (
	subscribeMethod   = "subscribe"
	unsubscribeMethod = "unsubscribe"
)

// feedUpdate wraps the result from any particular stream with the subscription ID it's associated with
type feedUpdate struct {
	SubscriptionID string          `json:"subscription"`
	Result         json.RawMessage `json:"result"`
}

// subscribeParams exist because subscribe arguments usually look like ["streamName", {"some": "opts"}], which doesn't map elegantly to Go structs
type subscribeParams struct {
	StreamName string
	StreamOpts json.RawMessage
}

func (s subscribeParams) MarshalJSON() ([]byte, error) {
	nameB, err := json.Marshal(s.StreamName)
	if err != nil {
		return nil, err
	}

	params := []json.RawMessage{nameB, s.StreamOpts}
	return json.Marshal(params)
}

func (s *subscribeParams) UnmarshalJSON(b []byte) error {
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

// unsubscribeParams exist because unsubscribe arguments usually look like ["subscriptionID"], which doesn't map elegantly to Go structs
type unsubscribeParams struct {
	SubscriptionID string
}

func (s *unsubscribeParams) UnmarshalJSON(b []byte) error {
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
