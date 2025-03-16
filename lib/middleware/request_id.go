package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// Contexts best practice is to use typed keys, since
// using plain strings can lead to collisions with other packages
type contextKey string

const reqIDKey contextKey = "request_id"

func WithReqID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), reqIDKey, uuid.New())
		next.ServeHTTP(w, r.Clone(ctx))
	})
}
