package serverinterceptor

import (
	"context"

	"google.golang.org/grpc"
	"singer.com/basic/meta"
)

// GenerateMetadataInterceptor 对于没有requestID的请求自动生成requestID
// 并会生成messageHead，用来做消息日志的统一头部，通过meta.FromMessageHeadContext(ctx)获取
func UnaryGenerateMetadataInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ctx = meta.GenerateContextTraceMetadata(ctx, info.FullMethod)
	return handler(ctx, req)
}

// GenerateMetadataStreamInterceptor 对于没有requestID的请求自动生成requestID
// 并会生成messageHead，用来做消息日志的统一头部，通过meta.FromMessageHeadContext(ctx)获取
func StreamGenerateMetadataInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := meta.GenerateContextTraceMetadata(ss.Context(), info.FullMethod)
	return handler(srv, meta.NewServerStream(ctx, ss))
}
