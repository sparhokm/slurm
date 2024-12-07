package handlers

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sparhokm/slurm/file-register/internal/models"
	"github.com/sparhokm/slurm/file-register/internal/rpc/interceptor"
	desc "github.com/sparhokm/slurm/file-register/pkg/register_v1"
)

func (i Implementation) AddFile(ctx context.Context, fileIn *desc.AddFileIn) (*desc.AddFileOut, error) {
	log := i.log.WithFields(map[string]any{
		"method":     "add_file",
		"request_id": interceptor.GetReqID(ctx),
	})

	u, err := uuid.NewV6()
	if err != nil {
		log.WithError(err).Error("uuid generate error")
		return nil, ErrInternal
	}

	if fileIn == nil || fileIn.GetOwnerID() == 0 ||
		len(fileIn.GetContentType()) == 0 || fileIn.GetSize() == 0 || len(fileIn.GetFilepath()) == 0 {
		return nil, ErrInvalidInput
	}

	file := models.File{
		ID:          u.String(),
		OwnerID:     fileIn.GetOwnerID(),
		ContentType: fileIn.GetContentType(),
		Size:        fileIn.GetSize(),
		Filepath:    fileIn.GetFilepath(),
		CreatedAt:   time.Now(),
	}

	err = i.fileRepo.Create(ctx, &file)
	if err != nil {
		log.WithError(err).Error("repo error")
		return nil, err
	}

	return &desc.AddFileOut{Id: file.ID}, nil
}
