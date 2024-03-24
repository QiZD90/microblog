package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"microblog-backend/internal/db"
	"microblog-backend/internal/model"
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

func (r *Repository) GetUserById(ctx context.Context, userId int64) (*model.User, error) {
	query, args, err := sq.Select(
		"id",
		"username",
		"created_at",
		"invite_id",
	).
		From("users").
		Where(sq.Eq{
			"id": userId,
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var user model.User
	if err := r.db.Get(ctx).GetContext(ctx, &user, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}

	return &user, nil
}

func (r *Repository) CreateUser(ctx context.Context, username string) (int64, error) {
	query, args, err := sq.Insert("users").
		Columns(
			"username",
			"created_at",
		).
		Values(
			username,
			time.Now(),
		).
		Suffix("returning id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build query: %w", err)
	}

	var userId int64
	if err := r.db.Get(ctx).GetContext(ctx, &userId, query, args...); err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return userId, nil
}
