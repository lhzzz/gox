### 错误处理

+ 优雅的错误处理方式，配合grpc interceptor实现错误日志和处理 

+ 借鉴了这篇文章 https://github.com/Mikaelemmmm/go-zero-looklook/blob/main/doc/chinese/10-%E9%94%99%E8%AF%AF%E5%A4%84%E7%90%86.md
+ 增加了将内部错误信息放到trailer的操作,经过grpc-web转发后可以由metadata中获取