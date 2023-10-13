package device

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
	serial "github.com/wwhai/goserial"
)

/*
*
* 有人G776串口4G模块
*
 */
type _G776CommonConfig struct {
	Tag         string `json:"tag" validate:"required" title:"数据Tag" info:"给数据打标签"`
	Frequency   int64  `json:"frequency" validate:"required" title:"采集频率"`
	AutoRequest bool   `json:"autoRequest" title:"启动轮询"`
}

type _G776Config struct {
	CommonConfig _G776CommonConfig       `json:"commonConfig" validate:"required"`
	UartConfig   common.CommonUartConfig `json:"uartConfig" validate:"required"`
}

// 这是有人G776型号的4G DTU模块，主要用来TCP远程透传数据, 实际上就是个很普通的串口读写程序
// 详细文档: https://www.usr.cn/Download/806.html
type UsrG776DTU struct {
	typex.XStatus
	status     typex.DeviceState
	RuleEngine typex.RuleX
	mainConfig _G776Config
	locker     sync.Locker
	serialPort serial.Port
	errCount   int
}

/*
*
* 有人4G DTU
*
 */
func NewUsrG776DTU(e typex.RuleX) typex.XDevice {
	uart := new(UsrG776DTU)
	uart.locker = &sync.Mutex{}
	uart.mainConfig = _G776Config{}
	uart.RuleEngine = e
	uart.serialPort = nil
	return uart
}

//  初始化
func (uart *UsrG776DTU) Init(devId string, configMap map[string]interface{}) error {
	uart.PointId = devId
	if err := utils.BindSourceConfig(configMap, &uart.mainConfig); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	if !utils.SContains([]string{"N", "E", "O"}, uart.mainConfig.UartConfig.Parity) {
		return errors.New("parity value only one of 'N','O','E'")
	}
	return nil
}

// 启动
func (uart *UsrG776DTU) Start(cctx typex.CCTX) error {
	uart.Ctx = cctx.Ctx
	uart.CancelCTX = cctx.CancelCTX
	config := serial.Config{
		Address:  uart.mainConfig.UartConfig.Uart,
		BaudRate: uart.mainConfig.UartConfig.BaudRate,
		DataBits: uart.mainConfig.UartConfig.DataBits,
		Parity:   uart.mainConfig.UartConfig.Parity,
		StopBits: uart.mainConfig.UartConfig.StopBits,
		Timeout:  time.Duration(uart.mainConfig.UartConfig.Timeout) * time.Millisecond,
	}
	serialPort, err := serial.Open(&config)
	if err != nil {
		glogger.GLogger.Error("Serial.Open failed:", err)
		return err
	}
	uart.errCount = 0
	uart.serialPort = serialPort
	uart.status = typex.DEV_UP
	return nil
}

/*
*
* 不支持读, 仅仅是个数据透传DTU
*
 */
func (uart *UsrG776DTU) OnRead(cmd []byte, data []byte) (int, error) {
	return 0, fmt.Errorf("UsrG776DTU not support read data")
}

/*
*
* 有人G776-DTU写入串口的数据会被不加修改的透传到上层
* data:ToUsrG776DTU("uuid", "DATA", "data-....")
*
 */
func (uart *UsrG776DTU) OnWrite(cmd []byte, b []byte) (int, error) {
	if string(cmd) != ("DATA") {
		return 0, nil
	}
	n, err := uart.serialPort.Write(b)
	if err != nil {
		uart.errCount++
		glogger.GLogger.Error(err)
		if uart.errCount > 5 {
			return n, err
		}
	}
	return n, nil
}

// 设备当前状态
func (uart *UsrG776DTU) Status() typex.DeviceState {
	if uart.serialPort != nil {
		// https://www.usr.cn/Download/806.html
		//  发送： AT\r
		//  接收： \r\nOK\r\n\r\n
		_, err := uart.serialPort.Write([]byte("AT\r"))
		if err != nil {
			uart.errCount++
			glogger.GLogger.Error(err)
			if uart.errCount > 5 {
				return typex.DEV_DOWN
			}
		}
	}
	return typex.DEV_UP
}

// 停止设备
func (uart *UsrG776DTU) Stop() {
	uart.CancelCTX()
	uart.status = typex.DEV_DOWN
	if uart.serialPort != nil {
		uart.serialPort.Close()
	}
}

// 设备属性，是一系列属性描述
func (uart *UsrG776DTU) Property() []typex.DeviceProperty {
	return []typex.DeviceProperty{}
}

// 真实设备
func (uart *UsrG776DTU) Details() *typex.Device {
	return uart.RuleEngine.GetDevice(uart.PointId)
}

// 状态
func (uart *UsrG776DTU) SetState(status typex.DeviceState) {
	uart.status = status

}

// 驱动
func (uart *UsrG776DTU) Driver() typex.XExternalDriver {
	return nil
}

func (uart *UsrG776DTU) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}
func (uart *UsrG776DTU) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	return []byte{}, nil
}
