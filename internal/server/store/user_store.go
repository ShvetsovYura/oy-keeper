package store

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ShvetsovYura/oykeeper/internal/logger"
	"github.com/ShvetsovYura/oykeeper/internal/server/store/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stephenafamo/scan"
	"github.com/stephenafamo/scan/pgxscan"
)

type UserStore struct {
	db *pgxpool.Pool
}

func NewUserStore(db *pgxpool.Pool) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) UpsertUser(ctx context.Context, record *models.UserDB) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error on create tx %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO item ("uuid", "login", otp_secret, otp_auth, otp_verified, is_active, created_at, updated_at) 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT ("uuid", "login") DO UPDATE 
		SET "otp_secret"=EXCLUDED."otp_secret", otp_auth=EXCLUDED.otp_auth, otp_verified=EXCLUDED.otp_verified, 
			is_active = EXCLUDED.is_active, updated_at=now()
			;
	`, record.Uuid, record.Login, record.Otp_secret, record.Otp_auth, record.Otp_verified, record.Created_at, record.Updated_at)

	if err != nil {
		logger.Log.Error(err.Error())
		return fmt.Errorf("error on exec tx %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error on commit %w", err)
	}
	return nil
}

func (s *UserStore) DeleteUser(ctx context.Context, userUUID string) error {
	user, err := s.GetUser(ctx, userUUID)
	if user == nil {
		logger.Log.Warn("user not found", slog.String("userUUID", userUUID), slog.String("msg", err.Error()))
		return nil
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error on create tx %w", err)
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, "DELETE FROM user WHERE uuid = $1", userUUID)
	if err != nil {
		return fmt.Errorf("error on delete user %w", err)
	}
	_, err = tx.Exec(ctx, "DELETE FROM item WHERE user_uuid = $1", userUUID)
	if err != nil {
		return fmt.Errorf("error on delete user records %w", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error on commit %w", err)
	}
	return nil
}

func (s *UserStore) GetUser(ctx context.Context, userUUID string) (*models.UserDB, error) {
	rec, err := pgxscan.One(ctx, s.db, scan.StructMapper[*models.UserDB](), "SELECT * FROM user where uuid = $1", userUUID)
	if err != nil {
		return nil, err
	}

	return rec, nil
}
