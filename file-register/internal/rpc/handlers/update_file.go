package handlers

import (
	"context"
	"errors"

	"github.com/sparhokm/slurm/file-register/internal/db/repository"
	"github.com/sparhokm/slurm/file-register/internal/rpc/interceptor"
	desc "github.com/sparhokm/slurm/file-register/pkg/register_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i Implementation) UpdateFile(ctx context.Context, fileIn *desc.UpdateFileIn) (*emptypb.Empty, error) {
	log := i.log.WithFields(map[string]any{
		"method":     "update_file",
		"request_id": interceptor.GetReqID(ctx),
	})

	if len(fileIn.GetId()) == 0 || fileIn.GetSize() == 0 || fileIn.GetVersion() == 0 {
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

	if file.Version != fileIn.GetVersion() {
		return nil, ErrMissFileVersion
	}

	file.Size = fileIn.GetSize()

	err = i.fileRepo.Update(ctx, file)
	if err != nil && errors.Is(err, repository.ErrMissVersion) {
		return nil, ErrMissFileVersion
	}
	if err != nil {
		log.WithError(err).Error("repo error")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
