package dao_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"otelDemo/db/dao"
	"testing"
)

func TestNewDBTXQuery(t *testing.T) {
	conn, err := dao.NewDBConn(context.Background(), "postgresql://postgres:secret@localhost:5432/tracing?sslmode=disable")
	require.NoError(t, err)
	query := dao.NewDBTXQuery(conn)
	pattern, err := query.SelectAllPattern(context.Background())
	require.NoError(t, err)
	fmt.Printf("info: %v", pattern)
}
