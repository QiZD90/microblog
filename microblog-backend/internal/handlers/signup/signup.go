package signup

import (
	"context"
	"fmt"
	"microblog-backend/internal/db"
	"microblog-backend/internal/mw/api"
)

type Handler struct {
	db              *db.DB
	userRepo        UserRepo
	credentialsRepo CredentialsRepo
}

func New(db *db.DB, userRepo UserRepo, credentialsRepo CredentialsRepo) *Handler {
	return &Handler{
		db:              db,
		userRepo:        userRepo,
		credentialsRepo: credentialsRepo,
	}
}

type Input struct {
	Username string `json:"username"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Output struct {
	Success bool `json:"success"`
}

func (h *Handler) Handle(ctx context.Context, input *Input) api.Response {
	err := h.db.WithTx(ctx, func(ctx context.Context) error {
		userId, err := h.userRepo.CreateUser(ctx, input.Username)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		if err := h.credentialsRepo.StoreCredentials(ctx, userId, input.Login, input.Password); err != nil {
			return fmt.Errorf("failed to store credentials: %w", err)
		}

		return nil
	})
	if err != nil {
		return api.InternalServerError(fmt.Errorf("failed to sign up: %w", err))
	}

	return api.Ok(Output{Success: true})
}
