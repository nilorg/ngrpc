package ngrpc

import (
	"fmt"
	"time"

	"github.com/nilorg/sdk/log"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// GrpcClient grpc客户端
type GrpcClient struct {
	serviceName string
	conn        *grpc.ClientConn // 连接
	Log         log.Logger
}

// GetConn 获取客户端连接
func (c *GrpcClient) GetConn() *grpc.ClientConn {
	return c.conn
}

// Close 关闭
func (c *GrpcClient) Close() {
	if c.conn == nil {
		c.Log.Warningf("close %s grpc client is nil", c.serviceName)
		return
	}
	err := c.conn.Close()
	if err != nil {
		c.Log.Errorf("close %s grpc client: %v", err)
		return
	}
}

// NewGrpcClient 创建Grpc客户端
func NewGrpcClient(serviceName string, port int, streamClientInterceptors []grpc.StreamClientInterceptor, unaryClientInterceptors []grpc.UnaryClientInterceptor) *GrpcClient {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
				Time:                10 * time.Second,
				Timeout:             100 * time.Millisecond,
				PermitWithoutStream: true,
			},
		),
	}
	if len(streamClientInterceptors) > 0 {
		for _, v := range streamClientInterceptors {
			opts = append(opts, grpc.WithStreamInterceptor(v))
		}
	}
	if len(unaryClientInterceptors) > 0 {
		for _, v := range unaryClientInterceptors {
			opts = append(opts, grpc.WithUnaryInterceptor(v))
		}
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", serviceName, port), opts...)
	if err != nil {
		logrus.Errorf("%s grpc client dial error: %v", serviceName, err)
	}
	return &GrpcClient{
		serviceName: serviceName,
		conn:        conn,
		Log:         logrus.StandardLogger(),
	}
}
