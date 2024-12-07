package responses

import (
	"encoding/json"
	"net/http"

	FileStorage "github.com/sparhokm/slurm/file-storage/internal/generated/schema"
)

func InternalServerError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusInternalServerError)

	if len(message) == 0 {
		message = "internal server error"
	}

	_ = json.NewEncoder(w).Encode(FileStorage.Error{Message: &message})
}
