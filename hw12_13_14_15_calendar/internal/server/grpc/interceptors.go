package internalgrpc

import (
	"context"
	"fmt"
	"time"

	"github.com/pyltsin/otusGolang/hw12_13_14_15_calendar/internal/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	h, err := handler(ctx, req)
	p, ok := peer.FromContext(ctx)
	if !ok {
		logger.Log.Error("unable to log grpc request")

		return h, err
	}

	logger.Log.Info(fmt.Sprintf("ip: %s, method: %s, latency: %s",
		p.Addr.String(),
		info.FullMethod,
		time.Since(start).String()))

	return h, err
}
