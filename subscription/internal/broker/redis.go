package broker

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisBroker struct {
	conn *redis.Client
}

func NewRedisBroker(conn *redis.Client) *RedisBroker {
	return &RedisBroker{conn: conn}
}

func (r RedisBroker) RunRead( //nolint:gocognit
	ctx context.Context,
	topic string,
	reader func(context.Context, map[string]interface{}) error,
) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			cmd := r.conn.XRead(ctx, &redis.XReadArgs{
				Streams: []string{topic, "$"},
				Block:   time.Second,
			})

			if cmd.Err() != nil {
				if errors.Is(cmd.Err(), redis.Nil) {
					continue // return nil
				}
				return cmd.Err()
			}

			for _, v := range cmd.Val() {
				for _, m := range v.Messages {
					err := reader(ctx, m.Values)
					if err != nil {
						return err
					}
				}
			}
		}
	}
}
