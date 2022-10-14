package meta

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
	"singer.com/basic/trace"
)

type messageHeadKey struct{}

type MessageHead struct {
	InterfaceName string
	RequestID     string
}

const (
	kMetadataKeyRequestID string = "requestid"
)

func MessageHeadContext(ctx context.Context) MessageHead {
	raw, ok := FromMessageHeadContext(ctx)
	if ok && raw != nil {
		return *raw
	}
	return MessageHead{
		RequestID: newTraceId(ctx),
	}
}

// FromMessageHeadContext
// 不保证一定会有值返回，所以需要判断值是否存在
func FromMessageHeadContext(ctx context.Context) (*MessageHead, bool) {
	raw, ok := ctx.Value(messageHeadKey{}).(*MessageHead)
	return raw, ok
}

func newTraceId(ctx context.Context) string {
	traceID := trace.TraceIDFromContext(ctx)
	if !traceID.IsValid() {
		return uuid.New().String()
	}
	return traceID.String()
}

func newMessageHeadContext(ctx context.Context, mh *MessageHead) context.Context {
	return context.WithValue(ctx, messageHeadKey{}, mh)
}

func newMessageHead(interfaceName, requestID string) *MessageHead {
	return &MessageHead{
		RequestID:     requestID,
		InterfaceName: interfaceName,
	}
}

func keyExistMD(md metadata.MD, key string) bool {
	return len(md.Get(key)) > 0
}

func generateMetaTraceID(ctx context.Context) (context.Context, string) {
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
	return ctx, traceId
}
