package middleware

import (
	"context"
	"net/http"
	"strings"

	"university-events/pkg/jwt"
	"university-events/pkg/response"
)

type contextKey string

const (
	
	ContextUserID contextKey = "user_id"
	ContextRole contextKey = "role"
)


func Auth(jwtManager *jwt.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" || !strings.HasPrefix(header, "Bearer ") {
				response.Error(w, http.StatusUnauthorized, "missing or malformed authorization header")
				return
			}

			tokenStr := strings.TrimPrefix(header, "Bearer ")
			claims, err := jwtManager.Parse(tokenStr)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), ContextUserID, claims.UserID)
			ctx = context.WithValue(ctx, ContextRole, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}


func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			got, _ := r.Context().Value(ContextRole).(string)
			if got != role {
				response.Error(w, http.StatusForbidden, "you do not have permission to perform this action")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
