package handlers

import (
	"io"
	"net/http"

	"github.com/sparhokm/slurm/file-storage/internal/domain"
	"github.com/sparhokm/slurm/file-storage/internal/http/responses"

	"github.com/sparhokm/slurm/file-storage/internal/http/middleware"
)

func (h *FileStorageImpl) FileUpdate(w http.ResponseWriter, r *http.Request, fileID string) {
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

	form, err := r.MultipartReader()
	if err != nil {
		log.WithError(err).Warn("can't parse multipart form")
		responses.BadRequest(w, "")
		return
	}

	var (
		fileInfo *domain.FileInfo
		fileR    *domain.FileRegister
	)
	for {
		part, errPart := form.NextPart()
		if errPart == io.EOF {
			break
		}
		if part.FormName() != "file" {
			continue
		}

		contentType := part.Header.Get("Content-Type")
		if contentType == "" {
			responses.BadRequest(w, "Content-Type required")
			return
		}

		fileR, err = h.fr.GetFile(r.Context(), fileID)
		if err != nil {
			log.WithError(err).Error("error check exist file in register")
			responses.InternalServerError(w, "")
			return
		}

		fileInfo, err = h.fs.UploadFile(r.Context(), *userID, fileR.Filepath, part, fileR.ContentType)
		if err != nil {
			log.WithError(err).Error("upload file error")
			responses.InternalServerError(w, "")
			return
		}

		part.Close()
		break
	}

	err = h.fr.UpdateFile(r.Context(), fileID, fileInfo.Size, fileR.Version)
	if err != nil {
		log.WithError(err).Error("update file to register error")
		responses.InternalServerError(w, "")
		return
	}

	log.Debug("Success")

	responses.FileUpdated(w)
}
