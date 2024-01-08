package utils

import (
	_ "embed"
)

//go:embed version.txt
var version string

func Version() string {
	return version
}
