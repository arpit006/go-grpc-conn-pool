package grpc

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-co-op/gocron"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type clientConnPool struct {
	target      string
	opts        *options
	currIndex   Selector
	conns       []*clientConn
	connsMu     sync.Mutex
	refreshMu   sync.Mutex
	lastDialErr atomic.Value
	_closed     atomic.Uint32
}

func (pool *clientConnPool) Invoke(ctx context.Context, method string, args any, reply any, opts ...grpc.CallOption) error {
	c, err := pool.get()
	if err != nil {
		return err
	}
	return c.conn.Invoke(ctx, method, args, reply, opts...)
}

func (pool *clientConnPool) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	// TODO: fill the implementation of new stream over grpc connection pool
	panic("not implemented yet")
}

func newConnPool(target string, opts ...Option) (*clientConnPool, error) {
	p := &clientConnPool{
		target: target,
		// TODO: add grpc dial option support
		opts:      wrapToOptions(opts),
		currIndex: &RoundRobinSelector{mu: sync.Mutex{}},
	}

	// initialize the client connection pool
	p.init()

	// refresh connection in background
	err := p.asyncRefresh()
	if err != nil {
		return nil, fmt.Errorf("[%s]. errors is: [%s]", asyncRefreshInitErr, err)
	}

	return p, nil
}

func (pool *clientConnPool) init() {
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

func (pool *clientConnPool) storeLastDialErr(err error) {
	pool.lastDialErr.Store(
		ErrMap{
			Err:        err,
			OccurredAt: time.Now(),
		})
}

func (pool *clientConnPool) dialConn() (*clientConn, error) {
	conn, err := pool.opts.dialer(context.Background(), pool.target, pool.opts.dialOptions...)
	if err != nil {
		pool.storeLastDialErr(err)
		return nil, grpcDialErr
	}

	c := wrapToClientConn(conn)
	c.setDeadline(pool.connLifeTimeout())

	return c, nil
}

func (pool *clientConnPool) connLifeTimeout() time.Duration {
	newRandomNo := rand.New(rand.NewSource(time.Now().UnixNano())).Float64() + 0.5

	return time.Duration(float64(pool.opts.maxLifeTimeout) + (newRandomNo * float64(pool.opts.stdDev)))
}

func (pool *clientConnPool) shouldRefresh(c *clientConn) bool {
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

func (pool *clientConnPool) asyncRefresh() error {
	s := gocron.NewScheduler(time.UTC)
	task := func(job gocron.Job) {
		pool.refreshInBackground()
	}
	_, err := s.Every(30).Seconds().DoWithJobDetails(task)
	if err != nil {
		return fmt.Errorf("[%s], error is: [%s]", cronErr, err)
	}
	s.StartAsync()
	return nil
}

func (pool *clientConnPool) refreshConnection(c *clientConn) error {
	// check if a connection should be refreshed
	if !pool.shouldRefresh(c) {
		return nil
	}

	newConn, err := pool.opts.dialer(context.Background(), pool.target, pool.opts.dialOptions...)
	if err != nil {
		pool.storeLastDialErr(err)
		return fmt.Errorf("[%s], error is: [%s]", grpcDialErr, err)
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

func (pool *clientConnPool) refreshInBackground() {
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

func (pool *clientConnPool) get() (*clientConn, error) {
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
			return nil, fmt.Errorf("[%s], error is: [%s]", connRefreshErr, err)
		}
	}
	// if current connection is healthy, return this connection
	return conn, nil
}

func (pool *clientConnPool) isHealthyConn(c *clientConn) bool {
	now := time.Now()

	if now.Sub(c.createdAt) >= c.deadline() {
		return false
	}

	if state := c.conn.GetState(); state == connectivity.Ready {
		return true
	}

	return false
}

func (pool *clientConnPool) getNextHealthyConn(curr, max int) (*clientConn, error) {
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
	return nil, noHealthyConnAvailableErr
}

// Close will close all the active connections in the clientConnPool
func (pool *clientConnPool) Close() error {
	if !pool._closed.CompareAndSwap(0, 1) {
		return connPoolCloseErr
	}

	pool.connsMu.Lock()
	defer pool.connsMu.Unlock()

	wg := &sync.WaitGroup{}

	// close all the connections in the background
	for _, c := range pool.conns {
		wg.Add(1)
		go func(wg *sync.WaitGroup, c *clientConn) {
			defer wg.Done()

			_ = c.conn.Close()
		}(wg, c)
	}

	wg.Wait()

	pool.conns = nil

	return nil
}

func (pool *clientConnPool) closed() bool {
	return pool._closed.Load() == 1
}

func (pool *clientConnPool) DialErr() error {
	em := pool.lastDialErr.Load().(ErrMap)
	return em.Err
}
