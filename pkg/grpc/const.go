package grpc

import (
	"time"
)

const (
	defaultConnectionPoolSize = 1
)

const (
	defaultConnMaxTimeout   = 10 * time.Minute
	defaultConnStdDeviation = 30 * time.Second
)

var (
	maxDuration = (time.Unix(1<<63-62135596801, 999999999)).Sub(time.Now())
)
