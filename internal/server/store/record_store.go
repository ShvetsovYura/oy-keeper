package store

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/ShvetsovYura/oykeeper/internal/logger"
	pb "github.com/ShvetsovYura/oykeeper/proto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RecordStore struct {
	db *pgxpool.Pool
	// conn *pgx.Conn
}

func NewRecordStore(db *pgxpool.Pool) *RecordStore {
	return &RecordStore{
		db: db,
		// conn: &pgx.new,
	}
}

func (s *RecordStore) NewRecord(ctx context.Context, record *pb.RecordReq, attributes []*pb.AttributeInfo, files []*pb.FileInfo) error {
	if record == nil {
		return errors.New("reqord equal nil")
	}

	logger.Log.Debug("incoming attributes", slog.Any("attrs", attributes))
	tx, _ := s.db.Begin(ctx)
	defer tx.Rollback(ctx)

	_, err := tx.Exec(ctx, `
		INSERT INTO item ("uuid", "name", username, "password", url, expired_at, created_at, updated_at, user_uuid, cardnum, description, "version") 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);
	`, record.Uuid, record.Name, record.Username, record.Password, record.Url, record.ExpiredAt, time.Now(), time.Now(), record.UserUuid, record.CardNum, record.Description, record.Version)

	if err != nil {
		logger.Log.Error(err.Error())
	}

	for _, a := range attributes {
		_, err = tx.Exec(ctx, `
			INSERT INTO "attribute" ("uuid", item_uuid, "name", value) 
			VALUES($1, $2, $3, $4) --ON CONFLICT DO UPDATE 
			--SET "name" = EXCLUDED."name", value = EXCLUDED.value
			;`,
			a.Uuid, record.Uuid, a.Name, a.Value)

		if err != nil {
			fmt.Printf("error on exec sub %v", err)
		}
	}

	for _, f := range files {
		_, err = tx.Exec(ctx, `
		INSERT INTO file ("uuid", item_uuid, "path", hash, "size", "name", created_at, updated_at, meta) 
		VALUES($1, $2,$3,$4,$5,$6,$7,$8,$9);
		`, f.Uuid, record.Uuid, f.Path, f.Hash, f.Size, f.Name, time.Now(), time.Now(), f.Meta)

		if err != nil {
			fmt.Printf("error on exec f %v", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		fmt.Printf("error on commit %v", err)
	}
	return nil
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
