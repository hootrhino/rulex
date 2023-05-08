package driver

//
// RS232/RS485 控制继电器模块，可利用电脑通过串口（没有串口的可利用 USB 转
// 串口）连接控制器进行对设备的控制，接口采用开关输出，有常开常闭点。
// 资料首页：http://www.yi-kun.com
//
import (
	"encoding/json"

	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/typex"
	modbus "github.com/wwhai/gomodbus"
)

type YK8RelayControllerDriver struct {
	state      typex.DriverState
	handler    *modbus.RTUClientHandler
	client     modbus.Client
	RuleEngine typex.RuleX
	Registers  []common.RegisterRW
	device     *typex.Device
}

func NewYK8RelayControllerDriver(d *typex.Device, e typex.RuleX,
	registers []common.RegisterRW,
	handler *modbus.RTUClientHandler,
	client modbus.Client) typex.XExternalDriver {
	return &YK8RelayControllerDriver{
		state:      typex.DRIVER_STOP,
		device:     d,
		RuleEngine: e,
		handler:    handler,
		Registers:  registers,
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
	return typex.DRIVER_UP
}

//

type yk08sw struct {
	Sw1 uint8 `json:"sw1"`
	Sw2 uint8 `json:"sw2"`
	Sw3 uint8 `json:"sw3"`
	Sw4 uint8 `json:"sw4"`
	Sw5 uint8 `json:"sw5"`
	Sw6 uint8 `json:"sw6"`
	Sw7 uint8 `json:"sw7"`
	Sw8 uint8 `json:"sw8"`
}

/*
*
* 读出来的是个JSON, 记录了8个开关的状态
*
 */
func (yk8 *YK8RelayControllerDriver) Read(cmd []byte, data []byte) (int, error) {
	dataMap := map[string]common.RegisterRW{}
	for _, r := range yk8.Registers {
		yk8.handler.SlaveId = r.SlaverId
		results, err := yk8.client.ReadCoils(0x00, 0x08)
		if err != nil {
			return 0, err
		}
		if len(results) == 1 {
			yks := yk08sw{
				Sw1: common.BitToUint8(results[0], 0),
				Sw2: common.BitToUint8(results[0], 1),
				Sw3: common.BitToUint8(results[0], 2),
				Sw4: common.BitToUint8(results[0], 3),
				Sw5: common.BitToUint8(results[0], 4),
				Sw6: common.BitToUint8(results[0], 5),
				Sw7: common.BitToUint8(results[0], 6),
				Sw8: common.BitToUint8(results[0], 7),
			}
			bytes, _ := json.Marshal(yks)
			value := common.RegisterRW{
				Tag:      r.Tag,
				Function: r.Function,
				SlaverId: r.SlaverId,
				Address:  r.Address,
				Quantity: r.Quantity,
				Value:    string(bytes),
			}
			dataMap[r.Tag] = value
		}
	}

	bytes, _ := json.Marshal(dataMap)
	copy(data, bytes)
	return len(bytes), nil
}

// 写入数据
func (yk8 *YK8RelayControllerDriver) Write(cmd []byte, data []byte) (int, error) {
	dataMap := []common.RegisterRW{}
	if err := json.Unmarshal(data, &dataMap); err != nil {
		return 0, err
	}
	for _, r := range dataMap {
		yk8.handler.SlaveId = r.SlaverId
		bytes, err0 := common.BitStringToBytes(string(r.Value))
		if err0 != nil {
			return 0, err0
		}
		_, err1 := yk8.client.WriteMultipleCoils(0, 1, bytes)
		if err1 != nil {
			return 0, err1
		}
	}
	return 0, nil
}

// ---------------------------------------------------
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
