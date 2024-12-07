package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sparhokm/slurm/file-storage/internal/otel"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	"github.com/riandyrn/otelchi"
	"github.com/sparhokm/slurm/file-register/pkg/register_v1"
	"github.com/yarlson/chiprom"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/sparhokm/slurm/file-storage/internal/clients/minio"
	"github.com/sparhokm/slurm/file-storage/internal/clients/register"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sparhokm/slurm/file-storage/internal/config"
	fileStorage "github.com/sparhokm/slurm/file-storage/internal/generated/schema"
	"github.com/sparhokm/slurm/file-storage/internal/http/handlers"
	"github.com/sparhokm/slurm/file-storage/internal/http/middleware"
	"github.com/sparhokm/slurm/file-storage/pkg/logger"
)

const (
	timeoutForShutdown = 5
)

func main() {
	cfg := config.MustLoad()

	log := logger.New(logger.WithMinLogLevel(cfg.Log.MinLevel))
	log = log.WithField("env", cfg.Env)

	minioClient, err := minio.New(cfg.Minio)
	if err != nil {
		log.WithError(err).Error("minio init")
		os.Exit(1)
	}

	conn, err := grpc.NewClient(
		cfg.Register.URL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithIdleTimeout(cfg.Register.IdleTimeout),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithUnaryInterceptor(register.RequestID),
	)

	if err != nil {
		log.WithError(err).Error("register client init")
		os.Exit(1)
	}
	defer conn.Close()

	fileRegisterClient := register.New(register_v1.NewRegisterV1Client(conn))

	swagger, err := fileStorage.GetSwagger()
	if err != nil {
		log.WithError(err).Error("get swagger")
		os.Exit(1) //nolint:gocritic
	}

	handle := handlers.New(log, minioClient, fileRegisterClient)

	r := chi.NewRouter()
	r.Use(nethttpmiddleware.OapiRequestValidatorWithOptions(
		swagger, &nethttpmiddleware.Options{
			// для проверки авторизации будет использоваться отдельная middleware
			Options: openapi3filter.Options{AuthenticationFunc: openapi3filter.NoopAuthenticationFunc},
		}),
	)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)
	r.Use(chiprom.NewMiddleware(cfg.ServiceName))
	r.Use(otelchi.Middleware(cfg.ServiceName, otelchi.WithChiRoutes(r)))
	if cfg.HTTPServer.Debug {
		log.Info("http logger enable")
		r.Use(middleware.NewLogger(log))
	}

	r.Use(middleware.NewAuthentication(log))

	fileStorage.HandlerFromMux(handle, r)

	baseRouter := chi.NewRouter()
	baseRouter.Handle("/metrics", promhttp.Handler())
	baseRouter.Mount("/", r)

	s := &http.Server{
		Handler:      baseRouter,
		Addr:         cfg.HTTPServer.ListenAddrAndPort(),
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	otelShutdown, err := otel.SetupOTelSDK(
		ctx,
		cfg.ServiceName,
		cfg.Tracer.Version,
		cfg.Tracer.Endpoint,
		cfg.Tracer.Enable,
	)
	if err != nil {
		return
	}
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	go func() {
		log.WithField("addr", cfg.HTTPServer.ListenAddrAndPort()).Info("Listen start")
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.WithError(err).Error("listen")
			stop()
		}
	}()

	<-ctx.Done()
	log.Info("Listen stopped")

	ctx, cancel := context.WithTimeout(ctx, timeoutForShutdown*time.Second)
	defer func() {
		cancel()
	}()

	if err := s.Shutdown(ctx); err != nil {
		log.WithError(err).Error("Shutdown error")
		os.Exit(1)
	}
	log.Info("Shutdown completed")
}
