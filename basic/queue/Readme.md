## 各种队列

### BaseQueue 
+ 有序 FIFO
+ 去重: 相同元素在同一时间不会被重复处理，例如一个元素在处理之前被添加多次，它只会被处理一次
+ 并发性： 支持多生产多消费者

### DelayingQueue
+ 延时队列：基于BaseQueue, 延迟一段时间后再将元素放到队列中

### RateLimitingQueue
+ 限速队列：基于DelayingQueue, 支持元素存入队列时进行速率限制
