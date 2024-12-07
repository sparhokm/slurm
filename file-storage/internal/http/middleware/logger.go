package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/sparhokm/slurm/file-storage/pkg/logger"
)

func NewLogger(log logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := log.WithFields(map[string]any{
				"method":      r.Method,
				"path":        r.URL.Path,
				"remote_addr": r.RemoteAddr,
				"user_agent":  r.UserAgent(),
				"request_id":  middleware.GetReqID(r.Context()),
			})

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				entry.WithFields(map[string]any{
					"status":   ww.Status(),
					"bytes":    ww.BytesWritten(),
					"duration": time.Since(t1).String(),
				}).Debug("request completed")
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
