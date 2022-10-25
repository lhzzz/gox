package queue

import (
	"sync"
	"testing"
	"time"
)

func TestFIFO(t *testing.T) {
	q := NewBaseQueue()
	q.Add(1)
	q.Add(2)
	q.Add(3)
	q.Add(4)
	q.Add(5)

	it, shutdown := q.Get()
	t.Log(it, shutdown, q.Len())
	it, shutdown = q.Get()
	t.Log(it, shutdown, q.Len())
	it, shutdown = q.Get()
	t.Log(it, shutdown, q.Len())
	it, shutdown = q.Get()
	t.Log(it, shutdown, q.Len())
	it, shutdown = q.Get()
	t.Log(it, shutdown, q.Len())

	q.ShutDown()
}

func TestDuplicate(t *testing.T) {
	q := NewBaseQueue()
	q.Add(1)
	q.Add(1)
	q.Add(1)
	q.Add(1)
	q.Add(2)

	it, shutdown := q.Get()
	t.Log(it, shutdown, q.Len())
	it, shutdown = q.Get()
	t.Log(it, shutdown, q.Len())

	q.ShutDown()
}

func TestBasic(t *testing.T) {
	tests := []struct {
		queue         BaseQueue
		queueShutDown func(BaseQueue)
	}{
		{
			queue:         NewBaseQueue(),
			queueShutDown: BaseQueue.ShutDown,
		},
		{
			queue:         NewBaseQueue(),
			queueShutDown: BaseQueue.ShutDownWithDrain,
		},
	}
	for _, test := range tests {
		// If something is seriously wrong this test will never complete.

		// Start producers
		const producers = 10
		producerWG := sync.WaitGroup{}
		producerWG.Add(producers)
		for i := 0; i < producers; i++ {
			go func(i int) {
				defer producerWG.Done()
				for j := 0; j < 10; j++ {
					test.queue.Add(i)
					time.Sleep(time.Millisecond)
				}
			}(i)
		}

		// Start consumers
		const consumers = 10
		consumerWG := sync.WaitGroup{}
		consumerWG.Add(consumers)
		for i := 0; i < consumers; i++ {
			go func(i int) {
				defer consumerWG.Done()
				for {
					item, quit := test.queue.Get()
					if item == "added after shutdown!" {
						t.Errorf("Got an item added after shutdown.")
					}
					if quit {
						return
					}
					t.Logf("Worker %v: begin processing %v", i, item)
					time.Sleep(3 * time.Millisecond)
					t.Logf("Worker %v: done processing %v", i, item)
					test.queue.Done(item)
				}
			}(i)
		}

		producerWG.Wait()
		test.queueShutDown(test.queue)
		test.queue.Add("added after shutdown!")
		consumerWG.Wait()
		if test.queue.Len() != 0 {
			t.Errorf("Expected the queue to be empty, had: %v items", test.queue.Len())
		}
	}
}
