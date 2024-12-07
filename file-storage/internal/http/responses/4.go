package responses

import (
	"encoding/json"
	"net/http"

	fileStorage "github.com/sparhokm/slurm/file-storage/internal/generated/schema"
)

func BadRequest(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)

	if len(message) == 0 {
		message = "bad request"
	}

	_ = json.NewEncoder(w).Encode(fileStorage.Error{Message: &message})
}

func Forbidden(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
}

func Conflict(w http.ResponseWriter) {
	w.WriteHeader(http.StatusConflict)
}

func NotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}
