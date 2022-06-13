package driver

import (
	"encoding/json"
	"rulex/typex"

	"github.com/goburrow/modbus"
)

type YK8RelayControllerDriver struct {
	state      typex.DriverState
	client     modbus.Client
	device     *typex.Device
	RuleEngine typex.RuleX
}

func NewYK8RelayControllerDriver(d *typex.Device, e typex.RuleX,
	client modbus.Client) typex.XExternalDriver {
	return &YK8RelayControllerDriver{
		state:      typex.DRIVER_STOP,
		device:     d,
		RuleEngine: e,
		client:     client,
	}
}

func (rtu485 *YK8RelayControllerDriver) Test() error {
	return nil
}

func (rtu485 *YK8RelayControllerDriver) Init(map[string]string) error {
	return nil
}

func (rtu485 *YK8RelayControllerDriver) Work() error {
	return nil
}

func (rtu485 *YK8RelayControllerDriver) State() typex.DriverState {
	return typex.DRIVER_RUNNING
}

//
type yk08sw struct {
	Sw1 bool `json:"sw1"`
	Sw2 bool `json:"sw2"`
	Sw3 bool `json:"sw3"`
	Sw4 bool `json:"sw4"`
	Sw5 bool `json:"sw5"`
	Sw6 bool `json:"sw6"`
	Sw7 bool `json:"sw7"`
	Sw8 bool `json:"sw8"`
}

// １、命令格式：地址（01-FE，1 字节）、功能码（01，1 字节）、继电器起始地址（0000，2 字节）、继电器
// 数量（0010，2 字节）、校验位（2 字节，低位先行）
// ２、回复格式：地址（01-FE，1 字节）、功能码（01，1 字节）、状态字节数（02，1 字节）、继电器状态（0000，2
// 字节）、校验位（2 字节，低位先行）
func getABitOnByte(b byte, position uint8) (v uint8) {
	//  --------------->
	//  7 6 5 4 3 2 1 0
	// |.|.|.|.|.|.|.|.|
	//
	mask := 0b00000001
	if position == 0 {
		return (b & byte(mask)) >> position
	} else {
		return (b & (1 << mask)) >> position
	}
}
func byteToBool(data byte, index int) bool {
	return getABitOnByte(data, 0) == 0

}
func (rtu485 *YK8RelayControllerDriver) Read(data []byte) (int, error) {
	results, err := rtu485.client.ReadCoils(0x00, 1)
	if err != nil {
		return 0, err
	}
	if len(results) == 1 {
		yks := yk08sw{
			Sw1: byteToBool(results[0], 0),
			Sw2: byteToBool(results[0], 1),
			Sw3: byteToBool(results[0], 2),
			Sw4: byteToBool(results[0], 3),
			Sw5: byteToBool(results[0], 4),
			Sw6: byteToBool(results[0], 5),
			Sw7: byteToBool(results[0], 6),
			Sw8: byteToBool(results[0], 7),
		}
		bytes, _ := json.Marshal(yks)
		copy(data, bytes)
	}
	return len(data), err
}

//
// 1、命令格式：地址（01-FE，1 字节）、功能码（0F，1 字节）、继电器起始地址（0000，2 字节）、继电器
// 数量（0010，2 字节）、写入数据字节（02，1 字节）、写入字节（0000，2 字节）、校验位（2 字节，低位先
// 行）
// ２、回复格式：地址（01-FE，1 字节）、功能码（01，1 字节）、继电器起始地址（0000，2 字节）、继电器
// 数量（0010，2 字节）、校验位（2 字节，低位先行）
//
func (rtu485 *YK8RelayControllerDriver) Write(data []byte) (int, error) {
	// 0x10： 16 个字节 8 个寄存器
	results, err := rtu485.client.WriteMultipleCoils(0, 1, data)
	if err != nil {
		return 0, err
	}
	copy(data, results)
	return 0, err
}

//---------------------------------------------------
func (rtu485 *YK8RelayControllerDriver) DriverDetail() *typex.DriverDetail {
	return &typex.DriverDetail{
		Name:        "Temperature And Humidity Sensor Driver",
		Type:        "UART",
		Description: "RTU 485 Temperature And Humidity Sensor Driver",
	}
}

func (rtu485 *YK8RelayControllerDriver) Stop() error {
	rtu485.state = typex.DRIVER_STOP
	return nil
}
