package ngrpc

import (
	"os"
	"time"

	"github.com/nilorg/ngrpc/v2/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// ServerOptions 可选参数列表
type ServerOptions struct {
	Name                     string
	Address                  string
	Log                      Logger
	StreamServerInterceptors []grpc.StreamServerInterceptor
	UnaryServerInterceptors  []grpc.UnaryServerInterceptor
	register                 resolver.Registry
	RandomPort               bool
	KeepaliveParams          keepalive.ServerParameters
	KeepaliveEnforcement     keepalive.EnforcementPolicy
}

// ServerOption 为可选参数赋值的函数
type ServerOption func(*ServerOptions)

// NewServerOptions 创建可选参数
func NewServerOptions(opts ...ServerOption) ServerOptions {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	opt := ServerOptions{
		Name:       hostname,
		Address:    ":5000",
		RandomPort: false,
		Log:        new(StdLogger),
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func WithServerName(name string) ServerOption {
	return func(o *ServerOptions) {
		o.Name = name
	}
}

func WithServerAddress(address string) ServerOption {
	return func(o *ServerOptions) {
		o.Address = address
	}
}

func WithServerLogger(log Logger) ServerOption {
	return func(o *ServerOptions) {
		o.Log = log
	}
}

func WithServerStreamServerInterceptors(streamServerInterceptors ...grpc.StreamServerInterceptor) ServerOption {
	return func(o *ServerOptions) {
		o.StreamServerInterceptors = streamServerInterceptors
	}
}

func WithServerUnaryServerInterceptors(unaryServerInterceptors ...grpc.UnaryServerInterceptor) ServerOption {
	return func(o *ServerOptions) {
		o.UnaryServerInterceptors = unaryServerInterceptors
	}
}

func WithServerRegister(register resolver.Registry) ServerOption {
	return func(o *ServerOptions) {
		o.register = register
	}
}

func WithServerRandomPort(randomPort bool) ServerOption {
	return func(o *ServerOptions) {
		o.RandomPort = randomPort
	}
}

// WithServerKeepaliveParams 设置 gRPC keepalive 服务端参数
// time: 客户端 ping 的最小间隔，默认 5 分钟
// timeout: 服务端等待 ping 响应的超时，默认 20 秒
func WithServerKeepaliveParams(time, timeout time.Duration) ServerOption {
	return func(o *ServerOptions) {
		o.KeepaliveParams = keepalive.ServerParameters{
			Time:    time,
			Timeout: timeout,
		}
	}
}

// WithServerKeepaliveEnforcement 设置 gRPC keepalive 服务端强制策略
// minTime: 允许客户端 ping 的最小间隔，默认 5 分钟
// permitWithoutStream: 是否允许无活跃流时 ping
func WithServerKeepaliveEnforcement(minTime time.Duration, permitWithoutStream bool) ServerOption {
	return func(o *ServerOptions) {
		o.KeepaliveEnforcement = keepalive.EnforcementPolicy{
			MinTime:             minTime,
			PermitWithoutStream: permitWithoutStream,
		}
	}
}

// ClientOptions 可选参数列表
type ClientOptions struct {
	Name                     string
	Address                  string
	Log                      Logger
	StreamClientInterceptors []grpc.StreamClientInterceptor
	UnaryClientInterceptors  []grpc.UnaryClientInterceptor
	discovery                resolver.Discovery
	dialOptions              []grpc.DialOption
}

// ClientOption 为可选参数赋值的函数
type ClientOption func(*ClientOptions)

// NewClientOptions 创建可选参数
func NewClientOptions(opts ...ClientOption) ClientOptions {
	opt := ClientOptions{
		Name:    "unknown",
		Address: ":5000",
		Log:     new(StdLogger),
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func WithClientName(name string) ClientOption {
	return func(o *ClientOptions) {
		o.Name = name
	}
}

func WithClientAddress(address string) ClientOption {
	return func(o *ClientOptions) {
		o.Address = address
	}
}

func WithClientLogger(log Logger) ClientOption {
	return func(o *ClientOptions) {
		o.Log = log
	}
}

func WithClientStreamClientInterceptors(streamClientInterceptors ...grpc.StreamClientInterceptor) ClientOption {
	return func(o *ClientOptions) {
		o.StreamClientInterceptors = streamClientInterceptors
	}
}

func WithClientUnaryClientInterceptors(unaryClientInterceptors ...grpc.UnaryClientInterceptor) ClientOption {
	return func(o *ClientOptions) {
		o.UnaryClientInterceptors = unaryClientInterceptors
	}
}

func WithClientDiscovery(discovery resolver.Discovery) ClientOption {
	return func(o *ClientOptions) {
		o.discovery = discovery
	}
}

func WithClientDialOptions(dialOptions ...grpc.DialOption) ClientOption {
	return func(o *ClientOptions) {
		o.dialOptions = dialOptions
	}
}
