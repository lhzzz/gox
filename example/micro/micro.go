package main

import (
	"google.golang.org/grpc"
	"singer.com/rpc/micro"
)

type app struct {
}

func (a *app) Regist(s *grpc.Server) {

}

func main() {
	a := app{}
	s := micro.NewService(&a, micro.Name("demo"))
	s.Run()
}
