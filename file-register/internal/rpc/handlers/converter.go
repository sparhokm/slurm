package handlers

import (
	"github.com/sparhokm/slurm/file-register/internal/models"
	"github.com/sparhokm/slurm/file-register/pkg/register_v1"
)

func toFileOut(file *models.File) *register_v1.File {
	return &register_v1.File{
		Id:          file.ID,
		OwnerID:     file.OwnerID,
		ContentType: file.ContentType,
		Size:        file.Size,
		Filepath:    file.Filepath,
		Version:     file.Version,
	}
}
