package broker

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/sparhokm/slurm/file-register/internal/models"
)

const (
	maxLen = 1000
)

type RedisBroker struct {
	conn *redis.Client
}

func NewRedisBroker(conn *redis.Client) *RedisBroker {
	return &RedisBroker{conn: conn}
}

func (r RedisBroker) Send(ctx context.Context, event models.FileEvent) error {
	return r.conn.XAdd(ctx, &redis.XAddArgs{
		Stream: "files",
		MaxLen: maxLen,
		Approx: true,
		Values: map[string]interface{}{
			"eventID":     event.ID,
			"eventType":   event.Type,
			"time":        event.CreatedAt.Unix(),
			"fileID":      event.File.ID,
			"ownerID":     event.File.OwnerID,
			"filepath":    event.File.Filepath,
			"size":        event.File.Size,
			"contentType": event.File.ContentType,
			"requestID":   event.RequestID,
			"traceID":     event.TraceID,
			"spanID":      event.SpanID,
		},
	}).Err()
}
