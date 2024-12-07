package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/sparhokm/slurm/file-storage/internal/domain"
	fileStorage "github.com/sparhokm/slurm/file-storage/internal/generated/schema"
	"github.com/sparhokm/slurm/file-storage/internal/http/middleware"
	"github.com/sparhokm/slurm/file-storage/internal/http/responses"
)

func (h *FileStorageImpl) FileUpload( //nolint:funlen,gocognit
	w http.ResponseWriter,
	r *http.Request,
	params fileStorage.FileUploadParams,
) {
	log := h.log.WithFields(map[string]any{
		"method":     r.Method,
		"request_id": middleware.GetReqID(r.Context()),
	})

	userID := middleware.GetUserID(r.Context())
	if userID == nil {
		responses.Forbidden(w)
		return
	}

	if params.Filepath == "" || !strings.HasPrefix(params.Filepath, "/") || filepath.Ext(params.Filepath) == "" {
		log.WithField("filepath", params.Filepath).Warn("wrong filepath")
		responses.BadRequest(w, "wrong filepath")
		return
	}

	if len(params.Filepath) != len(clearFilepath(params.Filepath)) {
		log.WithField("filepath", params.Filepath).Warn("wrong filepath")
		responses.BadRequest(w, "wrong filepath")
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
	)
	for {
		part, errPart := form.NextPart()
		if errPart == io.EOF {
			break
		}
		if part.FormName() == "file" {
			contentType := part.Header.Get("Content-Type")
			if contentType == "" {
				log.WithError(err).Error("content-type is empty")
				responses.BadRequest(w, "Content-Type required")
				return
			}

			exist, err := h.fr.FileExist(r.Context(), *userID, params.Filepath) //nolint:govet
			if err != nil {
				log.WithError(err).Error("error check exist file in register")
				responses.InternalServerError(w, "")
				return
			}
			if exist {
				log.Warn("file exist")
				responses.Conflict(w)
				return
			}

			fileInfo, err = h.uploadFile(r.Context(), *userID, params.Filepath, part, contentType)
			if err != nil {
				log.WithError(err).Error("upload file error")
				responses.InternalServerError(w, "")
				return
			}

			part.Close()
			break
		}
	}

	fileID, err := h.fr.AddFile(r.Context(), *userID, *fileInfo)
	if err != nil {
		log.WithError(err).Error("add file to register error")
		responses.InternalServerError(w, "")
		return
	}

	responses.FileCreated(w, fileID, fileInfo)
}

func (h *FileStorageImpl) uploadFile(
	ctx context.Context,
	userID int64,
	filepath string,
	buf io.Reader,
	contentType string,
) (*domain.FileInfo, error) {
	fi, err := h.fs.FileInfo(ctx, userID, filepath)
	if err != nil && !errors.Is(err, domain.ErrFileNotFound) {
		return nil, err
	}

	if fi != nil {
		return fi, nil
	}

	fi, err = h.fs.UploadFile(ctx, userID, filepath, buf, contentType)
	if err != nil {
		return nil, err
	}

	return fi, nil
}

func clearFilepath(fp string) string {
	separator := string(filepath.Separator)
	clearFp := filepath.Dir(fp)
	if clearFp != separator {
		clearFp += separator
	}
	clearFp += filepath.Base(fp)

	return clearFp
}
