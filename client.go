package ngrpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// GrpcClient grpc客户端
type GrpcClient struct {
	conn *grpc.ClientConn // 连接
	opts ClientOptions
}

// GetConn 获取客户端连接
func (c *GrpcClient) GetConn() *grpc.ClientConn {
	return c.conn
}

// Close 关闭
func (c *GrpcClient) Close(ctx context.Context) {
	if c.conn == nil {
		c.opts.Log.Fatalf(ctx, "close %s grpc client is nil", c.opts.Address)
		return
	}
	err := c.conn.Close()
	if err != nil {
		c.opts.Log.Fatalf(ctx, "close %s grpc client: %v", c.opts.Address, err)
		return
	}
	if c.opts.discovery != nil {
		err = c.opts.discovery.Close()
		if err != nil {
			c.opts.Log.Fatalf(ctx, "close %s grpc client: %v", c.opts.Address, err)
			return
		}
	}
}

// NewGrpcClient 创建Grpc客户端
func NewGrpcClient(ctx context.Context, opts ...ClientOption) *GrpcClient {
	client := new(GrpcClient)
	client.opts = NewClientOptions(opts...)
	grpcClientOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
				Time:                2 * time.Minute,  // 每2分钟发送一次 ping
				Timeout:             10 * time.Second, // 超过10秒无响应就断开
				PermitWithoutStream: true,             // 即使没有 RPC 也发送 ping
			},
		),
	}
	if len(client.opts.StreamClientInterceptors) > 0 {
		for _, v := range client.opts.StreamClientInterceptors {
			grpcClientOptions = append(grpcClientOptions, grpc.WithStreamInterceptor(v))
		}
	}
	if len(client.opts.UnaryClientInterceptors) > 0 {
		for _, v := range client.opts.UnaryClientInterceptors {
			grpcClientOptions = append(grpcClientOptions, grpc.WithUnaryInterceptor(v))
		}
	}
	if len(client.opts.dialOptions) > 0 {
		grpcClientOptions = append(grpcClientOptions, client.opts.dialOptions...)
	}
	if client.opts.discovery != nil {
		builder, err := client.opts.discovery.Discover(client.opts.Name)
		if err != nil {
			client.opts.Log.Fatalf(ctx, "%s grpc client discovery error: %v", client.opts.Name, err)
		}
		grpcClientOptions = append(grpcClientOptions, grpc.WithResolvers(builder))
		conn, err := grpc.Dial(client.opts.discovery.Address(), grpcClientOptions...)
		if err != nil {
			client.opts.Log.Fatalf(ctx, "%s grpc client dial error: %v", client.opts.Name, err)
		}
		client.conn = conn
	} else {
		conn, err := grpc.Dial(client.opts.Address, grpcClientOptions...)
		if err != nil {
			client.opts.Log.Fatalf(ctx, "%s grpc client dial error: %v", client.opts.Name, err)
		}
		client.conn = conn
	}
	return client
}
