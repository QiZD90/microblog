package signin

import (
	"context"
	"fmt"
	"microblog-backend/internal/mw/api"
	"microblog-backend/internal/mw/auth"
	"net/http"
	"time"
)

type Handler struct {
	credentialsRepo CredentialsRepo
	sessionRepo     SessionRepo
}

func New(credentialsRepo CredentialsRepo, sessionRepo SessionRepo) *Handler {
	return &Handler{
		credentialsRepo: credentialsRepo,
		sessionRepo:     sessionRepo,
	}
}

type Input struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Output struct {
	Success bool `json:"success"`
}

func (r *Handler) Handle(ctx context.Context, input *Input) api.Response {
	userId, err := r.credentialsRepo.GetUserIdByLoginAndPassword(ctx, input.Login, input.Password)
	if err != nil {
		return api.InternalServerError(fmt.Errorf("failed to get credentials: %w", err))
	}
	if userId == nil {
		return api.Unauthorized()
	}

	session, err := r.sessionRepo.CreateSession(ctx, *userId, time.Hour)
	if err != nil {
		return api.InternalServerError(fmt.Errorf("failed to get session: %w", err))
	}

	return api.Ok(Output{Success: true}, api.SetCookie(&http.Cookie{
		Name:  auth.SessionCookieKey,
		Value: session.Key,
	}))
}
