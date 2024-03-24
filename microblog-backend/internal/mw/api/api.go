package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func MethodNotAllowed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		errorResponse{
			err:          fmt.Errorf("method not allowed"),
			statusCode:   http.StatusMethodNotAllowed,
			errorMessage: "method not allowed",
		}.Write(r.Context(), w)
	}
}

func NotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		errorResponse{
			err:          fmt.Errorf("not found"),
			statusCode:   http.StatusNotFound,
			errorMessage: "not found",
		}.Write(r.Context(), w)
	}
}

func HandlerMW[I any](handler func(ctx context.Context, input *I) Response) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		response := func() Response {
			var input I

			if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
				return BadRequest(fmt.Errorf("failed to unmarshal input: %w", err))
			}

			return handler(ctx, &input)
		}()

		response.Write(ctx, w)
	}
}
