# 输入资源开发
### 接口
```go
type InEnd struct {
	UUID        string          `json:"uuid"`
	State       ResourceState   `json:"state"`
	Type        InEndType       `json:"type"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Binds       map[string]Rule `json:"-"`
	Config   map[string]interface{} `json:"config"`
	Resource XResource              `json:"-"`
}

```