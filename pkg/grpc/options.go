package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

type Option interface {
	apply(*options)
}

type Dialer func(context.Context, string, ...grpc.DialOption) (*grpc.ClientConn, error)

// To override grpc dialer for connection
func (d Dialer) apply(o *options) { o.dialer = d }

type options struct {
	dialer            Dialer
	withBlockGrpcDial bool
	dialOptions       []grpc.DialOption
	poolSize          int
	maxLifeTimeout    time.Duration
	stdDev            time.Duration
}

type ConnectionMaxLifeTime time.Duration

func (cml ConnectionMaxLifeTime) apply(o *options) { o.maxLifeTimeout = time.Duration(cml) }

type ConnectionStandardDeviation time.Duration

func (d ConnectionStandardDeviation) apply(o *options) { o.stdDev = time.Duration(d) }

type PoolSize int

func (s PoolSize) apply(o *options) { o.poolSize = int(s) }

type WithBlockGrpcDial bool

// TODO: if blocking call, add this to options
func (block WithBlockGrpcDial) apply(o *options) {
	o.withBlockGrpcDial = bool(block)
	if o.dialOptions == nil {
		o.dialOptions = make([]grpc.DialOption, 0)
	}
	o.dialOptions = append(o.dialOptions, grpc.WithBlock())
}

// optionFunc is a helper function which appends all the grpc dialOptions to options in the list
type optionFunc func(*options)

func (f optionFunc) apply(opts *options) { f(opts) }

func wrapToOptions(opts []Option) *options {
	opt := &options{
		dialer:            grpc.DialContext,
		withBlockGrpcDial: false,
		poolSize:          defaultConnectionPoolSize,
		maxLifeTimeout:    defaultConnMaxTimeout,
		stdDev:            defaultConnStdDeviation,
	}

	for _, o := range opts {
		o.apply(opt)
	}

	return opt
}

func getPoolOptions(cfg *ClientConfig, opts []grpc.DialOption) []Option {
	return []Option{
		WithDialOptions(getGRPCDialOptions(cfg, opts)...),
		PoolSize(cfg.ConnectionPoolSize),
		ConnectionMaxLifeTime(cfg.ConnectionMaxLifeTime),
		ConnectionStandardDeviation(cfg.ConnectionLifeTimeDeviation),
		// TODO: add dialer override support in ClientConfig
	}
}
