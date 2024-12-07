package notification

import (
	"context"

	"github.com/sparhokm/slurm/subscription/internal/models"
)

type SubscriptionRepo interface {
	FindUserSubscribers(ctx context.Context, filepath string) ([]models.SubscriptionShort, error)
}

type Broker interface {
	RunRead(ctx context.Context, topic string, reader func(context.Context, map[string]interface{}) error) error
}

type Pusher interface {
	Send(ctx context.Context, receiverID int64, event models.Event) error
}
