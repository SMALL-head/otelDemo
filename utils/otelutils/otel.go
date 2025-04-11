package otelutils

import (
	"encoding/base64"
	"errors"
	"fmt"
)

// OtelIDToString converts a byte slice representing a span ID or trace ID to a hexadecimal string.
func OtelIDToString(spanID []byte) string {
	return fmt.Sprintf("%x", spanID)
}

func DecodeOtelID(encodedSpanID string) (string, error) {
	if encodedSpanID == "" {
		return "", errors.New("encoded span ID is empty")
	}

	decodeBytes, err := base64.StdEncoding.DecodeString(encodedSpanID)
	if err != nil {
		return "", err
	}
	return OtelIDToString(decodeBytes), nil
}
