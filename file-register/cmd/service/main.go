package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sparhokm/slurm/file-register/internal/otel"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sparhokm/slurm/file-register/internal/db/repository"
	"github.com/sparhokm/slurm/file-register/internal/rpc/handlers"
	"github.com/sparhokm/slurm/file-register/internal/rpc/interceptor"
	"github.com/sparhokm/slurm/file-register/pkg/register_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/sparhokm/slurm/file-register/internal/config"
	"github.com/sparhokm/slurm/file-storage/pkg/logger"
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

	fileRepo := repository.NewFileRepository(conn)
	handler := handlers.New(fileRepo, log)

	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.05, 0.1, 0.5, 1, 2, 3, 5}),
		),
	)
	prometheus.MustRegister(srvMetrics)

	logInterceptor := interceptor.NewLoggerInterceptor(log)

	gRPCServer := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()), grpc.ChainUnaryInterceptor(interceptor.RequestID, logInterceptor.Logger))
	reflection.Register(gRPCServer)

	register_v1.RegisterRegisterV1Server(gRPCServer, handler)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	http.Handle("/metrics", promhttp.Handler())

	otelShutdown, err := otel.SetupOTelSDK(ctx, cfg.ServiceName, cfg.Tracer.Version, cfg.Tracer.Endpoint, cfg.Tracer.Enable)
	if err != nil {
		return
	}
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

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

	go func() {
		err := http.ListenAndServe(cfg.Prom.ListenAddrAndPort(), nil) //nolint:gosec
		if err != nil {
			log.WithError(err).Error("prom listen")
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
