// TC-S200 系列空气质量监测仪内置 PM2.5、TVOC、甲醛、CO2，温湿度等高精
// 度传感器套件，可通过吸顶式或壁挂安装，RS-485 接口通过 Modbus-RTU 协议进行
// 数据输出，通过网关组网，或配合联动模块可以用于新风联动控制。
// 该驱动是V0.2版本
package driver

import (
	"encoding/json"
	"rulex/typex"
	"rulex/utils"

	"github.com/goburrow/modbus"
)

type tss200_v_0_2_Driver struct {
	state      typex.DriverState
	client     modbus.Client
	device     *typex.Device
	RuleEngine typex.RuleX
}

func NewTSS200_v_0_2_Driver(d *typex.Device, e typex.RuleX,
	client modbus.Client) typex.XExternalDriver {
	return &tss200_v_0_2_Driver{
		state:      typex.DRIVER_STOP,
		device:     d,
		RuleEngine: e,
		client:     client,
	}
}
func (tss *tss200_v_0_2_Driver) Test() error {
	return nil
}

func (tss *tss200_v_0_2_Driver) Init(map[string]string) error {
	return nil
}

func (tss *tss200_v_0_2_Driver) Work() error {
	return nil
}

func (tss *tss200_v_0_2_Driver) State() typex.DriverState {
	return typex.DRIVER_RUNNING
}


type _sensor_data struct {
	TEMP float32 `json:"temp"` //系数: 0.01
	HUM  float32 `json:"hum"`  //系数: 0.01
	PM1  uint16  `json:"pm1"`
	PM25 uint16  `json:"pm25"`
	PM10 uint16  `json:"pm10"`
	CO2  uint16  `json:"co2"`
	TOVC float32 `json:"tovc"` //系数: 0.001
	CHOH float32 `json:"choh"` //系数: 0.001
}

func (tss *tss200_v_0_2_Driver) Read(data []byte) (int, error) {
	// 获取全部传感器数据：
	// |地址码|功能码|寄存器地址|寄存器长度|校验码|校验码
	// |XX    |03   |17       | 长度     |CRC  |  CRC
	// -----------------------------------------------
	// 01 03 00 11 00 08 14 09 09 45 1A F7 00 6F 00 89  00 89 FF FF FF FF 00 0B
	// TEMP  HUM     PM1   PM2.5  Pm10  CO2   TOVC  CHOH
	//
	result, err := tss.client.ReadHoldingRegisters(17, 8)
	if len(result) == 16 {

		sd := _sensor_data{
			TEMP: float32(utils.BToU16(result, 0, 2)) * 0.01,
			HUM:  float32(utils.BToU16(result, 2, 4)) * 0.01,
			PM1:  utils.BToU16(result, 4, 6),
			PM25: utils.BToU16(result, 6, 8),
			PM10: utils.BToU16(result, 8, 10),
			CO2:  utils.BToU16(result, 10, 12),
			TOVC: float32(utils.BToU16(result, 12, 14)) * 0.01,
			CHOH: float32(utils.BToU16(result, 14, 16)) * 0.01,
		}
		bytes, _ := json.Marshal(sd)
		copy(data, bytes)
	}

	return len(data), err
}

func (tss *tss200_v_0_2_Driver) Write(_ []byte) (int, error) {
	return 0, nil
}

//---------------------------------------------------
func (tss *tss200_v_0_2_Driver) DriverDetail() *typex.DriverDetail {
	return &typex.DriverDetail{
		Name:        "TC-S200",
		Type:        "UART",
		Description: "TC-S200 系列空气质量监测仪",
	}
}

func (tss *tss200_v_0_2_Driver) Stop() error {
	tss.state = typex.DRIVER_STOP
	return nil
}
