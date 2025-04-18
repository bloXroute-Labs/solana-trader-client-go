package utils

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetTimestamp returns a Protocol Buffer Timestamp object based on the current time.
func GetTimestamp() *timestamppb.Timestamp {
	now := time.Now()
	ts := timestamppb.New(now)
	return ts
}

// GetRFC3339Timestamp returns the current time as an RFC 3339 formatted string suitable for
// Protocol Buffer Timestamp JSON serialization.
func GetRFC3339Timestamp() string {
	now := time.Now().UTC()
	formatted := now.Format(time.RFC3339Nano)
	return formatted
}
