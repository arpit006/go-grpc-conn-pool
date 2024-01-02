package grpc

import (
	"errors"
)

var (
	grpcDialErr         = errors.New("go-grpc:error while grpc dial to the server")
	asyncRefreshInitErr = errors.New("go-grpc:error while initializing background connection refresh job")
	connRefreshErr      = errors.New("go-grpc:error while refreshing connection")
	cronErr             = errors.New("go-grpc:error in initializing cron job")
	clientInitErr       = errors.New("go-grpc:error while initializing client")
)

var (
	noHealthyConnAvailableErr = errors.New("go-grpc:error no healthy connection available")
	connPoolCloseErr          = errors.New("go-grpc:error connection pool is already closed")
)
