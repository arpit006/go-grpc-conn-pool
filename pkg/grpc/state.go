package grpc

import (
	"google.golang.org/grpc/connectivity"
)

type StateMap map[connectivity.State]bool

var refreshStateMap StateMap = map[connectivity.State]bool{
	connectivity.Shutdown:         true,
	connectivity.TransientFailure: true,
	connectivity.Idle:             true,
}

func isRefreshState(state connectivity.State) bool {
	if b, _ := refreshStateMap[state]; b {
		return true
	}
	return false
}
