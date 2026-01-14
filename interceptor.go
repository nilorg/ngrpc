package ngrpc

import (
	"context"

	"google.golang.org/grpc"
)

// GrpcContextHandler ...
type GrpcContextHandler func(ctx context.Context) context.Context

// UnaryServerInterceptor ...
func UnaryServerInterceptor(f GrpcContextHandler) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx = f(ctx)
		return handler(ctx, req)
	}
}

// StreamServerInterceptor ...
func StreamServerInterceptor(f GrpcContextHandler) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		wrapped := WrapServerStream(stream)
		wrapped.WrappedContext = f(stream.Context())
		return handler(srv, wrapped)
	}
}

// WrappedServerStream is a thin wrapper around grpc.ServerStream that allows modifying context.
// Copied from github.com/grpc-ecosystem/go-grpc-middleware/v2
type WrappedServerStream struct {
	grpc.ServerStream
	// WrappedContext is the wrapper's own Context. You can assign it.
	WrappedContext context.Context
}

// Context returns the wrapper's WrappedContext, overwriting the nested grpc.ServerStream.Context()
func (w *WrappedServerStream) Context() context.Context {
	return w.WrappedContext
}

// WrapServerStream returns a ServerStream that has the ability to overwrite context.
// Copied from github.com/grpc-ecosystem/go-grpc-middleware/v2
func WrapServerStream(stream grpc.ServerStream) *WrappedServerStream {
	if existing, ok := stream.(*WrappedServerStream); ok {
		return existing
	}
	return &WrappedServerStream{ServerStream: stream, WrappedContext: stream.Context()}
}
