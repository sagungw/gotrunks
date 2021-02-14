package log

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()
		r, err := handler(ctx, req)

		WithFields(logrus.Fields{
			"method": info.FullMethod,
			"took":   time.Since(start).Milliseconds(),
		}).Info("gRPC log")

		return r, err
	}
}
