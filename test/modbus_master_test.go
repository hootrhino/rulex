package test

import (
	"encoding/json"
	"testing"
)

type registerParam struct {
	Tag      string `json:"tag" validate:"required"`      // Function
	Function int    `json:"function" validate:"required"` // Function
	Address  uint16 `json:"address" validate:"required"`  // Address
	Quantity uint16 `json:"quantity" validate:"required"` // Quantity
}

type ModBusConfig struct {
	Mode     string `json:"mode" title:"工作模式" info:"可以在 UART/TCP 两个模式之间切换"`
	Timeout  int    `json:"timeout" validate:"required" title:"连接超时"`
	SlaverId byte   `json:"slaverId" validate:"required" title:"TCP端口"`
	//
	// Weather allow AutoRequest?
	AutoRequest bool `json:"autoRequest" title:"启动轮询"`
	// Request Frequency, default 5 second
	Frequency      int64           `json:"frequency" validate:"required" title:"采集频率"`
	Config         interface{}     `json:"config" validate:"required" title:"工作模式配置"`
	RegisterParams []registerParam `json:"registerParams" validate:"required" title:"寄存器配置"`
}

type RtuConfig struct {
	Uart     string `json:"uart" validate:"required" title:"串口路径" info:"本地系统的串口路径"`
	BaudRate int    `json:"baudRate" validate:"required" title:"波特率" info:"串口通信波特率"`
	DataBits int    `json:"dataBits" validate:"required" title:"数据位" info:"串口通信数据位"`
	Parity   string `json:"parity" validate:"required" title:"校验位" info:"串口通信校验位"`
	StopBits int    `json:"stopBits" validate:"required" title:"停止位" info:"串口通信停止位"`
}

type TcpConfig struct {
	Ip   string `json:"ip" validate:"required" title:"IP地址"`
	Port int    `json:"port" validate:"required" title:"端口"`
}

func TestMaster(t *testing.T) {
	d1 := ModBusConfig{
		Mode:      "TCP",
		Timeout:   5,
		SlaverId:  1,
		Frequency: 5,
		Config: TcpConfig{
			Ip:   "127.0.0.1",
			Port: 1254,
		},
		RegisterParams: []registerParam{
			{
				Tag:      "A",
				Function: 3,
				Address:  0,
				Quantity: 10,
			},
		},
	}
	d2 := ModBusConfig{
		Mode:      "UART",
		Timeout:   5,
		SlaverId:  1,
		Frequency: 5,
		Config: RtuConfig{
			Uart:     "com1",
			BaudRate: 115200,
			DataBits: 8,
			Parity:   "N",
			StopBits: 1,
		},
		RegisterParams: []registerParam{
			{
				Tag:      "A",
				Function: 3,
				Address:  0,
				Quantity: 10,
			},
		},
	}
	b1, _ := json.MarshalIndent(d1, "", " ")
	b2, _ := json.MarshalIndent(d2, "", " ")
	t.Log(string(b1))
	t.Log(string(b2))
}
