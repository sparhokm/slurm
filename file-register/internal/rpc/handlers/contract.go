package handlers

import (
	"context"

	"github.com/sparhokm/slurm/file-register/internal/models"
)

type FileRepository interface {
	Create(ctx context.Context, file *models.File) error
	GetByID(ctx context.Context, id string) (*models.File, error)
	FindByFilepath(ctx context.Context, ownerID int64, filepath string) (*models.File, error)
	Update(ctx context.Context, file *models.File) error
	Delete(ctx context.Context, file *models.File) error
}
