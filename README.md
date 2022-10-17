# gox
一个基于grpc的rpc微服务框架，包含以下功能

### 服务治理
+ 熔断
+ 限流
+ 重试
+ 超时
+ 分布式追踪
+ 元数据传递
+ 慢日志
+ 异常捕获
+ 错误处理


### 基础库
+ 常用数据结构 
    + 各种queue (优先队列、延时队列等)
    + set
+ orm 数据库
+ redis 相关数据结构和功能
    + bloom
    + geo
    + hyperloglog
    + 分布式锁
+ nsq
+ jaeger 分布式调用追踪
+ prometheus 监控与统计
+ pprof 性能分析

### Feature List
+ 负载均衡 (目前都是依赖k8s的service实现负载均衡)
