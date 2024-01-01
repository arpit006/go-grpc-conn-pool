package errors

import "time"

type ErrMap struct {
	Message    string
	OccurredAt time.Time
}
