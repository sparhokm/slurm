package handlers

import (
	"net/http"

	"github.com/sparhokm/slurm/file-storage/internal/http/middleware"
)

func (h *FileStorageImpl) AuthRefreshToken(w http.ResponseWriter, r *http.Request) {
	log := h.log.WithFields(map[string]any{
		"method":     r.Method,
		"request_id": middleware.GetReqID(r.Context()),
	})

	// todo отправляем по gRPC refresh token в сервис Auth
	//  если в сервисе auth видим, что этот refresh token не использовали,
	//  то формируем новый jwt и refresh token и возвращаем их клиенту,
	//  а старый refresh token помечаем как использованный

	log.Debug("Success")
	w.WriteHeader(http.StatusOK)
}
