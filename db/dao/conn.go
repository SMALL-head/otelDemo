package dao

import (
	"context"
	"github.com/jackc/pgx/v5"
	"otelDemo/db/sqlc"
)

type DBTXQuery struct {
	*sqlc.Queries
	conn *pgx.Conn
}

func NewDBConn(ctx context.Context, dbHost string) (*pgx.Conn, error) {
	config, err := pgx.ParseConfig(dbHost)
	if err != nil {
		return nil, err
	}
	return pgx.ConnectConfig(ctx, config)
}

func NewDBTXQuery(conn *pgx.Conn) *DBTXQuery {
	return &DBTXQuery{
		Queries: sqlc.New(conn),
		conn:    conn,
	}
}
