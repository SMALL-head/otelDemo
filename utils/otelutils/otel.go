package otelutils

import "fmt"

// OtelIDToString converts a byte slice representing a span ID or trace ID to a hexadecimal string.
func OtelIDToString(spanID []byte) string {
	return fmt.Sprintf("%x", spanID)
}
