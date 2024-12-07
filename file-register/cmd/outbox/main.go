package main //nolint:cyclop

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/sparhokm/slurm/file-register/internal/broker"
	"github.com/sparhokm/slurm/file-register/internal/db/repository"
	"github.com/sparhokm/slurm/file-register/internal/service/outbox"

	"github.com/sparhokm/slurm/file-register/internal/config"
	"github.com/sparhokm/slurm/file-storage/pkg/logger"
)

const (
	redisBrokerProtocol = 3
	maxErrorCount       = 3
)

func main() {
	cfg := config.MustLoad()

	log := logger.New(logger.WithMinLogLevel(cfg.Log.MinLevel))
	log = log.WithField("env", cfg.Env)

	ctx := context.Background()
	conn, err := pgxpool.New(ctx, cfg.Database.URL)
	if err != nil {
		log.WithError(err).Error("database connect")
		os.Exit(1)
	}
	defer conn.Close()

	err = conn.Ping(ctx)
	if err != nil {
		log.WithError(err).Error("database ping")
		os.Exit(1) //nolint:gocritic
	}

	rdbBroker := redis.NewClient(&redis.Options{
		Addr:     cfg.Broker.URL,
		Password: cfg.Broker.Password,
		Protocol: redisBrokerProtocol,
	})
	if err := rdbBroker.Ping(ctx).Err(); err != nil {
		log.WithError(err).Error("redis broker ping: " + cfg.Broker.URL)
		os.Exit(1)
	}
	defer func() {
		err := rdbBroker.Close()
		if err != nil {
			log.WithError(err).Error("redis close error")
		}
	}()

	eventRepo := repository.NewEventRepository(conn)
	redisBroker := broker.NewRedisBroker(rdbBroker)
	outboxService := outbox.New(eventRepo, redisBroker, log)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	shutdown := make(chan os.Signal, 1)
	go func() {
		ticker := time.NewTicker(cfg.Outbox.Duration)
		defer func() {
			ticker.Stop()
		}()

		errorTicker := time.NewTicker(time.Minute)
		defer func() {
			errorTicker.Stop()
		}()

		errCount := 0
		for {
			select {
			case <-ticker.C:
				{
					err := outboxService.Process(ctx, cfg.Outbox.BatchSize)
					if err != nil {
						log.WithError(err).Error("process error")
						time.Sleep(time.Duration(errCount) * time.Second)
						errCount++
					}
					if errCount > maxErrorCount {
						log.WithError(err).Error("many errors")
						close(shutdown)
						return
					}
				}
			case <-errorTicker.C:
				{
					errCount = 0
				}
			case <-ctx.Done():
				{
					log.Info("context done")
					close(shutdown)
					return
				}
			case <-done:
				{
					log.Info("signal stop")
					close(shutdown)
					return
				}
			}
		}
	}()

	<-shutdown
}
