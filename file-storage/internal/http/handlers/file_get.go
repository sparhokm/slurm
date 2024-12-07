package handlers

import (
	"errors"
	"net/http"

	"github.com/sparhokm/slurm/file-storage/internal/domain"
	"github.com/sparhokm/slurm/file-storage/internal/http/responses"

	"github.com/sparhokm/slurm/file-storage/internal/http/middleware"
)

func (h *FileStorageImpl) FileGet(w http.ResponseWriter, r *http.Request, fileID string) {
	log := h.log.WithFields(map[string]any{
		"method":     r.Method,
		"file_id":    fileID,
		"request_id": middleware.GetReqID(r.Context()),
	})

	fileR, err := h.fr.GetFile(r.Context(), fileID)
	if err != nil && errors.Is(err, domain.ErrFileNotFound) {
		responses.NotFound(w)
		return
	}
	if err != nil {
		log.WithError(err).Error("error find file in register")
		responses.InternalServerError(w, "")
		return
	}

	log.Debug("Success")
	responses.FileInfo(w, fileR)
}
