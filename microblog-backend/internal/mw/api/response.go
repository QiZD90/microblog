package api

import (
	"context"
	"encoding/json"
	"fmt"
	"microblog-backend/internal/logger"
	"net/http"
)

type Response interface {
	Write(ctx context.Context, w http.ResponseWriter)
}

type errorResponse struct {
	err          error
	statusCode   int64
	errorMessage string
}

func (r errorResponse) Write(ctx context.Context, w http.ResponseWriter) {
	log := logger.FromContext(ctx).WithFields(map[string]any{
		"status_code": r.statusCode,
	})
	log.Errorf(r.err.Error())

	w.WriteHeader(int(r.statusCode))
	w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, r.errorMessage)))
}

func BadRequest(err error) Response {
	return errorResponse{
		err:          err,
		statusCode:   http.StatusBadRequest,
		errorMessage: "bad request",
	}
}

func InternalServerError(err error) Response {
	return errorResponse{
		err:          err,
		statusCode:   http.StatusInternalServerError,
		errorMessage: "internal server error",
	}
}

func Unauthorized() Response {
	return errorResponse{
		err:          fmt.Errorf("user is unauthorized"),
		statusCode:   http.StatusUnauthorized,
		errorMessage: "unauthorized",
	}
}

type okResponse[T any] struct {
	result  T
	options []Option
}

func (r okResponse[T]) Write(ctx context.Context, w http.ResponseWriter) {
	outputBytes, err := json.Marshal(r.result)
	if err != nil {
		InternalServerError(err).Write(ctx, w)
		return
	}

	log := logger.FromContext(ctx).WithFields(map[string]any{
		"status_code": http.StatusOK,
	})
	log.Infof("ok")

	for _, option := range r.options {
		option(w)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(outputBytes)
}

func Ok[T any](result T, options ...Option) Response {
	return okResponse[T]{
		result:  result,
		options: options,
	}
}
