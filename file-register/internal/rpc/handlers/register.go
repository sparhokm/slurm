package handlers

import (
	desc "github.com/sparhokm/slurm/file-register/pkg/register_v1"
	"github.com/sparhokm/slurm/file-storage/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrInvalidInput = status.Error(codes.InvalidArgument, "InvalidInput")
var ErrInternal = status.Error(codes.Internal, "Internal error")
var ErrNotFound = status.Error(codes.NotFound, "Not found")
var ErrMissFileVersion = status.Error(codes.Aborted, "Miss file version")

type Implementation struct {
	desc.UnimplementedRegisterV1Server
	fileRepo FileRepository
	log      logger.Logger
}

func New(fileRepo FileRepository, log logger.Logger) *Implementation {
	return &Implementation{fileRepo: fileRepo, log: log}
}
