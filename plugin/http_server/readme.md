# Http Server
## 简介
HTTP Server 是 RULEX 的 WEB API 提供者，主要用来支持 Dashboard 以及部分性能监控。
## 开发指南
### UI渲染
资源接口包含了一个回调函数：
```go
Configs() []XConfig
```
其中返回值表示的是前端渲染的界面描述，下面是 XConfig:
```go
type XConfig struct {
	Field     string `json:"field"`     // 字段名
	FieldType string `json:"fieldType"` // 字段类型
	Title     string `json:"title"`     // 标题
	Info      string `json:"info"`      // 提示信息
}
```
如何返回给前端信息，看下面这个案例：
```go
//
type config struct {
	A int32   `json:"a" validate:"required" title:"a" info:"aaaa"`
	B int64   `json:"b" validate:"required" title:"b" info:"bbbb"`
	C string  `json:"c" validate:"required" title:"c" info:"cccc"`
	D float32 `json:"d" validate:"required" title:"d" info:"dddd"`
	F []int   `json:"f" validate:"required" title:"f" info:"ffff"`
}
func Test() {
	xcfg, err := httpserver.RenderConfig(config{})
	if err != nil {
		t.Fatal(err)
	} else {
		b, _ := json.Marshal(xcfg)
		t.Log(string(b))
	}
}
```

输出
```json
[
    {
        "field":"a",
        "fieldType":"int32",
        "title":"a",
        "info":"aaaa"
    },
    {
        "field":"b",
        "fieldType":"int64",
        "title":"b",
        "info":"bbbb"
    },
    {
        "field":"c",
        "fieldType":"string",
        "title":"c",
        "info":"cccc"
    },
    {
        "field":"d",
        "fieldType":"float32",
        "title":"d",
        "info":"dddd"
    },
    {
        "field":"f",
        "fieldType":"[]int",
        "title":"f",
        "info":"ffff"
    }
]
```
***注意点***：
> Tag必须包含 title 字段, info 字段是可选的，在界面上的体现是那个小【？】，起到帮助信息作用.