# ngrpc

[![Go Report Card](https://goreportcard.com/badge/github.com/nilorg/ngrpc/v2)](https://goreportcard.com/report/github.com/nilorg/ngrpc/v2)
[![GoDoc](https://godoc.org/github.com/nilorg/ngrpc/v2?status.svg)](https://godoc.org/github.com/nilorg/ngrpc/v2)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

一个功能丰富的 gRPC 服务端/客户端包装库，提供简化的 gRPC 服务开发体验。

## 功能特性

- 🚀 **简化的 gRPC 服务端/客户端创建**
- 🔧 **灵活的配置选项**
- 🌐 **内置服务发现支持**（基于 etcd）
- 🔍 **服务健康检查**
- 🛡️ **拦截器支持**
- 📝 **自定义日志接口**
- 🎯 **反射服务支持**
- ⚡ **Keep-alive 连接管理**
- 🔄 **随机端口分配**

## 安装

```bash
go get github.com/nilorg/ngrpc/v2
```

## 快速开始

### 创建 gRPC 服务端

```go
package main

import (
    "context"
    "github.com/nilorg/ngrpc/v2"
)

func main() {
    ctx := context.Background()
    
    // 创建服务端
    server := ngrpc.NewGrpcServer(ctx,
        ngrpc.WithServerName("my-service"),
        ngrpc.WithServerAddress(":8080"),
    )
    
    // 注册你的服务
    // pb.RegisterYourServiceServer(server.GetSrv(), &yourServiceImpl{})
    
    // 启动服务
    server.Run()
}
```

### 创建 gRPC 客户端

```go
package main

import (
    "context"
    "github.com/nilorg/ngrpc/v2"
)

func main() {
    ctx := context.Background()
    
    // 创建客户端
    client := ngrpc.NewGrpcClient(ctx,
        ngrpc.WithClientAddress("localhost:8080"),
    )
    defer client.Close(ctx)
    
    // 使用连接创建服务客户端
    // serviceClient := pb.NewYourServiceClient(client.GetConn())
}
```

## 配置选项

### 服务端配置

```go
server := ngrpc.NewGrpcServer(ctx,
    ngrpc.WithServerName("my-service"),           // 服务名称
    ngrpc.WithServerAddress(":8080"),             // 监听地址
    ngrpc.WithServerRandomPort(true),             // 随机端口
    ngrpc.WithServerLog(customLogger),            // 自定义日志
    ngrpc.WithServerRegistry(etcdRegistry),       // 服务注册
)
```

### 客户端配置

```go
client := ngrpc.NewGrpcClient(ctx,
    ngrpc.WithClientAddress("localhost:8080"),    // 服务地址
    ngrpc.WithClientLog(customLogger),            // 自定义日志
    ngrpc.WithClientDiscovery(etcdDiscovery),     // 服务发现
)
```

## 服务发现

ngrpc 支持基于 etcd 的服务发现：

```go
import (
    clientv3 "go.etcd.io/etcd/client/v3"
    "github.com/nilorg/ngrpc/v2/resolver"
)

// 创建 etcd 客户端
etcdClient, err := clientv3.New(clientv3.Config{
    Endpoints: []string{"localhost:2379"},
})
if err != nil {
    panic(err)
}

// 服务注册
registry := resolver.NewEtcdRegistry(etcdClient, "my-domain")
server := ngrpc.NewGrpcServer(ctx,
    ngrpc.WithServerRegistry(registry),
)

// 服务发现
discovery := resolver.NewEtcdDiscovery(etcdClient, "my-domain")
client := ngrpc.NewGrpcClient(ctx,
    ngrpc.WithClientDiscovery(discovery),
)
```

## 健康检查

项目包含标准的 gRPC 健康检查服务：

```go
import "github.com/nilorg/ngrpc/v2/health/grpc_health_v1"

// 在服务端注册健康检查服务
grpc_health_v1.RegisterHealthServer(server.GetSrv(), healthImpl)
```

## 拦截器

支持自定义拦截器：

```go
// 创建上下文处理器
contextHandler := func(ctx context.Context) context.Context {
    // 添加自定义逻辑
    return ctx
}

// 应用拦截器
server := ngrpc.NewGrpcServer(ctx,
    ngrpc.WithServerUnaryInterceptor(ngrpc.UnaryServerInterceptor(contextHandler)),
    ngrpc.WithServerStreamInterceptor(ngrpc.StreamServerInterceptor(contextHandler)),
)
```

## 自定义日志

实现 `Logger` 接口来使用自定义日志：

```go
type Logger interface {
    Fatalf(ctx context.Context, format string, args ...interface{})
}

type MyLogger struct{}

func (l *MyLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {
    // 自定义日志实现
}

server := ngrpc.NewGrpcServer(ctx,
    ngrpc.WithServerLog(&MyLogger{}),
)
```

## 依赖项

- [gRPC-Go](https://github.com/grpc/grpc-go)
- [grpc-middleware](https://github.com/grpc-ecosystem/go-grpc-middleware)
- [etcd client](https://go.etcd.io/etcd/client/v3)

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 贡献

欢迎贡献代码！请先查看贡献指南。

## 相关项目

- [nilorg](https://github.com/nilorg) - 更多开源项目
