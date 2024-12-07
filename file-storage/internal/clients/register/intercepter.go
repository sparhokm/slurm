package register

import (
	"context"

	"github.com/sparhokm/slurm/file-storage/internal/http/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func RequestID(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	existingMD, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		existingMD = metadata.New(nil) // Create empty metadata if none exists
	}
	newMD := metadata.Pairs("request_id", middleware.GetReqID(ctx))
	combinedMD := metadata.Join(existingMD, newMD)

	ctx = metadata.NewOutgoingContext(ctx, combinedMD)

	err := invoker(ctx, method, req, reply, cc, opts...)

	if err != nil {
		return err
	}

	return nil
}
