package meta

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"singer.com/basic/trace"
)

const (
	kMetadataKeyRequestID string = "requestid"
)

func GetRequestId(ctx context.Context) (requestId string) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return
	}
	if len(md.Get(kMetadataKeyRequestID)) > 0 {
		return md.Get(kMetadataKeyRequestID)[0]
	}
	return
}

func newTraceId(ctx context.Context) string {
	traceID := trace.TraceIDFromContext(ctx)
	if !traceID.IsValid() {
		return uuid.New().String()
	}
	return traceID.String()
}

func keyExistMD(md metadata.MD, key string) bool {
	return len(md.Get(key)) > 0
}

func generateMetaTraceID(ctx context.Context) context.Context {
	var traceId string
	md, exist := metadata.FromIncomingContext(ctx)
	if !exist {
		md = metadata.New(nil)
	}
	if !keyExistMD(md, kMetadataKeyRequestID) {
		traceId = newTraceId(ctx)
		md.Set(kMetadataKeyRequestID, traceId)
	} else {
		traceId = md.Get(kMetadataKeyRequestID)[0]
	}
	ctx = metadata.NewIncomingContext(ctx, md)
	grpc.SetTrailer(ctx, metadata.Pairs(kMetadataKeyRequestID, traceId))
	return ctx
}
