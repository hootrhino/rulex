# 输出资源开发
### 接口
```go
type OutEnd struct {
	UUID        string        `json:"uuid"`
	State       ResourceState `json:"state"`
	Type        TargetType    `json:"type"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Config map[string]interface{} `json:"config"`
	Target XTarget                `json:"-"`
}
```