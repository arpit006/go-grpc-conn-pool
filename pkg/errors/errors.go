package errors

import (
	"errors"
)

var (
	GrpcDialErr         = errors.New("go-grpc:errors while grpc dial to the server")
	AsyncRefreshInitErr = errors.New("go-grpc:errors while initializing background connection refresh job")
	ConnRefreshErr      = errors.New("go-grpc:error while refreshing connection")
	CronErr             = errors.New("go-grpc:error in initializing cron job")
	ClientInitErr       = errors.New("go-grpc:errors while initializing client")
	PoolInitErr         = errors.New("go-grpc:errors while initializing client pool")
)

var (
	NoHealthyConnAvailableErr = errors.New("go-grpc:error no healthy connection available")
)
