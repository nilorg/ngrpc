package ngrpc

import (
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/nilorg/pkg/logger"
	"github.com/nilorg/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// GrpcServer 服务端
type GrpcServer struct {
	serviceName string
	address     string
	server      *grpc.Server
	Log         log.Logger
}

// GetSrv 获取rpc server
func (s *GrpcServer) GetSrv() *grpc.Server {
	return s.server
}

func (s *GrpcServer) register() {
	// 在gRPC服务器上注册反射服务。
	reflection.Register(s.server)
}

// Run ...
func (s *GrpcServer) Run() {
	s.register()

	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		s.Log.Errorf("%s grpc server failed to listen: %v", s.serviceName, err)
		return
	}
	if err := s.server.Serve(lis); err != nil {
		s.Log.Errorf("%s grpc server failed to serve: %v", s.serviceName, err)
	}
}

// Start 启动
func (s *GrpcServer) Start() {
	go func() {
		s.Run()
	}()
}

// Stop ...
func (s *GrpcServer) Stop() {
	if s.server == nil {
		s.Log.Warningf("stop %s grpc server is nil", s.serviceName)
		return
	}
	s.server.Stop()
}

// NewGrpcServer 创建Grpc服务端
func NewGrpcServer(name string, address string, streamServerInterceptors []grpc.StreamServerInterceptor, unaryServerInterceptors []grpc.UnaryServerInterceptor) *GrpcServer {
	var opts []grpc.ServerOption
	if len(streamServerInterceptors) > 0 {
		opts = append(opts, grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(streamServerInterceptors...)))
	}
	if len(unaryServerInterceptors) > 0 {
		opts = append(opts, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryServerInterceptors...)))
	}
	server := grpc.NewServer(opts...)
	if logger.Default() == nil {
		logger.Init()
	}
	return &GrpcServer{
		serviceName: name,
		server:      server,
		address:     address,
		Log:         logger.Default(),
	}
}
