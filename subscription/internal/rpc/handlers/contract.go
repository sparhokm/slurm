package handlers

import (
	"context"

	"github.com/sparhokm/slurm/subscription/internal/models"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, s *models.Subscription) error
}
