package grpc

import (
	"time"
)

type ClientConfig struct {
	name                        string
	target                      string
	clientID                    string
	requestTimeout              time.Duration
	connectionPoolSize          int
	connectionMaxLifeTime       time.Duration
	connectionLifeTimeDeviation time.Duration
}

type clientConfigBuilder struct {
	name            string
	target          string
	clientID        string
	requestTimeout  time.Duration
	poolSize        int
	connMaxLifetime time.Duration
	stdDev          time.Duration
	//grpcDialer      Dialer
}

func ClientConfigBuilder() *clientConfigBuilder {
	return &clientConfigBuilder{}
}

func (b *clientConfigBuilder) WithName(name string) *clientConfigBuilder {
	b.name = name
	return b
}

func (b *clientConfigBuilder) WithTarget(target string) *clientConfigBuilder {
	b.target = target
	return b
}

func (b *clientConfigBuilder) WithClientId(client string) *clientConfigBuilder {
	b.clientID = client
	return b
}

func (b *clientConfigBuilder) WithRequestTimeout(timeout time.Duration) *clientConfigBuilder {
	b.requestTimeout = timeout
	return b
}

func (b *clientConfigBuilder) WithPoolSize(size int) *clientConfigBuilder {
	b.poolSize = size
	return b
}

func (b *clientConfigBuilder) WithConnMaxLifetime(d time.Duration) *clientConfigBuilder {
	b.connMaxLifetime = d
	return b
}

func (b *clientConfigBuilder) WithStdDeviation(d time.Duration) *clientConfigBuilder {
	b.stdDev = d
	return b
}

func (b *clientConfigBuilder) Build() *ClientConfig {
	return &ClientConfig{
		name:                        GetOrDefault[string](b.name, "grpc-v2-client"),
		target:                      GetOrDefault[string](b.target, "0.0.0.0:80"),
		clientID:                    GetOrDefault[string](b.clientID, "grpc-v2-client"),
		requestTimeout:              GetOrDefault[time.Duration](b.requestTimeout, 5*time.Minute),
		connectionPoolSize:          GetOrDefault[int](b.poolSize, defaultConnectionPoolSize),
		connectionMaxLifeTime:       GetOrDefault[time.Duration](b.connMaxLifetime, defaultConnMaxTimeout),
		connectionLifeTimeDeviation: GetOrDefault[time.Duration](b.stdDev, defaultConnStdDeviation),
	}
}

func (c *ClientConfig) Name() string { return c.name }

func (c *ClientConfig) Target() string { return c.target }

func (c *ClientConfig) ClientID() string { return c.clientID }

func (c *ClientConfig) RequestTimeout() time.Duration { return c.requestTimeout }

func (c *ClientConfig) PoolSize() int { return c.connectionPoolSize }

func (c *ClientConfig) ConnMaxLifetime() time.Duration { return c.connectionMaxLifeTime }

func (c *ClientConfig) ConnLifetimeDeviation() time.Duration { return c.connectionLifeTimeDeviation }
