package test

import (
	"context"
	"microblog-backend/internal/mw/api"
	"microblog-backend/internal/mw/auth"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

type Input struct{}

type Output struct {
	Greeting string `json:"greeting"`
}

func (h *Handler) Handle(ctx context.Context, input *Input) api.Response {
	user := auth.FromContext(ctx)
	if user == nil {
		return api.Ok(Output{Greeting: "Hello, new friend!"})
	}

	return api.Ok(Output{
		Greeting: "wasssuuuuuup, " + user.Username,
	})
}
