package client

import (
	"context"
	"encoding/base64"
	"math/rand"
	"testing"

	"singer.com/example/micro/demo"
	"singer.com/rpc/client"
)

func TestClient(t *testing.T) {
	c, err := client.NewClient(":50051")
	if err != nil {
		t.Error(err)
		return
	}
	rpcc := demo.NewDemoServiceClient(c.Conn())

	testcase := []*demo.Request{
		{Msg: "slow"},
		{Msg: "timeout"},
		{Msg: "maxSend"},
		{Msg: genRandomString(128)},
		{Msg: "error"},
		{Msg: "error-msg"},
		{Msg: "wrap-error"},
		{Msg: "wrapf-error"},
		{Msg: "panic"},
	}

	for _, tc := range testcase {
		resp, err := rpcc.Get(context.Background(), tc)
		if err != nil {
			t.Log(err)
			continue
		}
		t.Log(resp)
	}
}

func genRandomString(len int) string {
	b := make([]byte, len)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}
