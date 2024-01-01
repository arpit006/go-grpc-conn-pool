package grpc

import (
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
)

type clientConn struct {
	conn      *grpc.ClientConn
	createdAt time.Time
	dl        atomic.Int64 // this will be atomic value
}

func wrapToClientConn(cc *grpc.ClientConn) *clientConn {
	return &clientConn{conn: cc, createdAt: time.Now()}
}

func (c *clientConn) close() error {
	return c.conn.Close()
}

func (c *clientConn) setDeadline(d time.Duration) { c.dl.Store(int64(d)) }

func (c *clientConn) deadline() time.Duration { return time.Duration(c.dl.Load()) }
