package csv

import (
	"encoding/csv"
	"fmt"
	"os"
)

type LinesSegment interface {
	FormatCSV() [][]string
}

func Write[T LinesSegment](outputFile string, header []string, linesSegments []T, filter func(line []string) bool) error {
	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	w := csv.NewWriter(f)
	defer w.Flush()

	err = w.Write(header)
	if err != nil {
		return err
	}

	for _, segment := range linesSegments {
		lines := segment.FormatCSV()

	LineWrite:
		for _, line := range lines {
			if filter(line) {
				continue LineWrite
			}

			if len(line) != len(header) {
				return fmt.Errorf("invalid CSV: line length (%v) differed from header (%v)", len(line), len(header))
			}

			err := w.Write(line)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
