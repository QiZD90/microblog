package auth

import (
	"context"
	"errors"
	"fmt"
	"microblog-backend/internal/model"
	"microblog-backend/internal/mw/api"
	"net/http"
	"time"
)

const SessionCookieKey = "session_key"

type Middleware struct {
	sessionRepo SessionRepo
	userRepo    UserRepo
}

func New(sessionRepo SessionRepo, userRepo UserRepo) *Middleware {
	return &Middleware{
		sessionRepo: sessionRepo,
		userRepo:    userRepo,
	}
}

type ctxKey struct{}

func (m *Middleware) AuthMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, err := m.authContextFromRequest(r.Context(), r)
		if err != nil {
			api.InternalServerError(err).Write(ctx, w)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) authContextFromRequest(ctx context.Context, r *http.Request) (context.Context, error) {
	if ctx.Value(ctxKey{}) != nil {
		return ctx, nil
	}

	cookie, err := r.Cookie(SessionCookieKey)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return ctx, nil
		}

		return ctx, fmt.Errorf("failed to get cookie: %w", err)
	}

	session, err := m.sessionRepo.GetSession(ctx, cookie.Value)
	if err != nil {
		return ctx, fmt.Errorf("failed to get session: %w", err)
	}
	if session == nil || session.ExpiresAt.Before(time.Now()) {
		return ctx, nil
	}

	user, err := m.userRepo.GetUserById(ctx, session.UserId)
	if err != nil {
		return ctx, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return ctx, fmt.Errorf("failed to get user: user doesn't exist")
	}

	return context.WithValue(ctx, ctxKey{}, user), nil
}

func FromContext(ctx context.Context) *model.User {
	switch v := ctx.Value(ctxKey{}).(type) {
	case *model.User:
		return v
	default:
		return nil
	}
}
