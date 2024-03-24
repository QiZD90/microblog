package session

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"
	"microblog-backend/internal/model"
	"time"

	"github.com/redis/go-redis/v9"
)

type Repository struct {
	client *redis.Client
}

func New(client *redis.Client) *Repository {
	return &Repository{
		client: client,
	}
}

func (r *Repository) GetSession(ctx context.Context, sessionKey string) (*model.Session, error) {
	key := fmt.Sprintf("session_%s", sessionKey)

	userId, err := r.client.Get(ctx, key).Int64()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to get session from redis: %w", err)
	}

	ttl := r.client.TTL(ctx, key).Val()

	return &model.Session{
		UserId:    userId,
		Key:       sessionKey,
		ExpiresAt: time.Now().Add(ttl),
	}, nil
}

func (r *Repository) CreateSession(ctx context.Context, userId int64, ttl time.Duration) (*model.Session, error) {
	sessionKey := generateSessionKey()
	key := fmt.Sprintf("session_%s", sessionKey)

	err := r.client.Set(ctx, key, userId, ttl).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to write session to redis: %w", err)
	}

	return &model.Session{
		UserId:    userId,
		Key:       sessionKey,
		ExpiresAt: time.Now().Add(ttl),
	}, nil
}

func generateSessionKey() string {
	b := make([]byte, 48)
	rand.Read(b)

	return base32.StdEncoding.EncodeToString(b)
}
