package interceptor

import (
	"context"
	"time"

	"github.com/sparhokm/slurm/file-storage/pkg/logger"

	"google.golang.org/grpc"
)

type LoggerInterceptor struct {
	log logger.Logger
}

func NewLoggerInterceptor(log logger.Logger) *LoggerInterceptor {
	return &LoggerInterceptor{log}
}

func (l LoggerInterceptor) Logger(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	now := time.Now()

	res, err := handler(ctx, req)
	if err != nil {
		l.log.WithError(err).WithFields(map[string]any{
			"method":     info.FullMethod,
			"request_id": GetReqID(ctx),
		})
	}

	l.log.WithFields(map[string]any{
		"method":     info.FullMethod,
		"request_id": GetReqID(ctx),
		"duration":   time.Since(now),
	}).Debug("request")

	return res, err
}
