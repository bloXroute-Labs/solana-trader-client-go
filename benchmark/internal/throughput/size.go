package throughput

import "fmt"

// FormatSize turns a byte count into a human readable string.
func FormatSize(size int) string {
	if size > 1024*1024*1024 {
		return fmt.Sprintf("%v MB", size/1024/1024)
	} else if size > 1024*1024 {
		return fmt.Sprintf("%v kB", size/1024)
	}
	return fmt.Sprintf("%v B", size)
}
