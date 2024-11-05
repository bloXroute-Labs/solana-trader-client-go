package utils

import "github.com/urfave/cli/v2"

var (
	OutputFileFlag = &cli.StringFlag{
		Name:     "output",
		Usage:    "file to output CSV results to",
		Required: true,
	}
)
