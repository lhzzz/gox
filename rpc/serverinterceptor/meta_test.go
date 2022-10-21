package serverinterceptor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestUnaryGenerateMetadataInterceptor(t *testing.T) {
	UnaryGenerateMetadataInterceptor(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: "/Unary/Meta",
	},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			md, exist := metadata.FromIncomingContext(ctx)
			assert.EqualValues(t, exist, true)
			t.Log(md)
			return nil, nil
		})
}

func TestStreamGenerateMetadataInterceptor(t *testing.T) {
	StreamGenerateMetadataInterceptor(context.Background(), mockedStream{ctx: context.Background()}, &grpc.StreamServerInfo{
		FullMethod: "/Stream/Meta",
	}, func(srv interface{}, stream grpc.ServerStream) error {
		md, exist := metadata.FromIncomingContext(stream.Context())
		assert.EqualValues(t, exist, true)
		t.Log(md)
		return nil
	})
}

type mockedStream struct {
	ctx context.Context
}

func (m mockedStream) SetHeader(md metadata.MD) error {
	return nil
}

func (m mockedStream) SendHeader(md metadata.MD) error {
	return nil
}

func (m mockedStream) SetTrailer(md metadata.MD) {
}

func (m mockedStream) Context() context.Context {
	return m.ctx
}

func (m mockedStream) SendMsg(v interface{}) error {
	return nil
}

func (m mockedStream) RecvMsg(v interface{}) error {
	return nil
}
