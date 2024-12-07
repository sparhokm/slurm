package register

import (
	"context"

	"github.com/sparhokm/slurm/file-register/pkg/register_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sparhokm/slurm/file-storage/internal/domain"
)

type Client struct {
	conn register_v1.RegisterV1Client
}

func New(conn register_v1.RegisterV1Client) *Client {
	return &Client{conn: conn}
}

func (c Client) GetFile(ctx context.Context, fileID string) (*domain.FileRegister, error) {
	fileOut, err := c.conn.GetFile(ctx, &register_v1.GetFileIn{Id: fileID})
	if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
		return nil, domain.ErrFileNotFound
	}

	return &domain.FileRegister{
		ID:          fileID,
		Filepath:    fileOut.GetFile().GetFilepath(),
		Size:        fileOut.GetFile().GetSize(),
		Version:     fileOut.GetFile().GetVersion(),
		ContentType: fileOut.GetFile().GetContentType(),
	}, nil
}

func (c Client) FileExist(ctx context.Context, userID int64, filepath string) (bool, error) {
	_, err := c.conn.FindFileByPath(ctx, &register_v1.FindFileByPathIn{
		OwnerID: userID,
		Path:    filepath,
	})

	if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (c Client) AddFile(ctx context.Context, userID int64, fi domain.FileInfo) (string, error) {
	fileOut, err := c.conn.AddFile(ctx, &register_v1.AddFileIn{
		OwnerID:     userID,
		Filepath:    fi.Filepath,
		ContentType: fi.ContentType,
		Size:        fi.Size,
	})

	if err != nil {
		return "", err
	}

	return fileOut.GetId(), nil
}

func (c Client) DeleteFile(ctx context.Context, fileID string) error {
	_, err := c.conn.DeleteFile(ctx, &register_v1.DeleteFileIn{
		Id: fileID,
	})

	if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
		return domain.ErrFileNotFound
	}

	if err != nil {
		return err
	}

	return nil
}

func (c Client) UpdateFile(ctx context.Context, fileID string, size int64, version int64) error {
	_, err := c.conn.UpdateFile(ctx, &register_v1.UpdateFileIn{
		Id:      fileID,
		Size:    size,
		Version: version,
	})

	if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
		return domain.ErrFileNotFound
	}

	if err != nil {
		return err
	}

	return nil
}
