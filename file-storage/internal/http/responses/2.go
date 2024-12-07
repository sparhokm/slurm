package responses

import (
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"

	"github.com/sparhokm/slurm/file-storage/internal/domain"
	fileStorage "github.com/sparhokm/slurm/file-storage/internal/generated/schema"
)

func FileCreated(w http.ResponseWriter, fileID string, fi *domain.FileInfo) {
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(fileStorage.FileInfo{Size: int(fi.Size), Id: fileID})
}

func FileUpdated(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func FileDownload(w http.ResponseWriter, file io.Reader, fi *domain.FileInfo) {
	w.Header().Set("Content-Type", fi.ContentType)
	w.Header().Set("Content-Disposition", "attachment; "+filepath.Base(fi.Filepath))
	w.WriteHeader(http.StatusOK)

	io.Copy(w, file) //nolint:errcheck
}

func FileInfo(w http.ResponseWriter, f *domain.FileRegister) {
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(fileStorage.FileInfo{Size: int(f.Size), Id: f.ID})
}

func FileDeleted(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
