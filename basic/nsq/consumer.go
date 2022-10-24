package nsq

import (
	"fmt"

	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
)

// Logger interface.
type logger interface {
	Output(int, string) error
}

type Consumer interface {
	SetMap(options map[string]interface{})
	Set(option string, value interface{})
	Start(handler nsq.Handler) error
	Stop() error
	NsqConsumer() *nsq.Consumer
}

// Consumer convenience layer.
type consumer struct {
	client      *nsq.Consumer
	config      *nsq.Config
	nsqds       []string
	nsqlookupds []string
	concurrency int
	channel     string
	topic       string
	level       nsq.LogLevel
	log         logger
	err         error
}

// NewConsumer returns a new consumer of `topic` and `channel`.
func NewConsumer(topic, channel string) Consumer {
	l, lvl := NewNSQLogrusLogger(logrus.GetLevel())
	return &consumer{
		log:         l,
		config:      nsq.NewConfig(),
		level:       lvl,
		channel:     channel,
		topic:       topic,
		concurrency: 1,
	}
}

func NewConsumerWithConfig(topic, channel string, config *nsq.Config) Consumer {
	l, lvl := NewNSQLogrusLogger(logrus.GetLevel())
	return &consumer{
		log:         l,
		config:      config,
		level:       lvl,
		channel:     channel,
		topic:       topic,
		concurrency: 1,
	}
}

// SetMap applies all `options`.
func (c *consumer) SetMap(options map[string]interface{}) {
	for k, v := range options {
		c.Set(k, v)
	}
}

// Set `option` to `value`, any error will be returned in `.Start()`.
//
// Custom options implemented:
//
//  - `topic` consumer topic
//  - `channel` consumer channel
//  - `nsqd` nsqd address
//  - `nsqds` nsqd addresses
//  - `nsqlookupd` nsqlookupd address
//  - `nsqlookupds` nsqlookupd addresses
//  - `concurrency` concurrent handlers [1]
//
//
func (c *consumer) Set(option string, value interface{}) {
	switch option {
	case "topic":
		c.topic = value.(string)
	case "channel":
		c.channel = value.(string)
	case "concurrency":
		c.concurrency = value.(int)
	case "nsqd":
		c.nsqds = []string{value.(string)}
	case "nsqlookupd":
		c.nsqlookupds = []string{value.(string)}
	case "nsqds":
		s, err := toStrings(value)
		if err != nil {
			c.err = fmt.Errorf("%q: %v", option, err)
			return
		}
		c.nsqds = s
	case "nsqlookupds":
		s, err := toStrings(value)
		if err != nil {
			c.err = fmt.Errorf("%q: %v", option, err)
			return
		}
		c.nsqlookupds = s
	default:
		err := c.config.Set(option, value)
		if err != nil {
			c.err = err
		}
	}
}

// Start consumer with `handler`.
func (c *consumer) Start(handler nsq.Handler) error {
	if c.err != nil {
		return c.err
	}

	client, err := nsq.NewConsumer(c.topic, c.channel, c.config)
	if err != nil {
		return err
	}
	c.client = client

	client.SetLogger(c.log, c.level)
	client.AddConcurrentHandlers(handler, c.concurrency)
	return c.connect()
}

// Stop and wait.
func (c *consumer) Stop() error {
	c.client.Stop()
	<-c.client.StopChan
	return nil
}

func (c *consumer) NsqConsumer() *nsq.Consumer {
	return c.client
}

// Connect to the configure nsqd(s) or nsqlookupd(s).
func (c *consumer) connect() error {
	if len(c.nsqds) == 0 && len(c.nsqlookupds) == 0 {
		return fmt.Errorf(`at least one "nsqd" or "nsqlookupd" address must be configured`)
	}

	if len(c.nsqds) > 0 {
		err := c.client.ConnectToNSQDs(c.nsqds)
		if err != nil {
			return err
		}
	}

	if len(c.nsqlookupds) > 0 {
		err := c.client.ConnectToNSQLookupds(c.nsqlookupds)
		if err != nil {
			return err
		}
	}

	return nil
}

// Returns a slice of strings or error.
//
// Primarily to allow for []interface{} from parsing configuration files.
func toStrings(v interface{}) ([]string, error) {
	switch v.(type) {
	case []string:
		return v.([]string), nil
	case []interface{}:
		var ret []string
		for _, e := range v.([]interface{}) {
			s, ok := e.(string)

			if !ok {
				return nil, fmt.Errorf("string expected, got %v", e)
			}

			ret = append(ret, s)
		}
		return ret, nil
	default:
		return nil, fmt.Errorf("strings expected")
	}
}
