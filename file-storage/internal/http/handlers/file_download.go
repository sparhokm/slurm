package handlers

import (
	"errors"
	"net/http"

	"github.com/sparhokm/slurm/file-storage/internal/domain"
	"github.com/sparhokm/slurm/file-storage/internal/http/middleware"
	"github.com/sparhokm/slurm/file-storage/internal/http/responses"
)

func (h *FileStorageImpl) FileDownload(w http.ResponseWriter, r *http.Request, fileID string) {
	log := h.log.WithFields(map[string]any{
		"method":     r.Method,
		"file_id":    fileID,
		"request_id": middleware.GetReqID(r.Context()),
	})

	userID := middleware.GetUserID(r.Context())
	if userID == nil {
		responses.Forbidden(w)
		return
	}

	fileR, err := h.fr.GetFile(r.Context(), fileID)
	if err != nil && errors.Is(err, domain.ErrFileNotFound) {
		responses.NotFound(w)
		return
	} else if err != nil {
		log.WithError(err).Error("error find file in register")
		responses.InternalServerError(w, "")
		return
	}

	fileInfo, file, err := h.fs.DownloadFile(r.Context(), *userID, fileR.Filepath)
	if err != nil {
		log.WithError(err).Error("download file error")
		responses.InternalServerError(w, "")
		return
	}

	responses.FileDownload(w, file, fileInfo)
}
