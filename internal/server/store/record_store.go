package store

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/ShvetsovYura/oykeeper/internal/logger"
	pb "github.com/ShvetsovYura/oykeeper/proto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stephenafamo/scan"
	"github.com/stephenafamo/scan/pgxscan"
)

var NoRecordError = errors.New("no record in result")
var OldRecordVersionError = errors.New("record has old versions")

type RecordStore struct {
	db *pgxpool.Pool
}

func NewRecordStore(db *pgxpool.Pool) *RecordStore {
	return &RecordStore{
		db: db,
	}
}
func (s *RecordStore) GetRecordVersion(ctx context.Context, recordUUID string, recordVersion uint32) (*RecordByUUID, error) {
	r, err := pgxscan.One(ctx, s.db, scan.StructMapper[*RecordByUUID](), "SELECT uuid, version FROM item WHERE uuid=$1", recordUUID)

	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *RecordStore) UpsertRecord(ctx context.Context, record *pb.RecordReq, attributes []*pb.AttributeInfo, files []*pb.FileInfo) error {
	if record == nil {
		return errors.New("reqord equal nil")
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error on create tx %w", err)
	}
	defer tx.Rollback(ctx)

	// TODO: Добавить lock
	_, err = tx.Exec(ctx, `
		INSERT INTO item ("uuid", "name", username, "password", url, expired_at, created_at, updated_at, user_uuid, cardnum, description, "version") 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT ("uuid") DO UPDATE 
		SET "name"=EXCLUDED."name", username=EXCLUDED.username, "password"=EXCLUDED."password", 
			expired_at = EXCLUDED.expired_at, cardnum=EXCLUDED.cardnum, description=EXCLUDED.description, 
			"version"=EXCLUDED."version", updated_at=now()
			;
	`, record.Uuid, record.Name, record.Username, record.Password, record.Url, nil, time.Now(), time.Now(), record.UserUuid, record.CardNum, record.Description, record.Version+1)

	if err != nil {
		logger.Log.Error(err.Error())
	}
	err = upsertRecordAttributes(ctx, tx, record, attributes)
	if err != nil {
		return err
	}
	err = upsertRecordFileInfos(ctx, tx, record, files)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error on commit %w", err)
	}
	return nil
}

func upsertRecordAttributes(ctx context.Context, tx pgx.Tx, record *pb.RecordReq, attributes []*pb.AttributeInfo) error {
	var attrUUIDs []string
	for _, a := range attributes {
		attrUUIDs = append(attrUUIDs, a.Uuid)
	}
	_, err := tx.Exec(ctx, `DELETE FROM attribute WHERE item_uuid = $1 AND "uuid" != any($2);`, record.Uuid, attrUUIDs)
	if err != nil {
		return fmt.Errorf("error on delete attributes %w", err)
	}
	for _, a := range attributes {
		_, err = tx.Exec(ctx, `
			INSERT INTO "attribute" ("uuid", item_uuid, "name", value) 
			VALUES($1, $2, $3, $4) ON CONFLICT ("uuid") DO UPDATE
			SET "name" = EXCLUDED."name", value = EXCLUDED.value
			;`,
			a.Uuid, record.Uuid, a.Name, a.Value)

		if err != nil {
			return fmt.Errorf("error on insert attribute %w", err)
		}
	}
	return nil
}

func upsertRecordFileInfos(ctx context.Context, tx pgx.Tx, record *pb.RecordReq, files []*pb.FileInfo) error {
	var filesUUIDs []string
	for _, f := range files {
		filesUUIDs = append(filesUUIDs, f.Uuid)
	}
	_, err := tx.Exec(ctx, `DELETE FROM file WHERE item_uuid = $1 AND "uuid" != any($2);`, record.Uuid, filesUUIDs)
	if err != nil {
		return fmt.Errorf("error on delete file infos %w", err)
	}
	for _, f := range files {
		_, err = tx.Exec(ctx, `
		INSERT INTO file ("uuid", item_uuid, "path", hash, "size", "name", created_at, updated_at, meta) 
		VALUES($1, $2,$3,$4,$5,$6,$7,$8,$9) ON CONFLICT ("uuid") DO UPDATE 
		SET "path"=EXCLUDED."path", hash=EXCLUDED.hash, size=EXCLUDED.size, updated_at=now(), meta=EXCLUDED.meta
		;
		`, f.Uuid, record.Uuid, f.Path, f.Hash, f.Size, f.Name, time.Now(), time.Now(), f.Meta)

		if err != nil {
			return fmt.Errorf("error on insert file_info %w", err)
		}
	}
	return nil
}

func (s *RecordStore) FetchUserRecords(ctx context.Context, userUUID string) ([]RecordDB, error) {
	recs, err := pgxscan.All(ctx, s.db, scan.StructMapper[RecordDB](),
		`SELECT *  FROM item WHERE user_uuid = $1`, userUUID)

	if err != nil {
		logger.Log.Info("error on list records", slog.String("error", err.Error()))

	}
	logger.Log.Info("count records", slog.Int("count", len(recs)))
	return recs, nil
}

func (s *RecordStore) FetchUserAttributes(ctx context.Context, userUUID string) ([]AttributeDB, error) {
	recs, err := pgxscan.All(ctx, s.db, scan.StructMapper[AttributeDB](),
		`SELECT a.*  FROM attribute a join item i on a.item_uuid = i.uuid WHERE i.user_uuid = $1`, userUUID)

	if err != nil {
		logger.Log.Info("error on list records", slog.String("error", err.Error()))

	}
	logger.Log.Info("count records", slog.Int("count", len(recs)))
	return recs, nil
}

func (s *RecordStore) FetchUserFilesInfos(ctx context.Context, userUUID string) ([]FileDB, error) {
	recs, err := pgxscan.All(ctx, s.db, scan.StructMapper[FileDB](),
		`SELECT f.*  FROM file f join item i on f.item_uuid = i.uuid WHERE i.user_uuid = $1`, userUUID)

	if err != nil {
		logger.Log.Info("error on list records", slog.String("error", err.Error()))
	}
	logger.Log.Info("count records", slog.Int("count", len(recs)))
	return recs, nil
}

func (s *RecordStore) RemoveRecord(ctx context.Context, uuid string) error {
	_, err := s.db.Exec(ctx, "DELETE FROM item WHERE uuid=$1", uuid)
	return err
}

func (s *RecordStore) RemoveAttribute(ctx context.Context, uuid string) error {
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
