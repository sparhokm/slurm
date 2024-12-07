package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sparhokm/slurm/subscription/internal/otel"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/sparhokm/slurm/file-storage/pkg/logger"
	"github.com/sparhokm/slurm/subscription/internal/broker"
	"github.com/sparhokm/slurm/subscription/internal/clients/push"
	"github.com/sparhokm/slurm/subscription/internal/config"
	"github.com/sparhokm/slurm/subscription/internal/db/repository"
	"github.com/sparhokm/slurm/subscription/internal/service/notification"
)

const (
	redisBrokerProtocol = 3
	timeoutForShutdown  = 5
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

	otelShutdown, err := otel.SetupOTelSDK(ctx, cfg.ServiceName, cfg.Tracer.Version, cfg.Tracer.Endpoint, cfg.Tracer.Enable)
	if err != nil {
		return
	}
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	subRepo := repository.NewSubscriptionRepository(conn)
	redisBroker := broker.NewRedisBroker(rdbBroker)
	pusher := push.New(log)

	s := notification.New(subRepo, redisBroker, pusher, log)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		log.Info("Run event listener")
		err := s.Run(ctx)
		if err != nil {
			log.WithError(err).Error("listen")
			close(done)
			return
		}
	}()

	<-done
	cancel()
	time.Sleep(timeoutForShutdown * time.Second)

	log.Info("Listen stopped")
}
