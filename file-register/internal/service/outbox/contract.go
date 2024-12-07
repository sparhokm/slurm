package outbox

import (
	"context"
	"time"

	"github.com/sparhokm/slurm/file-register/internal/models"
)

type EventRepo interface {
	GetLastNoneProcessed(ctx context.Context, count int64) ([]models.FileEvent, error)
	MarkProcessed(ctx context.Context, events []models.FileEvent, processedAt time.Time) error
}

type Broker interface {
	Send(ctx context.Context, event models.FileEvent) error
}
