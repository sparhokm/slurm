package handlers

import (
	fileStorage "github.com/sparhokm/slurm/file-storage/internal/generated/schema"
	"github.com/sparhokm/slurm/file-storage/pkg/logger"
)

type FileStorageImpl struct {
	log logger.Logger
	fs  FileStorage
	fr  FileRegister
}

var _ fileStorage.ServerInterface = (*FileStorageImpl)(nil)

func New(log logger.Logger, fs FileStorage, fr FileRegister) *FileStorageImpl {
	return &FileStorageImpl{
		log: log,
		fs:  fs,
		fr:  fr,
	}
}
