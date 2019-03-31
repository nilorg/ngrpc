/*
 Copyright 2019 Nilorg authors.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

	 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package ngrpc

import (
	"context"
	"crypto/tls"
	"net"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
)

var (
	// defaultServer 默认server
	defaultServer *Server
)

// Server 服务端
type Server struct {
	address   string
	tls       bool
	rpcServer *grpc.Server
}

// NewServer 创建服务端
func NewServer(address string) *Server {
	return newServer(address, nil, nil)
}

// NewServerTLSFromFile 创建服务端TLSFromFile
func NewServerTLSFromFile(address string, certFile, keyFile string) *Server {
	// TLS认证
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		grpclog.Fatalf("Failed to generate credentials %v", err)
	}
	return newServer(address, creds, nil)
}

// NewServerTLS 创建服务端TLS
func NewServerTLS(address string, cert *tls.Certificate) *Server {
	// 实例化grpc Server, 并开启TLS认证
	return newServer(address, credentials.NewServerTLSFromCert(cert), nil)
}

// ValidationFunc 验证方法
type ValidationFunc func(appKey, appSecret string) bool

// NewServerCustomAuthentication 创建服务端自定义服务验证
func NewServerCustomAuthentication(address string, validation ValidationFunc) *Server {
	return newServer(address, nil, validation)
}

// NewServerTLSCustomAuthentication 创建服务端TLS自定义服务验证
func NewServerTLSCustomAuthentication(address string, cert *tls.Certificate, validation ValidationFunc) *Server {
	return newServer(address, credentials.NewServerTLSFromCert(cert), validation)
}

// newServer 创建 grpc server
func newServer(address string, creds credentials.TransportCredentials, validation ValidationFunc) *Server {
	var opts []grpc.ServerOption
	if creds != nil {
		opts = append(opts, grpc.Creds(creds))
	}
	if validation != nil {
		unaryInterceptor := grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return nil, grpc.Errorf(codes.Unauthenticated, "无令牌认证信息")
			}
			var (
				appKey, appSecret string
			)

			if v, ok := md["app_key"]; ok {
				appKey = v[0]
			}
			if v, ok := md["app_secret"]; ok {
				appSecret = v[0]
			}
			if !validation(appKey, appSecret) {
				return nil, grpc.Errorf(codes.Unauthenticated, "无效的认证信息")
			}
			return handler(ctx, req)
		})
		opts = append(opts, unaryInterceptor)
	}
	rpcServer := grpc.NewServer(opts...)
	return &Server{
		rpcServer: rpcServer,
		address:   address,
		tls:       creds != nil,
	}
}

// Start 启动
func (s *Server) Start() {
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		grpclog.Fatalf("grpc failed to listen: %v", err)
	}
	// 在gRPC服务器上注册反射服务。
	reflection.Register(s.rpcServer)
	go func() {
		if err := s.rpcServer.Serve(lis); err != nil {
			grpclog.Fatalf("grpc failed to serve: %v", err)
		}
	}()
}

// GetSrv 获取rpc server
func (s *Server) GetSrv() *grpc.Server {
	return s.rpcServer
}

// Close 关闭
func (s *Server) Close() {
	s.rpcServer.Stop()
}

// Start 启动Grpc
func Start(address string) {
	defaultServer = NewServer(address)
	defaultServer.Start()
}

// GetSrv 获取rpc server
func GetSrv() *grpc.Server {
	return defaultServer.GetSrv()
}

// Close 关闭Grpc
func Close() {
	defaultServer.Close()
}
