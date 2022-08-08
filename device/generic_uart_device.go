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

type genericUartDevice struct {
	typex.XStatus
	status     typex.DeviceState
	RuleEngine typex.RuleX
	driver     typex.XExternalDriver
	mainConfig common.GenericUartConfig
	locker     sync.Locker
}

/*
*
* 通用串口透传
*
 */
func NewGenericUartDevice(e typex.RuleX) typex.XDevice {
	uart := new(genericUartDevice)
	uart.locker = &sync.Mutex{}
	uart.mainConfig = common.GenericUartConfig{}
	uart.RuleEngine = e
	return uart
}

//  初始化
func (uart *genericUartDevice) Init(devId string, configMap map[string]interface{}) error {
	uart.PointId = devId
	if err := utils.BindSourceConfig(configMap, &uart.mainConfig); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	if !contains([]string{"N", "E", "O"}, uart.mainConfig.Parity) {
		return errors.New("parity value only one of 'N','O','E'")
	}
	return nil
}

// 启动
func (uart *genericUartDevice) Start(cctx typex.CCTX) error {
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
		Timeout:  time.Duration(uart.mainConfig.Frequency) * time.Second,
	}
	serialPort, err := serial.Open(&config)
	if err != nil {
		glogger.GLogger.Error("rawUartDriver start failed:", err)
		return err
	}
	uart.driver = driver.NewRawUartDriver(uart.Ctx, uart.RuleEngine, uart.Details(), serialPort)
	go func(ctx context.Context) {
		ticker := time.NewTicker(time.Duration(uart.mainConfig.Frequency) * time.Second)
		buffer := make([]byte, common.T_64KB)
		uart.driver.Read(buffer) //清理缓存
		for {
			<-ticker.C
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			default:
				uart.locker.Lock()
				n, err := uart.driver.Read(buffer)
				uart.locker.Unlock()
				if err != nil {
					glogger.GLogger.Error(err)
				} else {
					mapV := map[string]interface{}{
						"tag":   uart.mainConfig.Tag,
						"value": string(buffer[:n]),
					}
					bytes, _ := json.Marshal(mapV)
					uart.RuleEngine.WorkDevice(uart.Details(), string(bytes))
				}
			}
		}

	}(uart.Ctx)
	return nil
}

// 从设备里面读数据出来
func (uart *genericUartDevice) OnRead(data []byte) (int, error) {
	return 0, nil
}

// 把数据写入设备
func (uart *genericUartDevice) OnWrite(b []byte) (int, error) {
	return uart.driver.Write(b)
}

// 设备当前状态
func (uart *genericUartDevice) Status() typex.DeviceState {
	return typex.DEV_UP
}

// 停止设备
func (uart *genericUartDevice) Stop() {
	if uart.driver != nil {
		uart.driver.Stop()
	}
	uart.CancelCTX()
	uart.status = typex.DEV_STOP
}

// 设备属性，是一系列属性描述
func (uart *genericUartDevice) Property() []typex.DeviceProperty {
	return []typex.DeviceProperty{}
}

// 真实设备
func (uart *genericUartDevice) Details() *typex.Device {
	return uart.RuleEngine.GetDevice(uart.PointId)
}

// 状态
func (uart *genericUartDevice) SetState(status typex.DeviceState) {
	uart.status = status

}

// 驱动
func (uart *genericUartDevice) Driver() typex.XExternalDriver {
	return uart.driver
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
