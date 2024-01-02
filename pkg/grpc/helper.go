package grpc

import (
	"github.com/arpit006/go-grpc-conn-pool/pkg/interceptors"

	"google.golang.org/grpc"
)

func getGRPCDialOptions(cfg *ClientConfig, opts []grpc.DialOption) []grpc.DialOption {
	irs := []grpc.UnaryClientInterceptor{
		interceptors.ClientTimeoutInterceptor(cfg.RequestTimeout),
	}
	// Add any other interceptor on the way

	dialOpts := []grpc.DialOption{grpc.WithChainUnaryInterceptor(irs...)}

	if opts != nil && len(opts) > 0 {
		dialOpts = append(dialOpts, opts...)
	}
	return dialOpts
}

// WithDialOptions returns a function which gets executed when Option interface apply-function is called,
// which inherently sets the options.dialOption in the options struct
func WithDialOptions(grpcDialOpts ...grpc.DialOption) Option {
	return optionFunc(func(o *options) {
		o.dialOptions = grpcDialOpts
	})
}
