## 标准库规范

- 最多传5个参数，需要资源的时候第一个参数永远是资源UUID
- 函数必须两个返回值：data，error

## 文档
提倡每个LUA函数都有文档，下面提供一种快速写文档的方法：
```go
//@desc:数据转发到HTTP服务器
func __RULEX_DataToHttp(
	uuid string, //@arg: HTTP UUID
	data string, //@arg: 数据
) error //@arg: 错误信息

```
其实是用 `go` 的语法来生成文档，@开头的文本会被复制进文档里面去。必须每一个字段写一行。