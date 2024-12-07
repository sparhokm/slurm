package handlers

import (
	"context"
	"io"

	"github.com/sparhokm/slurm/file-storage/internal/domain"
)

type FileStorage interface {
	FileInfo(ctx context.Context, userID int64, filepath string) (*domain.FileInfo, error)
	DeleteFile(ctx context.Context, userID int64, filepath string) error
	UploadFile(
		ctx context.Context,
		userID int64,
		filepath string,
		file io.Reader,
		contentType string,
	) (*domain.FileInfo, error)
	DownloadFile(ctx context.Context, userID int64, filepath string) (*domain.FileInfo, io.Reader, error)
}

type FileRegister interface {
	GetFile(ctx context.Context, fileID string) (*domain.FileRegister, error)
	FileExist(ctx context.Context, userID int64, filepath string) (bool, error)
	AddFile(ctx context.Context, userID int64, fi domain.FileInfo) (string, error)
	DeleteFile(ctx context.Context, fileID string) error
	UpdateFile(ctx context.Context, fileID string, size int64, version int64) error
}
