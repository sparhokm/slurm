package push

import (
	"context"
	"time"

	"github.com/sparhokm/slurm/file-storage/pkg/logger"
	"github.com/sparhokm/slurm/subscription/internal/models"
)

type Client struct {
	log logger.Logger
}

func New(log logger.Logger) *Client {
	return &Client{log: log}
}

func (c Client) Send(_ context.Context, receiverID int64, event models.Event) error {
	c.log.WithFields(map[string]any{
		"receiverID": receiverID,
		"event":      event,
	}).Info("send push")
	time.Sleep(100 * time.Millisecond)
	return nil
}
