// 485 温湿度传感器驱动案例
// 这是个很简单的485温湿度传感器驱动, 目的是为了演示厂商如何实现自己的设备底层驱动
// 本驱动完成于2022年4月28日, 温湿度传感器资料请移步文档
// 备注：THer 的含义是 ·temperature and humidity detector· 的简写
package driver

import (
	"rulex/typex"

	"github.com/goburrow/modbus"
)

// 协议：UART：485 baud=4800 无校验 数据位1 停止位1
// 功能码为: 3（ReadHoldingRegisters）
// 站号为:1
// 寄存器：0000H保存湿度 0001H保存温度，数据一共两个寄存器，4个字节(uint16*2)
// ---------------------
// | 00H 00H | 00H 01H |
// ---------------------
// |   湿度   |   温度  |
// ---------------------
// ** 其中低位保存小数
//
type rtu485_THer_Driver struct {
	state      typex.DriverState
	client     modbus.Client
	In         *typex.InEnd
	RuleEngine typex.RuleX
}

func NewRtu485_THer_Driver(in *typex.InEnd, e typex.RuleX,
	client modbus.Client) typex.XExternalDriver {
	return &rtu485_THer_Driver{
		state:      typex.DRIVER_STOP,
		In:         in,
		RuleEngine: e,
		client:     client,
	}
}
func (rtu485 *rtu485_THer_Driver) Test() error {
	return nil
}

func (rtu485 *rtu485_THer_Driver) Init() error {
	return nil
}

func (rtu485 *rtu485_THer_Driver) Work() error {
	return nil
}

func (rtu485 *rtu485_THer_Driver) State() typex.DriverState {
	return typex.DRIVER_RUNNING
}


func (rtu485 *rtu485_THer_Driver) Read(data []byte) (int, error) {
	// Example: 0x02 0x92 0xFF 0x98
	results, err := rtu485.client.ReadHoldingRegisters(0x00, 2)
	copy(data, results)
	return len(results), err
}

func (rtu485 *rtu485_THer_Driver) Write(_ []byte) (int, error) {
	return 0, nil

}

//---------------------------------------------------
func (rtu485 *rtu485_THer_Driver) DriverDetail() *typex.DriverDetail {
	return &typex.DriverDetail{
		Name:        "RTU 485 Temperature And Humidity Detector Driver",
		Type:        "UART",
		Description: "RTU 485 Temperature And Humidity Detector Driver",
	}
}

func (rtu485 *rtu485_THer_Driver) Stop() error {
	rtu485.state = typex.DRIVER_STOP
	return nil
}
