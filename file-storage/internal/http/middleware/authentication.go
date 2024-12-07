package middleware

import (
	"context"
	"net/http"

	"github.com/sparhokm/slurm/file-storage/pkg/logger"
)

type userIDKey int

const UserIDKey userIDKey = 0

func NewAuthentication(_ logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if authHeader := r.Header.Get("Authorization"); authHeader != "" {
				// todo: отправляем jwt токен в Auth сервис и он уже возвращает UserID
				var uID int64 = 45
				r = r.WithContext(context.WithValue(r.Context(), UserIDKey, uID))
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func GetUserID(ctx context.Context) *int64 {
	uID := ctx.Value(UserIDKey)

	if id, ok := uID.(int64); ok {
		return &id
	}

	return nil
}
