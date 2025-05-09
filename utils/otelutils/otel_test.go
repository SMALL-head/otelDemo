package otelutils_test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"otelDemo/utils/otelutils"
	"testing"
)

func TestDecodeOtelID(t *testing.T) {
	s, err := otelutils.DecodeOtelID("3Hj505iBlX8=")
	require.NoError(t, err)
	fmt.Printf("%s", s)
}
