package grpc

import "sync"

type Selector interface {
	Select(max int) int
}

type RoundRobinSelector struct {
	i  int
	mu sync.Mutex
}

func (rr *RoundRobinSelector) Select(max int) int {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	i := rr.i
	rr.i = (rr.i + 1) % max

	return i
}
