package ngrpc

import (
	"context"
	"fmt"
	"net"

	"github.com/nilorg/ngrpc/v2/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// GrpcServer 服务端
type GrpcServer struct {
	server *grpc.Server
	opts   ServerOptions
	ctx    context.Context
}

// GetSrv 获取rpc server
func (s *GrpcServer) GetSrv() *grpc.Server {
	return s.server
}

func (s *GrpcServer) register() {
	// 在gRPC服务器上注册反射服务。
	reflection.Register(s.server)
}

func (s *GrpcServer) Run() {
	s.register()
	var address string
	if s.opts.RandomPort {
		address = ":0"
	} else {
		address = s.opts.Address
	}
	lis, err := net.Listen("tcp", address)
	if err != nil {
		s.opts.Log.Fatalf(s.ctx, "%s grpc server failed to listen: %v", s.opts.Name, err)
		return
	}
	if s.opts.register != nil {
		serviceInfo := resolver.NewServiceInfo()
		serviceInfo.Name = s.opts.Name
		if s.opts.RandomPort {
			ipAddr, err := LocalIPv4()
			if err != nil {
				s.opts.Log.Fatalf(s.ctx, "LocalIPv4: %s", err)
				return
			}
			port := lis.Addr().(*net.TCPAddr).Port
			serviceInfo.Address = fmt.Sprintf("%s:%d", ipAddr, port)
		} else {
			serviceInfo.Address = address
		}
		err = s.opts.register.Register(serviceInfo)
		if err != nil {
			s.opts.Log.Fatalf(s.ctx, "%s grpc server failed to register: %v", s.opts.Name, err)
		}
	}
	if err := s.server.Serve(lis); err != nil {
		s.opts.Log.Fatalf(s.ctx, "%s grpc server failed to serve: %v", s.opts.Name, err)
	}
}

func (s *GrpcServer) Start() {
	go func() {
		s.Run()
	}()
}

func (s *GrpcServer) Stop() {
	if s.server == nil {
		s.opts.Log.Warnf(s.ctx, "stop %s grpc server is nil", s.opts.Name)
	} else {
		s.server.Stop()
	}
	if s.opts.register != nil {
		err := s.opts.register.Close()
		if err != nil {
			s.opts.Log.Errorf(s.ctx, "%s grpc server failed to unregister: %v", s.opts.Name, err)
		}
	}
}

// NewGrpcServer 创建Grpc服务端
func NewGrpcServer(ctx context.Context, opts ...ServerOption) *GrpcServer {
	server := new(GrpcServer)
	server.ctx = ctx
	server.opts = NewServerOptions(opts...)
	var grpcServerOptions []grpc.ServerOption
	if len(server.opts.StreamServerInterceptors) > 0 {
		grpcServerOptions = append(grpcServerOptions, grpc.ChainStreamInterceptor(server.opts.StreamServerInterceptors...))
	}
	if len(server.opts.UnaryServerInterceptors) > 0 {
		grpcServerOptions = append(grpcServerOptions, grpc.ChainUnaryInterceptor(server.opts.UnaryServerInterceptors...))
	}
	server.server = grpc.NewServer(grpcServerOptions...)
	return server
}
