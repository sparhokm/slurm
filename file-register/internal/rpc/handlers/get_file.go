package handlers

import (
	"context"
	"errors"

	"github.com/sparhokm/slurm/file-register/internal/db/repository"
	"github.com/sparhokm/slurm/file-register/internal/rpc/interceptor"
	desc "github.com/sparhokm/slurm/file-register/pkg/register_v1"
)

func (i Implementation) GetFile(ctx context.Context, fileIn *desc.GetFileIn) (*desc.GetFileOut, error) {
	log := i.log.WithFields(map[string]any{
		"method":     "get_file",
		"request_id": interceptor.GetReqID(ctx),
	})

	if len(fileIn.GetId()) == 0 {
		return nil, ErrInvalidInput
	}

	file, err := i.fileRepo.GetByID(ctx, fileIn.GetId())
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		log.WithError(err).Error("repo error")
		return nil, err
	}

	return &desc.GetFileOut{File: toFileOut(file)}, nil
}
