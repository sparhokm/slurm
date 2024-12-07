package minio

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sparhokm/slurm/file-storage/internal/domain"
	"github.com/sparhokm/slurm/file-storage/internal/otel"
)

type Minio struct {
	bucket string
	client *minio.Client
}

func New(cfg Config) (*Minio, error) {
	minioClient, err := minio.New(cfg.GetEndpoint(), &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.GetAccessKey(), cfg.GetSecretKey(), ""),
		Secure: cfg.GetUseSSL(),
	})
	if err != nil {
		return nil, err
	}

	exist, err := minioClient.BucketExists(context.Background(), cfg.GetBucket())
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, fmt.Errorf("bucket %s not exist", cfg.GetBucket())
	}

	return &Minio{client: minioClient, bucket: cfg.GetBucket()}, nil
}

func (m Minio) FileInfo(ctx context.Context, userID int64, filepath string) (*domain.FileInfo, error) {
	ctx, span := otel.GetTracer().Start(ctx, "FileInfo")
	defer span.End()

	info, err := m.client.StatObject(
		ctx,
		m.bucket,
		mimioFilepath(userID, filepath),
		minio.StatObjectOptions{Checksum: true},
	)
	if err != nil {
		resp := minio.ToErrorResponse(err)
		if resp.StatusCode == http.StatusNotFound {
			return nil, domain.ErrFileNotFound
		}
		return nil, err
	}

	return &domain.FileInfo{Filepath: filepath, Size: info.Size, ContentType: info.ContentType}, nil
}

func (m Minio) UploadFile(
	ctx context.Context,
	userID int64,
	filepath string,
	object io.Reader,
	contentType string,
) (*domain.FileInfo, error) {
	ctx, span := otel.GetTracer().Start(ctx, "minio.UploadFile")
	defer span.End()

	info, err := m.client.PutObject(
		ctx, m.bucket,
		mimioFilepath(userID, filepath),
		object,
		-1,
		minio.PutObjectOptions{ContentType: contentType, SendContentMd5: true},
	)
	if err != nil {
		return nil, err
	}

	return &domain.FileInfo{Filepath: filepath, Size: info.Size, ContentType: contentType}, nil
}

func (m Minio) DownloadFile(ctx context.Context, userID int64, filepath string) (*domain.FileInfo, io.Reader, error) {
	ctx, span := otel.GetTracer().Start(ctx, "minio.DownloadFile")
	defer span.End()

	file, err := m.client.GetObject(ctx, m.bucket, mimioFilepath(userID, filepath), minio.GetObjectOptions{})
	if err != nil {
		return nil, nil, err
	}

	info, err := file.Stat()
	if err != nil {
		return nil, nil, err
	}

	return &domain.FileInfo{Filepath: filepath, Size: info.Size, ContentType: info.ContentType}, file, nil
}

func (m Minio) DeleteFile(ctx context.Context, userID int64, filepath string) error {
	ctx, span := otel.GetTracer().Start(ctx, "minio.DeleteFile")
	defer span.End()

	err := m.client.RemoveObject(ctx, m.bucket, mimioFilepath(userID, filepath), minio.RemoveObjectOptions{})
	if err != nil {
		resp := minio.ToErrorResponse(err)
		if resp.StatusCode == http.StatusNotFound {
			return domain.ErrFileNotFound
		}
		return err
	}

	return nil
}

func mimioFilepath(userID int64, filepath string) string {
	return fmt.Sprintf("/%d%s", userID, filepath)
}
