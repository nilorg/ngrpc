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
	"errors"
	"net"
	"regexp"
	"strconv"

	"github.com/nilorg/pkg/consul/register"

	"github.com/nilorg/pkg/consul"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
)

var (
	portRegexp = regexp.MustCompile(`:[0-65535]{1,5}`)
	// LocalIP 本地IP
	LocalIP LocalIPFunc = consul.LocalIP
	// defaultServer 默认server
	defaultServer *Server
)

// eventFunc 事件函数
type eventFunc func()

// LocalIPFunc 获取本地IP函数
type LocalIPFunc func() string

// Server 服务端
type Server struct {
	ServiceName       string
	address           string
	tls               bool
	startBeforeEvents []eventFunc
	startAfterEvents  []eventFunc
	rpcServer         *grpc.Server
}

func (s *Server) startBefore() {
	for i := 0; i < len(s.startBeforeEvents); i++ {
		s.startBeforeEvents[i]()
	}
}

func (s *Server) startAfter() {
	for i := 0; i < len(s.startAfterEvents); i++ {
		s.startAfterEvents[i]()
	}
}

func getPort(address string) int {
	s := portRegexp.FindString(address)
	port, _ := strconv.Atoi(s[1:])
	return port
}

// RegisterConsul 注册Consul
func (s *Server) RegisterConsul(consulServerAddr string, sInfo *register.ServiceInfo) {
	if s.ServiceName == "" {
		panic(errors.New("register consul service name be empty"))
	}
	s.startBeforeEvents = append(s.startBeforeEvents, func() {
		consul.RegisterHealthServer(s.rpcServer, s.ServiceName)
	})
	s.startAfterEvents = append(s.startAfterEvents, func() {
		if sInfo == nil {
			sInfo := register.NewServiceInfo()
			sInfo.Name = s.ServiceName
			sInfo.Tags = []string{}
			sInfo.IP = LocalIP()
			sInfo.Port = getPort(s.address)
		}
		err := register.Register(consulServerAddr, sInfo)
		if err != nil {
			grpclog.Errorf("consul服务注册错误：%s", err)
		}
	})
}

// NewServer 创建服务端
func NewServer(address string, interceptor ...grpc.UnaryServerInterceptor) *Server {
	return newServer(address, nil, nil, interceptor...)
}

// NewServerTLSFromFile 创建服务端TLSFromFile
func NewServerTLSFromFile(address string, certFile, keyFile string, interceptor ...grpc.UnaryServerInterceptor) *Server {
	// TLS认证
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		grpclog.Fatalf("Failed to generate credentials %v", err)
	}
	return newServer(address, creds, nil, interceptor...)
}

// NewServerTLS 创建服务端TLS
func NewServerTLS(address string, cert *tls.Certificate, interceptor ...grpc.UnaryServerInterceptor) *Server {
	// 实例化grpc Server, 并开启TLS认证
	return newServer(address, credentials.NewServerTLSFromCert(cert), nil, interceptor...)
}

// ValidationFunc 验证方法
type ValidationFunc func(appKey, appSecret string) bool

// NewServerCustomAuthentication 创建服务端自定义服务验证
func NewServerCustomAuthentication(address string, validation ValidationFunc, interceptor ...grpc.UnaryServerInterceptor) *Server {
	return newServer(address, nil, validation, interceptor...)
}

// NewServerTLSCustomAuthentication 创建服务端TLS自定义服务验证
func NewServerTLSCustomAuthentication(address string, cert *tls.Certificate, validation ValidationFunc, interceptor ...grpc.UnaryServerInterceptor) *Server {
	return newServer(address, credentials.NewServerTLSFromCert(cert), validation, interceptor...)
}

// newServer 创建 grpc server
func newServer(address string, creds credentials.TransportCredentials, validation ValidationFunc, interceptor ...grpc.UnaryServerInterceptor) *Server {
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
	for _, v := range interceptor {
		opts = append(opts, grpc.UnaryInterceptor(v))
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
	s.startBefore()
	defer s.startAfter()
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
