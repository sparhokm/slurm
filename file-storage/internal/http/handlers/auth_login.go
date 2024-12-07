package handlers

import (
	"net/http"

	"github.com/sparhokm/slurm/file-storage/internal/http/middleware"
)

func (h *FileStorageImpl) AuthLogin(w http.ResponseWriter, r *http.Request) {
	log := h.log.WithFields(map[string]any{
		"method":     r.Method,
		"request_id": middleware.GetReqID(r.Context()),
	})

	// todo отправляем по gRPC username и password в сервис User для проверки
	//  в случаи успеха по gRPC отправляем userID в сервис Auth
	//  и получаем от него jwt token (ttl = 5m) и refresh token(ttl = 20d)
	//  в сервисе auth дополнительно сохраняем refresh token в БД

	log.Debug("Success")
	w.WriteHeader(http.StatusOK)
}
