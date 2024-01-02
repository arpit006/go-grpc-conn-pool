package grpc

import (
	"context"

	"google.golang.org/grpc"
)

type Client interface {
	grpc.ClientConnInterface
	Close()
}

type client struct {
	pool *clientConnPool
}

func NewClient(cfg *ClientConfig, opts ...grpc.DialOption) (Client, error) {
	if cfg.ConnectionPoolSize <= 0 {
		cfg.ConnectionPoolSize = defaultConnectionPoolSize
	}

	if cfg.ConnectionMaxLifeTime == 0 {
		cfg.ConnectionMaxLifeTime = maxDuration
	}

	pool, err := newConnPool(cfg.Target, getPoolOptions(cfg, opts)...)
	if err != nil {
		panic(clientInitErr)
	}
	return &client{
		pool: pool,
	}, nil
}

func (c client) Invoke(ctx context.Context, method string, args any, reply any, opts ...grpc.CallOption) error {
	return c.pool.Invoke(ctx, method, args, reply, opts...)
}

func (c client) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return c.pool.NewStream(ctx, desc, method, opts...)
}

func (c client) Close() { _ = c.pool.Close() }
