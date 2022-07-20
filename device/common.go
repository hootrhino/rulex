package device

import (
	"github.com/i4de/rulex/driver"
	"github.com/i4de/rulex/typex"
)

type modBusConfig struct {
	Timeout   *int      `json:"timeout" validate:"required" title:"连接超时" info:""`
	SlaverIds []byte    `json:"slaverIds" validate:"required" title:"TCP端口" info:""`
	Config    rtuConfig `json:"config" validate:"required" title:"工作模式" info:""`
}

type rtuConfig struct {
	Uart     string       `json:"uart" validate:"required" title:"串口路径" info:"本地系统的串口路径"`
	BaudRate int          `json:"baudRate" validate:"required" title:"波特率" info:"串口通信波特率"`
	DataBits int          `json:"dataBits" validate:"required" title:"数据位" info:"串口通信数据位"`
	Parity   typex.Parity `json:"parity" validate:"required" title:"校验位" info:"串口通信校验位"`
	StopBits int          `json:"stopBits" validate:"required" title:"停止位" info:"串口通信停止位"`
}
type S1200Config struct {
	Host          string              `json:"host" validate:"required" title:"IP地址" info:""`          // 127.0.0.1
	Port          *int                `json:"port" validate:"required" title:"端口号" info:""`           // 0
	Rack          *int                `json:"rack" validate:"required" title:"架号" info:""`            // 0
	Slot          *int                `json:"slot" validate:"required" title:"槽号" info:""`            // 1
	Model         string              `json:"model" validate:"required" title:"型号" info:""`           // s7-200 s7 1500
	Timeout       *int                `json:"timeout" validate:"required" title:"连接超时时间" info:""`     // 5s
	IdleTimeout   *int                `json:"idleTimeout" validate:"required" title:"心跳超时时间" info:""` // 5s
	ReadFrequency *int                `json:"readFrequency" validate:"required" title:"采集频率" info:""` // 5s
	Blocks        []driver.S1200Block `json:"blocks" validate:"required" title:"采集配置" info:""`        // Db
}
