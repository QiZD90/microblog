package credentials

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"errors"
	"fmt"
	"io"
	"microblog-backend/internal/db"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type Repository struct {
	db *db.DB
}

func New(db *db.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetUserIdByLoginAndPassword(ctx context.Context, login string, password string) (*int64, error) {
	query, args, err := sq.Select("user_id").
		From("credentials").
		Where(sq.Eq{
			"login":         login,
			"password_hash": hashPassword(password),
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var userId int64
	if err := r.db.Get(ctx).GetContext(ctx, &userId, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to find user by credentials: %w", err)
	}

	return &userId, nil
}

func (r *Repository) UserWithLoginExists(ctx context.Context, login string) (bool, error) {
	query, args, err := sq.Select("count(*) > 0").
		From("credentials").
		Where(sq.Eq{
			"login": login,
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("failed to build query: %w", err)
	}

	var exists bool
	if err := r.db.Get(ctx).GetContext(ctx, &exists, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, fmt.Errorf("failed to check if user exists: %w", err)
	}

	return exists, nil
}

func (r *Repository) StoreCredentials(ctx context.Context, userId int64, login string, password string) error { // todo: add invite link
	query, args, err := sq.Insert("credentials").
		Columns(
			"user_id",
			"login",
			"password_hash",
			"created_at",
		).
		Values(
			userId,
			login,
			hashPassword(password),
			time.Now(),
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	if _, err := r.db.Get(ctx).ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to store credentials: %w", err)
	}

	return nil
}

func hashPassword(password string) string {
	sha := sha256.New()
	io.WriteString(sha, password)
	b := sha.Sum(nil)

	return base32.StdEncoding.EncodeToString(b)
}
