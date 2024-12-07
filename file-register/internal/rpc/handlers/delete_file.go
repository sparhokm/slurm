package handlers

import (
	"context"
	"errors"

	"github.com/sparhokm/slurm/file-register/internal/db/repository"
	"github.com/sparhokm/slurm/file-register/internal/rpc/interceptor"
	desc "github.com/sparhokm/slurm/file-register/pkg/register_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i Implementation) DeleteFile(ctx context.Context, fileIn *desc.DeleteFileIn) (*emptypb.Empty, error) {
	log := i.log.WithFields(map[string]any{
		"method":     "delete_file",
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

	err = i.fileRepo.Delete(ctx, file)
	if err != nil {
		log.WithError(err).Error("repo error")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
