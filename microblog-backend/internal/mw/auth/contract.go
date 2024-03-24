package auth

import (
	"context"
	"microblog-backend/internal/model"
	"time"
)

type SessionRepo interface {
	GetSession(ctx context.Context, sessionKey string) (*model.Session, error)
	CreateSession(ctx context.Context, userId int64, ttl time.Duration) (*model.Session, error)
}

type UserRepo interface {
	GetUserById(ctx context.Context, userId int64) (*model.User, error)
}
