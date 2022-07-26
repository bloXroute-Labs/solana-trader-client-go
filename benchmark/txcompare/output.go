package main

import (
	"strconv"
	"time"
)

const tsFormat = "15:04:05.999999"

type Datapoint struct {
	Iteration    int
	CreationTime time.Time
	Signature    string
	Endpoint     string

	Executed      bool
	ExecutionTime time.Time
	Slot          uint64
	Position      int
}

func (d Datapoint) FormatCSV() [][]string {
	return [][]string{{
		strconv.Itoa(d.Iteration),
		d.CreationTime.Format(tsFormat),
		d.Signature,
		d.Endpoint,
		strconv.FormatBool(d.Executed),
		d.ExecutionTime.Format(tsFormat),
		strconv.Itoa(int(d.Slot)),
		strconv.Itoa(d.Position),
	}}
}
