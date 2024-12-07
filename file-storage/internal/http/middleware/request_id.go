package middleware

import (
	"context"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func GetReqID(ctx context.Context) string {
	return chiMiddleware.GetReqID(ctx)
}
