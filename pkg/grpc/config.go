package grpc

import "time"

type ClientConfig struct {
	Name                        string
	Target                      string
	ClientID                    string
	RequestTimeout              time.Duration
	ConnectionPoolSize          int
	ConnectionMaxLifeTime       time.Duration
	ConnectionLifeTimeDeviation time.Duration
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
		Name:                        WrapWithDefault[string](b.name, "grpc-v2-client"),
		Target:                      WrapWithDefault[string](b.target, "0.0.0.0:80"),
		ClientID:                    WrapWithDefault[string](b.clientID, "grpc-v2-client"),
		RequestTimeout:              WrapWithDefault[time.Duration](b.requestTimeout, 5*time.Minute),
		ConnectionPoolSize:          WrapWithDefault[int](b.poolSize, defaultConnectionPoolSize),
		ConnectionMaxLifeTime:       WrapWithDefault[time.Duration](b.connMaxLifetime, defaultConnMaxTimeout),
		ConnectionLifeTimeDeviation: WrapWithDefault[time.Duration](b.stdDev, defaultConnStdDeviation),
	}
}

// TODO: add getters over ClientConfig
