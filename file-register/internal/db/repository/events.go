package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sparhokm/slurm/file-register/internal/models"
)

type EventRepository struct {
	conn *pgxpool.Pool
}

func NewEventRepository(conn *pgxpool.Pool) *EventRepository {
	return &EventRepository{conn: conn}
}

func (r *EventRepository) GetLastNoneProcessed(ctx context.Context, count int64) ([]models.FileEvent, error) {
	rows, err := r.conn.Query(
		ctx,
		"SELECT id, event_type, request_id, trace_id, span_id, payload, created_at "+
			"FROM outbox WHERE processed_at is null ORDER BY id LIMIT $1",
		count,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]models.FileEvent, 0, count)
	for rows.Next() {
		event := models.FileEvent{}
		err := rows.Scan(
			&event.ID,
			&event.Type,
			&event.RequestID,
			&event.TraceID,
			&event.SpanID,
			&event.File,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		events = append(events, event)
	}

	return events, nil
}

func (r *EventRepository) MarkProcessed(
	ctx context.Context,
	events []models.FileEvent,
	processedAt time.Time,
) error {
	if len(events) == 0 {
		return errors.New("list events is empty")
	}

	params := make([]interface{}, 0, len(events)+1)
	params = append(params, processedAt)
	idPlaceholders := make([]string, 0, len(events))
	for i, e := range events {
		params = append(params, e.ID)
		idPlaceholders = append(idPlaceholders, fmt.Sprintf("$%d", i+2)) //nolint:mnd
	}

	sql := fmt.Sprintf(
		`UPDATE outbox SET processed_at = $1 WHERE id in (%s) AND processed_at is null`,
		strings.Join(idPlaceholders, ","),
	)

	result, err := r.conn.Exec(
		ctx,
		sql,
		params...,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() != int64(len(events)) {
		return errors.New("not all events mark as processed")
	}

	return nil
}
