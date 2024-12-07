package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/sparhokm/slurm/file-register/internal/otel"

	"github.com/sparhokm/slurm/file-register/internal/rpc/interceptor"

	"github.com/jackc/pgx/v5"
	"github.com/sparhokm/slurm/file-register/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FileRepository struct {
	conn *pgxpool.Pool
}

func NewFileRepository(conn *pgxpool.Pool) *FileRepository {
	return &FileRepository{conn: conn}
}

var ErrNotFound = errors.New("file not found")
var ErrMissVersion = errors.New("miss file version")

func (r *FileRepository) Create(ctx context.Context, file *models.File) error {
	ctx, span := otel.GetTracer().Start(ctx, "file.Create")
	defer span.End()

	eventPayload, err := json.Marshal(file)
	if err != nil {
		return err
	}

	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx) //nolint:errcheck
		} else {
			tx.Commit(ctx) //nolint:errcheck
		}
	}()

	_, err = tx.Exec(
		ctx,
		"INSERT INTO files (id, owner_id, content_type, size, filepath, version, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)", //nolint:lll
		file.ID, file.OwnerID, file.ContentType, file.Size, file.Filepath, 1, time.Now(),
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		ctx,
		"INSERT INTO outbox (event_type, payload, request_id, trace_id, span_id, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
		models.EventFileAdd, eventPayload, interceptor.GetReqID(ctx), span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String(), time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *FileRepository) GetByID(ctx context.Context, id string) (*models.File, error) {
	ctx, span := otel.GetTracer().Start(ctx, "file.GetByID")
	defer span.End()

	var file models.File
	sql := `SELECT id, owner_id, content_type, size, filepath, version, created_at, updated_at FROM files WHERE id = $1`

	err := r.conn.QueryRow(ctx, sql, id).Scan(
		&file.ID,
		&file.OwnerID,
		&file.ContentType,
		&file.Size,
		&file.Filepath,
		&file.Version,
		&file.CreatedAt,
		&file.UpdatedAt,
	)

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &file, nil
}

func (r *FileRepository) FindByFilepath(ctx context.Context, ownerID int64, filepath string) (*models.File, error) {
	ctx, span := otel.GetTracer().Start(ctx, "file.FindByFilepath")
	defer span.End()

	var file models.File
	sql := `SELECT id, owner_id, content_type, size, filepath, version, created_at, updated_at FROM files WHERE owner_id = $1 AND filepath = $2` //nolint:lll

	err := r.conn.QueryRow(ctx, sql, ownerID, filepath).Scan(
		&file.ID,
		&file.OwnerID,
		&file.ContentType,
		&file.Size,
		&file.Filepath,
		&file.Version,
		&file.CreatedAt,
		&file.UpdatedAt,
	)

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &file, nil
}

func (r *FileRepository) Update(ctx context.Context, file *models.File) error {
	ctx, span := otel.GetTracer().Start(ctx, "file.Update")
	defer span.End()

	curVersion := file.Version
	file.Version++

	eventPayload, err := json.Marshal(file)
	if err != nil {
		return err
	}

	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx) //nolint:errcheck
		} else {
			tx.Commit(ctx) //nolint:errcheck
		}
	}()

	result, err := tx.Exec(
		ctx,
		"UPDATE files SET size = $1, updated_at = $2, version = version + 1 WHERE id = $3 AND version = $4",
		file.Size, time.Now(), file.ID, curVersion,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrMissVersion
	}

	file.Version++

	_, err = tx.Exec(
		ctx,
		"INSERT INTO outbox (event_type, payload, request_id, trace_id, span_id, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
		models.EventFileUpdate, eventPayload, interceptor.GetReqID(ctx), span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String(), time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *FileRepository) Delete(ctx context.Context, file *models.File) error {
	ctx, span := otel.GetTracer().Start(ctx, "file.Delete")
	defer span.End()

	eventPayload, err := json.Marshal(file)
	if err != nil {
		return err
	}

	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx) //nolint:errcheck
		} else {
			tx.Commit(ctx) //nolint:errcheck
		}
	}()

	_, err = tx.Exec(
		ctx,
		"DELETE FROM files WHERE id = $1",
		file.ID,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		ctx,
		"INSERT INTO outbox (event_type, payload, request_id, trace_id, span_id, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
		models.EventFileDelete, eventPayload, interceptor.GetReqID(ctx), span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String(), time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}
