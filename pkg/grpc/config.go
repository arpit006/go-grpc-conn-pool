package grpc

import "time"

type ClientConfig struct {
	Name                        string
	Target                      string
	ClientID                    string
	RequestTimeout              time.Duration
	ConnectionPoolSize          int
	ConnectionMaxLifeTime       time.Duration
	ConnectionLifeTimeDeviation time.Duration
}
