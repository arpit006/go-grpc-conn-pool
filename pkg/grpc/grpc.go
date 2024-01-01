package grpc

import (
	"github.com/arpit006/go-grpc-conn-pool/pkg/middleware"

	"google.golang.org/grpc"
)

// NewClient creates a new grpc client with an internally managed connection pool.
func NewClient(cfg *ClientConfig) (*grpc.ClientConn, error) {

}

func initClient(cfg *ClientConfig, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(cfg.Target, getAdditionalOptions(cfg, opts)...)
}

func getAdditionalOptions(cfg *ClientConfig, opts []grpc.DialOption) []grpc.DialOption {
	// TODO: Add hystrix interceptor for circuit breaking support
	interceptors := []grpc.UnaryClientInterceptor{
		middleware.ClientTimeoutInterceptor(cfg.RequestTimeout),
	}
}
