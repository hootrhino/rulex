// TC-S200 系列空气质量监测仪内置 PM2.5、TVOC、甲醛、CO2，温湿度等高精
// 度传感器套件，可通过吸顶式或壁挂安装，RS-485 接口通过 Modbus-RTU 协议进行
// 数据输出，通过网关组网，或配合联动模块可以用于新风联动控制。
// 该驱动是V0.2版本
package driver

import (
	"rulex/typex"

	"github.com/goburrow/modbus"
)

type ts200_v_0_2_Driver struct {
	state      typex.DriverState
	client     modbus.Client
	device     *typex.Device
	RuleEngine typex.RuleX
}

func NewTS200_v_0_2_Driver(d *typex.Device, e typex.RuleX,
	client modbus.Client) typex.XExternalDriver {
	return &ts200_v_0_2_Driver{
		state:      typex.RUNNING,
		device:     d,
		RuleEngine: e,
		client:     client,
	}
}
func (tss *ts200_v_0_2_Driver) Test() error {
	return nil
}

func (tss *ts200_v_0_2_Driver) Init() error {
	return nil
}

func (tss *ts200_v_0_2_Driver) Work() error {
	return nil
}

func (tss *ts200_v_0_2_Driver) State() typex.DriverState {
	return typex.RUNNING
}

func (tss *ts200_v_0_2_Driver) SetState(state typex.DriverState) {
	tss.state = state
}

func (tss *ts200_v_0_2_Driver) Read(data []byte) (int, error) {
	// 获取全部传感器数据：
	// |地址码|功能码|寄存器地址|寄存器长度|校验码|校验码
	// |XX    |03   |17       | 长度     |CRC  |  CRC
	results, err := tss.client.ReadHoldingRegisters(17, 6)
	copy(data, results)
	return len(results), err
}

func (tss *ts200_v_0_2_Driver) Write(_ []byte) (int, error) {
	return 0, nil
}

//---------------------------------------------------
func (tss *ts200_v_0_2_Driver) DriverDetail() *typex.DriverDetail {
	return &typex.DriverDetail{
		Name:        "TC-S200",
		Type:        "UART",
		Description: "TC-S200 系列空气质量监测仪",
	}
}

func (tss *ts200_v_0_2_Driver) Stop() error {
	tss.state = typex.STOP
	return nil
}
