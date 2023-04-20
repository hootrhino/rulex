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

type _GUDCommonConfig struct {
	Tag         string `json:"tag" validate:"required" title:"数据Tag" info:"给数据打标签"`
	Frequency   int64  `json:"frequency" validate:"required" title:"采集频率" info:""`
	AutoRequest bool   `json:"autoRequest" title:"启动轮询" info:""`
	Separator   string `json:"separator" title:"协议分隔符" info:""`
}

type _GUDConfig struct {
	commonConfig _GUDCommonConfig        `json:"commonConfig" validate:"required"`
	uartConfig   common.CommonUartConfig `json:"uartConfig" validate:"required"`
}

type genericUartDevice struct {
	typex.XStatus
	status     typex.DeviceState
	RuleEngine typex.RuleX
	driver     typex.XExternalDriver
	mainConfig _GUDConfig
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
	uart.mainConfig = _GUDConfig{}
	uart.RuleEngine = e
	return uart
}

//  初始化
func (uart *genericUartDevice) Init(devId string, configMap map[string]interface{}) error {
	uart.PointId = devId
	// 检查配置
	if uart.mainConfig.commonConfig.Separator == "" {
		uart.mainConfig.commonConfig.Separator = "\n"
	}
	if err := utils.BindSourceConfig(configMap, &uart.mainConfig); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	if !utils.SContains([]string{"N", "E", "O"}, uart.mainConfig.uartConfig.Parity) {
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
		Address:  uart.mainConfig.uartConfig.Uart,
		BaudRate: uart.mainConfig.uartConfig.BaudRate,
		DataBits: uart.mainConfig.uartConfig.DataBits,
		Parity:   uart.mainConfig.uartConfig.Parity,
		StopBits: uart.mainConfig.uartConfig.StopBits,
		Timeout:  time.Duration(uart.mainConfig.uartConfig.Timeout) * time.Second,
	}
	serialPort, err := serial.Open(&config)
	if err != nil {
		glogger.GLogger.Error("rawUartDriver start failed:", err)
		return err
	}
	uart.driver = driver.NewRawUartDriver(uart.Ctx, uart.RuleEngine, uart.Details(), serialPort)
	if !uart.mainConfig.commonConfig.AutoRequest {
		uart.status = typex.DEV_UP
		return nil
	}

	go func(ctx context.Context) {
		buffer := make([]byte, common.T_64KB) // 默认缓冲区64KB, 应该够了
		offset := 0
		for {
			select {
			case <-ctx.Done():
				return
			default:
				{
				}
			}
			uart.locker.Lock()
			len, err := uart.driver.Read([]byte{}, buffer[offset:])
			uart.locker.Unlock()
			if err != nil {
				glogger.GLogger.Error(err)
				if uart.status == typex.DEV_STOP {
					return
				} else {
					continue
				}
			}
			// 检查是否读到了协议结束符号, 只要发现结束符就提交, 移动指针
			for _, Byte := range buffer[offset : offset+len] {
				// 换行符 == 10
				Separator := uart.mainConfig.commonConfig.Separator[0]
				if Byte == Separator {
					mapV := map[string]string{
						"tag":   uart.mainConfig.commonConfig.Tag,
						"value": string(buffer[:offset]),
					}
					bytes, _ := json.Marshal(mapV)
					uart.RuleEngine.WorkDevice(uart.Details(), string(bytes))
					offset = 0
					break
				} else {
					offset += 1 // 一个一个移动
				}
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

// 全局缓冲
var _ReadBuffer []byte = make([]byte, common.T_64KB) // 默认缓冲区64KB, 应该够了
var _ReadBufferOffset int = 0

func (uart *genericUartDevice) OnRead(cmd []byte, data []byte) (int, error) {

	uart.driver.Read([]byte{}, _ReadBuffer[_ReadBufferOffset:]) //清理缓存
	uart.locker.Lock()
	n, err := uart.driver.Read([]byte{}, _ReadBuffer[_ReadBufferOffset:])
	uart.locker.Unlock()
	if err != nil {
		glogger.GLogger.Error(err)
		return 0, err
	}
	// 检查是否读到了协议结束符号, 只要发现结束符就提交, 移动指针
	for i := 0; i < n; i++ {
		if _ReadBuffer[i] == uart.mainConfig.commonConfig.Separator[0] {
			mapV := map[string]string{
				"tag":   uart.mainConfig.commonConfig.Tag,
				"value": string(_ReadBuffer[:n]),
			}
			bytes, _ := json.Marshal(mapV)
			uart.RuleEngine.WorkDevice(uart.Details(), string(bytes))
			copy(data, bytes)
			_ReadBufferOffset = 0
			continue
		} else {
			_ReadBufferOffset += n
			continue
		}
	}
	return 0, nil
}

// 把数据写入设备
func (uart *genericUartDevice) OnWrite(cmd []byte, b []byte) (int, error) {
	return uart.driver.Write(cmd, b)
}

// 设备当前状态
func (uart *genericUartDevice) Status() typex.DeviceState {
	return typex.DEV_UP
}

// 停止设备
func (uart *genericUartDevice) Stop() {
	uart.status = typex.DEV_STOP
	uart.CancelCTX()
	if uart.driver != nil {
		uart.driver.Stop()
		uart.driver = nil

	}
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

// --------------------------------------------------------------------------------------------------
//
// --------------------------------------------------------------------------------------------------

func (uart *genericUartDevice) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}
func (uart *genericUartDevice) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	return []byte{}, nil
}
