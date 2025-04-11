package svc

import (
	"context"
	"otelDemo/db/dao"
	"otelDemo/db/sqlc"
)

type Svc struct {
	dao *dao.DBTXQuery
}

func New(db *dao.DBTXQuery) *Svc {
	return &Svc{
		dao: db,
	}
}

func (s *Svc) GetAllPattern(ctx context.Context) ([]sqlc.SelectAllPatternRow, error) {
	rows, err := s.dao.SelectAllPattern(ctx)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
