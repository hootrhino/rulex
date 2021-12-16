# HOOK开发
HOOK 用来实现LUA脚本引解决不了的, 或者LUA引擎性能不够的情况下的扩展功能。可以理解为功能和 LUA 脚本一样, 但是是go原生开发的功能。
## 接口
```go

//
// XHook for enhancement rulex with golang
//
type XHook interface {
	Work(data string) error
	Error(error)
	Name() string
}

```