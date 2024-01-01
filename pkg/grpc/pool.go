package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc/connectivity"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/arpit006/go-grpc-conn-pool/pkg/errors"

	"github.com/go-co-op/gocron/v2"
	"google.golang.org/grpc"
)

type ClientConnPool struct {
	target      string
	opts        *options
	currIndex   Selector
	conns       []*clientConn
	connsMu     sync.Mutex
	refreshMu   sync.Mutex
	lastDialErr atomic.Value
	_closed     uint32
}

func (pool *ClientConnPool) Invoke(ctx context.Context, method string, args any, reply any, opts ...grpc.CallOption) error {
	c, err := pool.get()
	if err != nil {
		return err
	}
	return c.conn.Invoke(ctx, method, args, reply, opts...)
}

func (pool *ClientConnPool) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	// TODO: fill the implementation of new stream over grpc connection pool
	panic("not implemented yet")
}

func (pool *ClientConnPool) init() {
	wg := &sync.WaitGroup{}
	// initialize all the grpc connections in async
	for i := 0; i < pool.opts.poolSize; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			if c, e := pool.dialConn(); e == nil {
				pool.connsMu.Lock()
				defer pool.connsMu.Unlock()
				pool.conns = append(pool.conns, c)
			}
		}(wg)
	}

	// wait till all the connections are initialized
	wg.Wait()
}

func (pool *ClientConnPool) storeLastDialErr(err error) {
	pool.lastDialErr.Store(
		errors.ErrMap{
			Message:    err.Error(),
			OccurredAt: time.Now(),
		})
}

func (pool *ClientConnPool) dialConn() (*clientConn, error) {
	conn, err := pool.opts.dialer(context.Background(), pool.target, pool.opts.dialOptions...)
	if err != nil {
		pool.storeLastDialErr(err)
		return nil, errors.GrpcDialErr
	}

	c := wrapToClientConn(conn)
	c.setDeadline(pool.connLifeTimeout())

	return c, nil
}

func (pool *ClientConnPool) connLifeTimeout() time.Duration {
	newRandomNo := rand.New(rand.NewSource(time.Now().UnixNano())).Float64() + 0.5

	return time.Duration(float64(pool.opts.maxLifeTimeout) + (newRandomNo * float64(pool.opts.stdDev)))
}

func (pool *ClientConnPool) shouldRefresh(c *clientConn) bool {
	now := time.Now()

	// check if deadline has been exceeded
	if now.Sub(c.createdAt) >= c.deadline() {
		return true
	}

	// check if connection is not in an unexpected healhty state
	if state := c.conn.GetState(); isRefreshState(state) {
		return true
	}
	return false
}

func (pool *ClientConnPool) asyncRefresh() error {
	s, err := gocron.NewScheduler()
	if err != nil {
		return err
	}

	task := func(job gocron.Job) {
		pool.refreshInBackground()
	}
	_, err = s.NewJob(
		gocron.DurationJob(30*time.Second),
		gocron.NewTask(task),
	)

	if err != nil {
		return fmt.Errorf("[%s]. error is: [%s]", errors.CronErr, err)
	}
	s.Start()

	// TODO: close down this background job
	return nil
}

func (pool *ClientConnPool) refreshConnection(c *clientConn) error {
	// check if a connection should be refreshed
	if !pool.shouldRefresh(c) {
		return nil
	}

	newConn, err := pool.opts.dialer(context.Background(), pool.target, pool.opts.dialOptions...)
	if err != nil {
		pool.storeLastDialErr(err)
		return fmt.Errorf("[%s], error is: [%s]", errors.GrpcDialErr, err)
	}
	pool.connsMu.Lock()

	// close in background
	go func(cc *grpc.ClientConn) {
		_ = cc.Close()
	}(c.conn)

	c.conn = newConn
	c.createdAt = time.Now()
	c.setDeadline(pool.connLifeTimeout())

	pool.connsMu.Unlock()
	return nil
}

func (pool *ClientConnPool) refreshInBackground() {
	if pool.conns == nil {
		return
	}

	// acquire lock
	pool.refreshMu.Lock()
	defer pool.refreshMu.Unlock()

	// get all unhealthy connections
	unhealthyConns := make([]*clientConn, 0)
	for _, connect := range pool.conns {
		c := connect
		if pool.shouldRefresh(c) {
			unhealthyConns = append(unhealthyConns, c)
		}
	}

	wg := &sync.WaitGroup{}
	// refresh all the connections in background
	for _, c := range unhealthyConns {
		wg.Add(1)
		go func(wg *sync.WaitGroup, c *clientConn) {
			defer wg.Done()
			_ = pool.refreshConnection(c)
		}(wg, c)
	}
	// release the lock
	wg.Wait()
}

func (pool *ClientConnPool) get() (*clientConn, error) {
	idx := pool.currIndex.Select(len(pool.conns))
	conn := pool.conns[idx]

	// if current connection is unhealthy, serve the RPC from next available healthy connection
	if !pool.isHealthyConn(conn) {
		if healthyConn, err := pool.getNextHealthyConn(idx, len(pool.conns)); err == nil {
			return healthyConn, nil
		}
		// since no healthy connection found. Dial in sync to create a healthy connection to serve this RPC
		err := pool.refreshConnection(conn)
		if err != nil {
			return nil, fmt.Errorf("[%s], error is: [%s]", errors.ConnRefreshErr, err)
		}
	}
	// if current connection is healthy, return this connection
	return conn, nil
}

func (pool *ClientConnPool) isHealthyConn(c *clientConn) bool {
	now := time.Now()

	if now.Sub(c.createdAt) >= c.deadline() {
		return false
	}

	if state := c.conn.GetState(); state == connectivity.Ready {
		return true
	}

	return false
}

func (pool *ClientConnPool) getNextHealthyConn(curr, max int) (*clientConn, error) {
	ptr := curr
	for {
		ptr = (ptr + 1) % max
		if ptr == curr {
			break
		}
		if pool.isHealthyConn(pool.conns[ptr]) {
			return pool.conns[ptr], nil
		}
	}
	return nil, errors.NoHealthyConnAvailableErr
}

func NewConnPool(target string) (*ClientConnPool, error) {
	p := &ClientConnPool{
		target: target,
	}

	// initialize the client connection pool
	p.init()

	// refresh connection in background
	err := p.asyncRefresh()
	if err != nil {
		return nil, fmt.Errorf("[%s]. errors is: [%s]", errors.AsyncRefreshInitErr, err)
	}
}
