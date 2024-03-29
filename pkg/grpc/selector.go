package grpc

import (
	"sync"
)

type Selector interface {
	Select(max int) int
}

// RoundRobinSelector implements Selector interface to decide on how to
// decide on the iteration logic over list of connections
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
