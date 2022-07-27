package main

import (
	"fmt"
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

func Print(iterations int, endpoints []string, bests []int, lost []int) {
	fmt.Println("Iterations: ", iterations)
	fmt.Println("Endpoints:")

	for _, endpoint := range endpoints {
		fmt.Println("    ", endpoint)
	}

	fmt.Println()
	fmt.Println("Win counts: ")

	for i, endpoint := range endpoints {
		fmt.Println(fmt.Sprintf("    %-3d  %v", bests[i], endpoint))
	}

	fmt.Println()
	fmt.Println("Lost transactions: ")

	for i, endpoint := range endpoints {
		fmt.Println(fmt.Sprintf("    %-3d  %v", lost[i], endpoint))
	}
}
