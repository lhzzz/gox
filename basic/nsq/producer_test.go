package nsq

import (
	"bytes"
	"testing"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/stretchr/testify/assert"
)

func TestNewProducer(t *testing.T) {
	config := nsq.NewConfig()
	p, err := NewProducer("10.0.0.120:4150", config)
	assert.Nil(t, err)

	err = p.Publish("test-topic", bytes.NewBufferString("hello consumer").Bytes())
	assert.Nil(t, err)

	err = p.DeferredPublish("test-topic", time.Second*5, bytes.NewBufferString("hello consumer delay").Bytes())
	assert.Nil(t, err)

	np := p.NsqProduer()
	err = np.Ping()
	assert.Nil(t, err)

	p.Stop()
}
