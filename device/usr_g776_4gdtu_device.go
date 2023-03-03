package device

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/i4de/rulex/common"
	"github.com/i4de/rulex/driver"
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"
	serial "github.com/wwhai/goserial"
)

// 这是有人G776型号的4G DTU模块，主要用来TCP远程透传数据, 实际上就是个很普通的串口读写程序
// 详细文档: https://www.usr.cn/Download/806.html
type UsrG776DTU struct {
	typex.XStatus
	status     typex.DeviceState
	RuleEngine typex.RuleX
	driver     typex.XExternalDriver
	mainConfig common.GenericUartConfig
	locker     sync.Locker
}

/*
*
* 有人4G DTU
*
 */
func NewUsrG776DTU(e typex.RuleX) typex.XDevice {
	uart := new(UsrG776DTU)
	uart.locker = &sync.Mutex{}
	uart.mainConfig = common.GenericUartConfig{}
	uart.RuleEngine = e
	return uart
}

//  初始化
func (uart *UsrG776DTU) Init(devId string, configMap map[string]interface{}) error {
	uart.PointId = devId
	if err := utils.BindSourceConfig(configMap, &uart.mainConfig); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	if !utils.SContains([]string{"N", "E", "O"}, uart.mainConfig.Parity) {
		return errors.New("parity value only one of 'N','O','E'")
	}
	return nil
}

// 启动
func (uart *UsrG776DTU) Start(cctx typex.CCTX) error {
	uart.Ctx = cctx.Ctx
	uart.CancelCTX = cctx.CancelCTX

	// 串口配置固定写法
	// 下面的参数是传感器固定写法
	config := serial.Config{
		Address:  uart.mainConfig.Uart,
		BaudRate: uart.mainConfig.BaudRate,
		DataBits: uart.mainConfig.DataBits,
		Parity:   uart.mainConfig.Parity,
		StopBits: uart.mainConfig.StopBits,
		Timeout:  time.Duration(uart.mainConfig.Timeout) * time.Second,
	}
	serialPort, err := serial.Open(&config)
	if err != nil {
		glogger.GLogger.Error("G776Driver start failed:", err)
		return err
	}
	uart.driver = driver.NewRawUartDriver(uart.Ctx, uart.RuleEngine, uart.Details(), serialPort)
	if !uart.mainConfig.AutoRequest {
		uart.status = typex.DEV_UP
		return nil
	}
	go func(ctx context.Context) {
		ticker := time.NewTicker(time.Duration(uart.mainConfig.Frequency) * time.Second)
		buffer := make([]byte, common.T_64KB)
		uart.driver.Read(0, buffer) //清理缓存
		for {
			<-ticker.C
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			default:
				uart.locker.Lock()
				n, err := uart.driver.Read(0, buffer)
				uart.locker.Unlock()
				if err != nil {
					glogger.GLogger.Error(err)
					continue
				}
				mapV := map[string]interface{}{
					"tag":   uart.mainConfig.Tag,
					"value": string(buffer[:n]),
				}
				bytes, _ := json.Marshal(mapV)
				uart.RuleEngine.WorkDevice(uart.Details(), string(bytes))
			}
		}

	}(uart.Ctx)
	uart.status = typex.DEV_UP
	return nil
}

// 从设备里面读数据出来:
//
//	{
//	    "tag":"data tag",
//	    "value":"value s"
//	}
func (uart *UsrG776DTU) OnRead(cmd int, data []byte) (int, error) {
	uart.locker.Lock()
	n, err := uart.driver.Read(cmd, data)
	uart.locker.Unlock()
	buffer := make([]byte, n)
	mapV := map[string]interface{}{
		"tag":   uart.mainConfig.Tag,
		"value": string(buffer[:n]),
	}
	bytes, _ := json.Marshal(mapV)
	copy(data, bytes)
	return n, err
}

// 把数据写入设备
func (uart *UsrG776DTU) OnWrite(cmd int, b []byte) (int, error) {
	return uart.driver.Write(cmd, b)
}

// 设备当前状态
func (uart *UsrG776DTU) Status() typex.DeviceState {
	return typex.DEV_UP
}

// 停止设备
func (uart *UsrG776DTU) Stop() {
	uart.status = typex.DEV_STOP
	uart.CancelCTX()
	if uart.driver != nil {
		uart.driver.Stop()
		uart.driver = nil
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
	return uart.driver
}

func (uart *UsrG776DTU) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}
