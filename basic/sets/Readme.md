### 集合

+ 目前支持线程安全和非现场安全的set
```golang
set.New(true) // threadsafe 
set.New(false) // non-threadsafe
```