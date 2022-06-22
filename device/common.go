package device

import "github.com/i4de/rulex/typex"

type modBusConfig struct {
	Timeout   *int       `json:"timeout" validate:"required" title:"连接超时" info:""`
	SlaverIds []byte     `json:"slaverIds" validate:"required" title:"TCP端口" info:""`
	Config    rtuConfig `json:"config" validate:"required" title:"工作模式" info:""`
}

type rtuConfig struct {
	Uart     string       `json:"uart" validate:"required" title:"串口路径" info:"本地系统的串口路径"`
	BaudRate int          `json:"baudRate" validate:"required" title:"波特率" info:"串口通信波特率"`
	DataBits int          `json:"dataBits" validate:"required" title:"数据位" info:"串口通信数据位"`
	Parity   typex.Parity `json:"parity" validate:"required" title:"校验位" info:"串口通信校验位"`
	StopBits int          `json:"stopBits" validate:"required" title:"停止位" info:"串口通信停止位"`
}
