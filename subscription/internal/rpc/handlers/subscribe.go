package handlers

import (
	"context"
	"time"

	"github.com/sparhokm/slurm/subscription/internal/models"
	"github.com/sparhokm/slurm/subscription/internal/rpc/interceptor"
	desc "github.com/sparhokm/slurm/subscription/pkg/subscription_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i Implementation) Subscribe(ctx context.Context, subIn *desc.SubscribeIn) (*emptypb.Empty, error) {
	log := i.log.WithFields(map[string]any{
		"method":     "subscribe",
		"request_id": interceptor.GetReqID(ctx),
	})

	if len(subIn.GetPrefix()) == 0 || subIn.GetUserID() == 0 ||
		(subIn.FilesOwnerID != nil && subIn.GetFilesOwnerID() == 0) {
		return nil, ErrInvalidInput
	}

	sub := models.Subscription{
		UserID:       subIn.GetUserID(),
		Prefix:       subIn.GetPrefix(),
		FilesOwnerID: subIn.FilesOwnerID, //nolint:protogetter
		CreatedAt:    time.Now(),
	}

	err := i.subscriptionRepo.Create(ctx, &sub)
	if err != nil {
		log.WithError(err).Error("repo error")
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
