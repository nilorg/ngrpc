package ngrpc

import (
	"os"

	"github.com/nilorg/ngrpc/v2/resolver"
	"google.golang.org/grpc"
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
