package signin

import (
	"context"
	"microblog-backend/internal/model"
	"time"
)

type CredentialsRepo interface {
	GetUserIdByLoginAndPassword(ctx context.Context, login string, password string) (*int64, error)
}

type SessionRepo interface {
	CreateSession(ctx context.Context, userId int64, ttl time.Duration) (*model.Session, error)
}
