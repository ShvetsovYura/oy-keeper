package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RecordStore struct {
	db *pgxpool.Pool
}

func NewRecordStore(db *pgxpool.Pool) *RecordStore {
	return &RecordStore{}
}

func (s *RecordStore) NewRecord(ctx context.Context, data any, attributes any, files any) error {
	tx, err := s.db.Begin(ctx)
	defer tx.Rollback(ctx)
	tx.Exec(`INSERT INTO "item"() values($1,$2,$3)`)
}

func (s *RecordStore) AddAttibute(ctx context.Context, data any) error {
	return nil
}

func (s *RecordStore) AddFileInfo(ctx context.Context, data any) error {
	return nil
}

func (s *RecordStore) AddAttributes(ctx context.Context, data any) error {
	return nil
}

func (s *RecordStore) AddFileInfos(ctx context.Context, data any) error {
	return nil
}
