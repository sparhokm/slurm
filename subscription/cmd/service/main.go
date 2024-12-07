package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sparhokm/slurm/subscription/internal/db/repository"
	"github.com/sparhokm/slurm/subscription/internal/rpc/handlers"
	"github.com/sparhokm/slurm/subscription/internal/rpc/interceptor"
	"github.com/sparhokm/slurm/subscription/pkg/subscription_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/sparhokm/slurm/file-storage/pkg/logger"
	"github.com/sparhokm/slurm/subscription/internal/config"
)

const (
	timeoutForShutdown = 5
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

	subscriptionRepo := repository.NewSubscriptionRepository(conn)
	handler := handlers.New(subscriptionRepo, log)

	logInterceptor := interceptor.NewLoggerInterceptor(log)

	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.RequestID, logInterceptor.Logger))
	reflection.Register(gRPCServer)

	subscription_v1.RegisterSubscriptionV1Server(gRPCServer, handler)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.WithField("addr", cfg.GRPCServer.ListenAddrAndPort()).Info("Listen")
		listener, err := net.Listen("tcp", cfg.GRPCServer.ListenAddrAndPort())
		if err != nil {
			log.WithError(err).Error("listen")
			close(done)
			return
		}
		if err := gRPCServer.Serve(listener); err != nil {
			log.WithError(err).Error("listen")
			close(done)
		}
	}()

	<-done
	log.Info("Listen stopped")

	_, cancel := context.WithTimeout(ctx, timeoutForShutdown*time.Second)
	defer func() {
		cancel()
	}()

	gRPCServer.GracefulStop()
	log.Info("Shutdown completed")
}
