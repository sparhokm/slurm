package notification

import (
	"context"
	"encoding/json"

	"github.com/sparhokm/slurm/subscription/internal/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/sparhokm/slurm/file-storage/pkg/logger"
	"github.com/sparhokm/slurm/subscription/internal/models"
)

const topic = "files"

type Service struct {
	subRepo SubscriptionRepo
	broker  Broker
	pusher  Pusher
	log     logger.Logger
}

func New(subRepo SubscriptionRepo, broker Broker, pusher Pusher, log logger.Logger) *Service {
	return &Service{subRepo: subRepo, broker: broker, pusher: pusher, log: log}
}

func (s Service) Run(ctx context.Context) error {
	return s.broker.RunRead(ctx, topic, s.process)
}

func (s Service) process(ctx context.Context, message map[string]interface{}) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	e := models.Event{}
	if err = json.Unmarshal(jsonData, &e); err != nil {
		return err
	}

	var span *trace.Span
	spanContext, err := otel.NewSpanContext(e.TraceID, e.SpanID)
	if err != nil {
		s.log.WithError(err).Error("create span context")
	} else {
		ctx = trace.ContextWithSpanContext(ctx, spanContext)
		spanCtx, newSpan := otel.GetTracer().Start(ctx, "process event", trace.WithAttributes(attribute.Int64("eventType", e.EventType)))
		ctx = spanCtx
		span = &newSpan
	}

	subs, err := s.subRepo.FindUserSubscribers(ctx, e.Filepath)
	if err != nil {
		return err
	}

	if len(subs) == 0 {
		return nil
	}

	for _, sub := range subs {
		err = s.pusher.Send(ctx, sub.UserID, e)
		if err != nil {
			return err
		}
	}

	if span != nil {
		(*span).End()
	}

	return nil
}
