package meta

import (
	"context"
	"sync"

	"google.golang.org/grpc"
)

// check opentracing
type serverStream struct {
	grpc.ServerStream
	Ctx context.Context
}

var generateMetadataMtx sync.Mutex

func NewServerStream(ctx context.Context, ss grpc.ServerStream) serverStream {
	return serverStream{
		ServerStream: ss,
		Ctx:          ctx,
	}
}

func (ss serverStream) Context() context.Context {
	return ss.Ctx
}

func GenerateContextTraceMetadata(ctx context.Context, method string) context.Context {
	generateMetadataMtx.Lock()
	defer generateMetadataMtx.Unlock()
	return generateMetaTraceID(ctx)
}
