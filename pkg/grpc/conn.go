package grpc

import (
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
)

type clientConn struct {
	conn      *grpc.ClientConn
	createdAt time.Time
	dl        int64 // this will be atomic value
	cMu       sync.Mutex
}

func wrapToClientConn(cc *grpc.ClientConn) *clientConn {
	return &clientConn{conn: cc, createdAt: time.Now()}
}

func (c *clientConn) close() error {
	return c.conn.Close()
}

func (c *clientConn) setDeadline(d time.Duration) { atomic.StoreInt64(&c.dl, int64(d)) }

func (c *clientConn) deadline() time.Duration { return time.Duration(atomic.LoadInt64(&c.dl)) }
