# 内部消息总线
类似于Nats一样的简单Pub\Sub框架
## 示例
```go

package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hootrhino/rulex/component/eventbus"
)

func TestEventBus(t *testing.T) {

	eventbus.InitEventBus()
	eventbus.Subscribe("hello", &eventbus.Subscriber{
		Callback: func(Topic string, Msg eventbus.EventMessage) {
			t.Log("hello:", Msg)
		},
	})
	start := time.Now()
	for i := 0; i < 100; i++ {
		eventbus.Publish("hello", eventbus.EventMessage{
			Payload: fmt.Sprintf("world:%d", i),
		})
	}
	duration := time.Since(start)
	t.Log("time.Since(start):", duration)
	time.Sleep(3 * time.Second)
	eventbus.Flush()
}

```