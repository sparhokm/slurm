package handlers

import (
	"context"
	"errors"

	"github.com/sparhokm/slurm/file-register/internal/db/repository"
	"github.com/sparhokm/slurm/file-register/internal/rpc/interceptor"
	desc "github.com/sparhokm/slurm/file-register/pkg/register_v1"
)

func (i Implementation) FindFileByPath(
	ctx context.Context,
	fileIn *desc.FindFileByPathIn,
) (*desc.FindFileByPathOut, error) {
	log := i.log.WithFields(map[string]any{
		"method":     "find_file",
		"request_id": interceptor.GetReqID(ctx),
	})

	if len(fileIn.GetPath()) == 0 || fileIn.GetOwnerID() == 0 {
		return nil, ErrInvalidInput
	}

	file, err := i.fileRepo.FindByFilepath(ctx, fileIn.GetOwnerID(), fileIn.GetPath())
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		log.WithError(err).Error("repo error")
		return nil, err
	}

	return &desc.FindFileByPathOut{File: toFileOut(file)}, nil
}
