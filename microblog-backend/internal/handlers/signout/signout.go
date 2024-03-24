package signout

import (
	"context"
	"microblog-backend/internal/mw/api"
	"microblog-backend/internal/mw/auth"
	"net/http"
)

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

type Input struct{}

type Output struct {
	Success bool `json:"success"`
}

func (r *Handler) Handle(ctx context.Context, input *Input) api.Response {
	return api.Ok(Output{Success: true}, api.SetCookie(&http.Cookie{
		Name:   auth.SessionCookieKey,
		Value:  "none",
		MaxAge: 0,
	}))
}
