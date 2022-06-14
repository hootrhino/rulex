package driver

//
// RS232/RS485 控制继电器模块，可利用电脑通过串口（没有串口的可利用 USB 转
// 串口）连接控制器进行对设备的控制，接口采用开关输出，有常开常闭点。
// 资料首页：http://www.yi-kun.com
//
import (
	"encoding/json"
	"errors"
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

func (yk8 *YK8RelayControllerDriver) Test() error {
	return nil
}

func (yk8 *YK8RelayControllerDriver) Init(map[string]string) error {
	return nil
}

func (yk8 *YK8RelayControllerDriver) Work() error {
	return nil
}

func (yk8 *YK8RelayControllerDriver) State() typex.DriverState {
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

/*
*
* 读数据
*
 */
func (yk8 *YK8RelayControllerDriver) Read(data []byte) (int, error) {
	results, err := yk8.client.ReadCoils(0x00, 1)
	if err != nil {
		return 0, err
	}
	if len(results) == 1 {
		yks := yk08sw{
			Sw1: byteToBool1(results[0], 0),
			Sw2: byteToBool1(results[0], 1),
			Sw3: byteToBool1(results[0], 2),
			Sw4: byteToBool1(results[0], 3),
			Sw5: byteToBool1(results[0], 4),
			Sw6: byteToBool1(results[0], 5),
			Sw7: byteToBool1(results[0], 6),
			Sw8: byteToBool1(results[0], 7),
		}
		bytes, _ := json.Marshal(yks)
		copy(data, bytes)
	}
	return len(data), err
}

//
// data = [1,1,1,1,1,1,1,1]
//
func (yk8 *YK8RelayControllerDriver) Write(data []byte) (int, error) {
	if len(data) != 8 {
		return 0, errors.New("操作继电器组最少8个布尔值")
	}
	for _, v := range data {
		if v == 0 {
			continue
		} else if v == 1 {
			continue
		} else {
			return 0, errors.New("必须是逻辑值")
		}
	}

	Sw1 := byteToBool2(data[0])
	Sw2 := byteToBool2(data[1])
	Sw3 := byteToBool2(data[2])
	Sw4 := byteToBool2(data[3])
	Sw5 := byteToBool2(data[4])
	Sw6 := byteToBool2(data[5])
	Sw7 := byteToBool2(data[6])
	Sw8 := byteToBool2(data[7])
	var value byte
	setABitOnByte(&value, 0, Sw1)
	setABitOnByte(&value, 1, Sw2)
	setABitOnByte(&value, 2, Sw3)
	setABitOnByte(&value, 3, Sw4)
	setABitOnByte(&value, 4, Sw5)
	setABitOnByte(&value, 5, Sw6)
	setABitOnByte(&value, 6, Sw7)
	setABitOnByte(&value, 7, Sw8)

	_, err := yk8.client.WriteMultipleCoils(0, 1, []byte{value})
	if err != nil {
		return 0, err
	}
	return 0, err
}

//---------------------------------------------------
func (yk8 *YK8RelayControllerDriver) DriverDetail() typex.DriverDetail {
	return typex.DriverDetail{
		Name:        "YK-08-RELAY CONTROLLER",
		Type:        "UART",
		Description: "一个支持RS232和485的国产8路继电器控制器",
	}
}

func (yk8 *YK8RelayControllerDriver) Stop() error {
	yk8.state = typex.DRIVER_STOP
	return nil
}

//--------------------------------------------------------------------------------------------------
// 内部函数
//--------------------------------------------------------------------------------------------------

/*
*
* 取某个字节上的位
*
 */
func getABitOnByte(b byte, position uint8) (v uint8) {
	mask := 0b00000001
	if position == 0 {
		return (b & byte(mask)) >> position
	} else {
		return (b & (1 << mask)) >> position
	}
}

/*
*
* 设置字节上的某个位
*
 */
func setABitOnByte(b *byte, position uint8, value bool) (byte, error) {
	if position > 7 {
		return 0, errors.New("下标必须是0-7, 高位在前, 低位在后")
	}
	if value {
		return *b & 0b1111_1111, nil
	} else {
		masks := []byte{
			0b11111110,
			0b11111101,
			0b11111011,
			0b11110111,
			0b11101111,
			0b11011111,
			0b10111111,
			0b01111111,
		}
		return *b & masks[position], nil
	}

}

/*
*
* 字节转逻辑
*
 */
func byteToBool1(data byte, index uint8) bool {
	return getABitOnByte(data, index) == 1
}
func byteToBool2(data byte) bool {
	return data == 1
}
