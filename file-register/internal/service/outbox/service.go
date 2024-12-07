package outbox

import (
	"context"
	"time"

	"github.com/sparhokm/slurm/file-storage/pkg/logger"
)

type Outbox struct {
	eventRepo EventRepo
	broker    Broker
	log       logger.Logger
}

func New(eventRepo EventRepo, broker Broker, log logger.Logger) *Outbox {
	return &Outbox{eventRepo: eventRepo, broker: broker, log: log}
}

func (s Outbox) Process(ctx context.Context, count int64) error {
	events, err := s.eventRepo.GetLastNoneProcessed(ctx, count)
	if err != nil {
		return err
	}

	if len(events) == 0 {
		return nil
	}

	for _, event := range events {
		s.log.WithFields(map[string]any{
			"filepath": event.File.Filepath,
			"ownerID":  event.File.OwnerID,
		}).Debug("send message")

		err = s.broker.Send(ctx, event)
		if err != nil {
			return err
		}
	}

	err = s.eventRepo.MarkProcessed(ctx, events, time.Now())
	if err != nil {
		return err
	}

	return nil
}
