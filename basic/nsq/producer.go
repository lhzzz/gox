package nsq

import (
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
)

type Producer interface {
	Publish(topic string, body []byte) error
	DeferredPublish(topic string, delay time.Duration, body []byte) error
	Stop()
	NsqProduer() *nsq.Producer
}

type producer struct {
	server *nsq.Producer
	config *nsq.Config
	nsqds  string
}

func NewProducer(nsqdAddr string, config *nsq.Config) (Producer, error) {
	p := producer{
		config: config,
		nsqds:  nsqdAddr,
	}
	server, err := nsq.NewProducer(p.nsqds, p.config)
	if err != nil {
		return nil, err
	}
	l, lvl := NewNSQLogrusLogger(logrus.GetLevel())
	server.SetLogger(l, lvl)
	p.server = server
	return &p, nil
}

func (p *producer) Stop() {
	p.server.Stop()
}

func (p *producer) Publish(topic string, body []byte) error {
	return p.server.Publish(topic, body)
}

func (p *producer) DeferredPublish(topic string, delay time.Duration, body []byte) error {
	return p.server.DeferredPublish(topic, delay, body)
}

func (p *producer) NsqProduer() *nsq.Producer {
	return p.server
}
