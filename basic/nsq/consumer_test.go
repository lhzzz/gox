package nsq

import (
	"testing"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/stretchr/testify/assert"
)

func TestNewConsumer(t *testing.T) {
	c := NewConsumer("test-topic", "ch-test")
	c.SetMap(map[string]interface{}{
		"concurrency":   4,
		"nsqd":          "10.0.0.120:4150",
		"max_in_flight": 15,
	})
	err := c.Start(nsq.HandlerFunc(func(message *nsq.Message) error {
		t.Log(string(message.Body))
		time.Sleep(3 * time.Second)
		return nil
	}))

	assert.Nil(t, err)
	time.Sleep(10 * time.Second)

	err = c.Stop()
	assert.Nil(t, err)
}

func TestNewConsumerWithConfig(t *testing.T) {
	config := nsq.NewConfig()
	config.MaxInFlight = 1000
	config.DefaultRequeueDelay = 15
	config.MaxAttempts = 10

	c := NewConsumerWithConfig("test-topic", "ch-test", config)
	c.Set("nsqd", "10.0.0.120:4150")
	c.Set("concurrency", 4)
	err := c.Start(nsq.HandlerFunc(func(message *nsq.Message) error {
		t.Log(message.ID)
		time.Sleep(5 * time.Second)
		return nil
	}))
	assert.Nil(t, err)
	time.Sleep(10 * time.Second)

	err = c.Stop()
	assert.Nil(t, err)
}
