package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

type Dialer func(context.Context, string, ...grpc.DialOption) (*grpc.ClientConn, error)

type options struct {
	dialer         Dialer
	dialOptions    []grpc.DialOption
	poolSize       int
	maxLifeTimeout time.Duration
	stdDev         time.Duration
}
