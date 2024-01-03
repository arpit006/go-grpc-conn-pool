package grpc

import (
	"time"
)

type ErrMap struct {
	Err        error
	OccurredAt time.Time
}
