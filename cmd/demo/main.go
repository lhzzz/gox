package main

import (
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
	pkg "singer.com/internal/demo"
	"singer.com/rpc/micro"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	a := pkg.NewDemoApp()
	s := micro.NewService(&a, micro.Name("demo"),
		micro.EnableMetric("demo"),
		micro.SlowThreshold(500*time.Millisecond),
		micro.Timeout(3*time.Second),
		micro.SetMaxConnectionIdle(5*time.Minute),
		micro.SetMaxRecvMsgSize(128),
		micro.SetMaxSendMsgSize(64),
		micro.PreRunHooks(func() error {
			logrus.Info("pre running, what to do")
			return nil
		}),
		micro.PreShutdownHooks(func() error {
			logrus.Info("pre shutdown, what to finish")
			time.Sleep(5 * time.Second)
			return nil
		}),
	)
	s.Run()
}
