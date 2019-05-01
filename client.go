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
	"crypto/x509"
	"time"

	"google.golang.org/grpc/balancer/roundrobin"

	"github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

// Client grpc客户端
type Client struct {
	conn *grpc.ClientConn // 连接
	tls  bool
}

// GetConn 获取客户端连接
func (c *Client) GetConn() *grpc.ClientConn {
	return c.conn
}

// NewClientWithBalancer 创建客户端使用负载均衡
func NewClientWithBalancer(b grpc.Balancer, interceptor ...grpc.UnaryClientInterceptor) *Client {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		// 开启 grpc 中间件的重试功能
		grpc.WithUnaryInterceptor(
			grpc_retry.UnaryClientInterceptor(
				grpc_retry.WithBackoff(grpc_retry.BackoffLinear(time.Duration(1)*time.Millisecond)), // 重试间隔时间
				grpc_retry.WithMax(3), // 重试次数
				grpc_retry.WithPerRetryTimeout(time.Duration(5)*time.Millisecond), // 重试时间
				// 返回码为如下值时重试
				grpc_retry.WithCodes(codes.ResourceExhausted, codes.Unavailable, codes.DeadlineExceeded),
			),
		),
		// 负载均衡
		grpc.WithBalancer(b),
	}
	conn, err := grpc.Dial("", opts...)
	if err != nil {
		grpclog.Errorln(err)
	}
	return &Client{
		conn: conn,
		tls:  false,
	}
}

// NewClientWithBalancerName 创建客户端使用负载均衡
func NewClientWithBalancerName(ctx context.Context, target string, interceptor ...grpc.UnaryClientInterceptor) *Client {
	opts := []grpc.DialOption{
		// 使用withBlock 但是不使用超时的话会不断的重试下去。
		grpc.WithBlock(),
		grpc.WithInsecure(),
		// 负载均衡
		grpc.WithBalancerName(roundrobin.Name),
	}
	conn, err := grpc.DialContext(ctx, target, opts...)
	if err != nil {
		grpclog.Errorln(err)
	}
	return &Client{
		conn: conn,
		tls:  false,
	}
}

// NewClient 创建grpc客户端
func NewClient(serverAddress string, interceptor ...grpc.UnaryClientInterceptor) *Client {
	return newClient(serverAddress, nil, nil, interceptor...)
}

// NewClientTLSFromFile 创建grpc客户端TLSFromFile
func NewClientTLSFromFile(serverAddress string, certFile, serverNameOverride string, interceptor ...grpc.UnaryClientInterceptor) *Client {
	// TLS连接
	creds, err := credentials.NewClientTLSFromFile(certFile, serverNameOverride)
	if err != nil {
		grpclog.Fatalf("Failed to create TLS credentials %v", err)
	}
	return newClient(serverAddress, creds, nil, interceptor...)
}

// NewClientTLS 创建grpc客户端
func NewClientTLS(serverAddress string, cp *x509.CertPool, serverNameOverride string, interceptor ...grpc.UnaryClientInterceptor) *Client {
	return newClient(serverAddress, credentials.NewClientTLSFromCert(cp, serverNameOverride), nil, interceptor...)
}

// CustomCredential 自定义凭证
type CustomCredential struct {
	AppKey, AppSecret string
	Security          bool
}

// NewCustomCredential 创建自定义凭证
func NewCustomCredential(appKey, appSecret string, tls bool) *CustomCredential {
	return &CustomCredential{
		AppKey:    appKey,
		AppSecret: appSecret,
		Security:  tls,
	}
}

// GetRequestMetadata Get请求元数据
func (c CustomCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"app_key":    c.AppKey,
		"app_secret": c.AppSecret,
	}, nil
}

// RequireTransportSecurity 是否安全传输
func (c CustomCredential) RequireTransportSecurity() bool {
	return c.Security
}

// GetCustomAuthenticationParameter 获取自定义参数
type GetCustomAuthenticationParameter func() (appID, appKey string)

// NewClientCustomAuthentication 创建grpc客户端自定义服务验证
func NewClientCustomAuthentication(serverAddress string, credential credentials.PerRPCCredentials, interceptor ...grpc.UnaryClientInterceptor) *Client {
	return newClient(serverAddress, nil, credential, interceptor...)
}

// NewClientTLSCustomAuthentication 创建grpc客户端TLS自定义服务验证
func NewClientTLSCustomAuthentication(serverAddress string, cp *x509.CertPool, serverNameOverride string, credential credentials.PerRPCCredentials, interceptor ...grpc.UnaryClientInterceptor) *Client {
	return newClient(serverAddress, credentials.NewClientTLSFromCert(cp, serverNameOverride), credential, interceptor...)
}

func newClient(serverAddress string, creds credentials.TransportCredentials, credential credentials.PerRPCCredentials, interceptor ...grpc.UnaryClientInterceptor) *Client {
	var opts []grpc.DialOption
	if creds != nil && credential != nil {
		// 使用自定义认证
		opts = append(opts, grpc.WithPerRPCCredentials(credential))
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else if creds != nil {
		// TLS
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else if credential != nil {
		opts = append(opts, grpc.WithInsecure())
		// 使用自定义认证
		opts = append(opts, grpc.WithPerRPCCredentials(credential))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	for _, v := range interceptor {
		opts = append(opts, grpc.WithUnaryInterceptor(v))
	}
	conn, err := grpc.Dial(serverAddress, opts...)
	if err != nil {
		grpclog.Errorln(err)
	}
	return &Client{
		conn: conn,
		tls:  creds != nil,
	}
}

// Close 关闭
func (c *Client) Close() {
	c.conn.Close()
}
