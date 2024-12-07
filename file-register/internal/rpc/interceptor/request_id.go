package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ReqKey int

const reqKey ReqKey = 1

func RequestID(
	ctx context.Context,
	req interface{}, _ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return handler(ctx, req)
	}

	requestID, ok := md["request_id"]
	if !ok || len(requestID) == 0 {
		return handler(ctx, req)
	}

	ctx = context.WithValue(ctx, reqKey, requestID[0])

	return handler(ctx, req)
}

func GetReqID(ctx context.Context) string {
	k := ctx.Value(reqKey)

	if reqID, ok := k.(string); ok {
		return reqID
	}

	return ""
}
