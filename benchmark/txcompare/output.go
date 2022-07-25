package main

import (
	"time"
)

type Datapoint struct {
	Iteration    int
	CreationTime time.Time
	Signature    string
	Endpoint     string

	Executed      bool
	ExecutionTime int
	Slot          uint64
	Position      int
}
