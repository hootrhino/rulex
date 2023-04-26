// 485 温湿度传感器驱动案例
// 这是个很简单的485温湿度传感器驱动, 目的是为了演示厂商如何实现自己的设备底层驱动
// 本驱动完成于2022年4月28日, 温湿度传感器资料请移步文档
// 备注：THer 的含义是 ·temperature and humidity detector· 的简写
package driver

import (
	"encoding/json"
	"sync"

	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
	modbus "github.com/wwhai/gomodbus"
)

// Example: 0x02 0x92 0xFF 0x98
type sensor_data struct {
	Temperature float32 `json:"t"` //系数: 0.1
	Humidity    float32 `json:"h"` //系数: 0.1
}

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
type rtu485_THer_Driver struct {
	state      typex.DriverState
	handler    *modbus.RTUClientHandler
	client     modbus.Client
	RuleEngine typex.RuleX
	Registers  []common.RegisterRW
	device     *typex.Device
	lock       sync.Mutex
}

func NewRtu485THerDriver(d *typex.Device, e typex.RuleX,
	registers []common.RegisterRW,
	handler *modbus.RTUClientHandler,
	client modbus.Client) typex.XExternalDriver {
	return &rtu485_THer_Driver{
		state:      typex.DRIVER_STOP,
		device:     d,
		RuleEngine: e,
		handler:    handler,
		client:     client,
		Registers:  registers,
		lock:       sync.Mutex{},
	}
}

func (rtu485 *rtu485_THer_Driver) Test() error {
	_, err := rtu485.client.ReadHoldingRegisters(0x00, 2)
	return err
}

func (rtu485 *rtu485_THer_Driver) Init(map[string]string) error {

	return nil
}

func (rtu485 *rtu485_THer_Driver) Work() error {
	return nil
}

func (rtu485 *rtu485_THer_Driver) State() typex.DriverState {
	return typex.DRIVER_UP
}

func (rtu485 *rtu485_THer_Driver) Read(cmd []byte, data []byte) (int, error) {
	dataMap := map[string]common.RegisterRW{}
	for _, r := range rtu485.Registers {
		rtu485.handler.SlaveId = r.SlaverId
		rtu485.lock.Lock()
		results, err := rtu485.client.ReadHoldingRegisters(0x00, 2)
		rtu485.lock.Unlock()
		if err != nil {
			return 0, err
		}
		if len(results) == 4 {
			sd := sensor_data{
				Humidity:    float32(utils.BToU16(results, 0, 2)) * 0.1,
				Temperature: float32(utils.BToU16(results, 2, 4)) * 0.1,
			}
			bytes, _ := json.Marshal(sd)
			value := common.RegisterRW{
				Tag:      r.Tag,
				Function: r.Function,
				SlaverId: r.SlaverId,
				Address:  r.Address,
				Quantity: r.Quantity,
				Value:    string(bytes),
			}
			dataMap[r.Tag] = value
			if err != nil {
				glogger.GLogger.Error(err)
			}
		}
	}
	bytes, _ := json.Marshal(dataMap)
	copy(data, bytes)
	return len(bytes), nil
}

func (rtu485 *rtu485_THer_Driver) Write(cmd []byte, _ []byte) (int, error) {
	return 0, nil

}

// ---------------------------------------------------
func (rtu485 *rtu485_THer_Driver) DriverDetail() typex.DriverDetail {
	return typex.DriverDetail{
		Name:        "Temperature And Humidity Sensor Driver",
		Type:        "UART",
		Description: "RTU 485 Temperature And Humidity Sensor Driver",
	}
}

func (rtu485 *rtu485_THer_Driver) Stop() error {
	rtu485.state = typex.DRIVER_STOP
	return nil
}
