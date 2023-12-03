# 内部消息总线
类似于Nats一样的简单Pub\Sub框架
## 示例
```go
    ebus.publish("a.b.c", "data")
    ebus.subscribe("a.b.c", func(data))
```