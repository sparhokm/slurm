package handlers

import (
	"github.com/sparhokm/slurm/file-storage/pkg/logger"
	desc "github.com/sparhokm/slurm/subscription/pkg/subscription_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrInvalidInput = status.Error(codes.InvalidArgument, "InvalidInput")

type Implementation struct {
	desc.UnsafeSubscriptionV1Server
	subscriptionRepo SubscriptionRepository
	log              logger.Logger
}

func New(subscriptionRepo SubscriptionRepository, log logger.Logger) *Implementation {
	return &Implementation{subscriptionRepo: subscriptionRepo, log: log}
}
