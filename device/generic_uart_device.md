# 通用串口收发器
## 配置
```go
type GenericUartConfig struct {
	Tag      string `json:"tag" validate:"required" title:"数据Tag" info:"给数据打标签"`
	Uart     string `json:"uart" validate:"required" title:"串口路径" info:"本地系统的串口路径"`
	BaudRate int    `json:"baudRate" validate:"required" title:"波特率" info:"串口通信波特率"`
	DataBits int    `json:"dataBits" validate:"required" title:"数据位" info:"串口通信数据位"`
	// 结束符, 默认是 '\n'；但是可以自己定义
	Decollator string `json:"decollator" title:"协议分隔符"`
	// Weather allow AutoRequest?
	AutoRequest bool `json:"autoRequest" title:"启动轮询"`
	// Request Frequency, default 5 second
	Frequency int64  `json:"frequency" validate:"required" title:"采集频率"`
	Timeout   int    `json:"timeout" validate:"required" title:"连接超时"`
	Parity    string `json:"parity" validate:"required" title:"奇偶校验" info:"奇偶校验"`
	StopBits  int    `json:"stopBits" validate:"required" title:"停止位" info:"串口通信停止位"`
}
```

## 注意
需要注意协议分隔符，默认是 '\n'；但是可以自己定义。