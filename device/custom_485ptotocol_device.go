package device

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"
	serial "github.com/wwhai/goserial"
	"net"
	"sync"
	"time"
)

// 传输形式：
// `rawtcp`, `rawudp`, `rs485rawserial`, `rs485rawtcp`
const rawtcp string = "rawtcp"
const rawudp string = "rawudp"
const rs485rawserial string = "rs485rawserial"
const rs485rawtcp string = "rs485rawtcp"

type _CommonConfig struct {
	Frequency   int    `json:"frequency" validate:"required"`
	AutoRequest bool   `json:"autoRequest" validate:"required"`
	Transport   string `json:"transport" validate:"required"`
	WaitTime    int    `json:"waitTime" validate:"required"`
	Timeout     int    `json:"timeout" validate:"required"`
}
type _UartConfig struct {
	Uart     string `json:"uart" validate:"required"`
	BaudRate int    `json:"baudRate" validate:"required"`
	DataBits int    `json:"dataBits" validate:"required"`
	Parity   string `json:"parity" validate:"required"`
	StopBits int    `json:"stopBits" validate:"required"`
}
type _ProtocolArg struct {
	In  string `json:"in" validate:"required"`
	Out string `json:"out" validate:"required"`
}
type _Protocol struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	ProtocolArg _ProtocolArg
}

/*
*
* 自定义协议
*
 */
type _CustomProtocolConfig struct {
	CommonConfig _CommonConfig        `json:"commonConfig" validate:"required"`
	UartConfig   _UartConfig          `json:"uartConfig" validate:"required"`
	DeviceConfig map[string]_Protocol `json:"deviceConfig" validate:"required"`
}
type CustomProtocolDevice struct {
	typex.XStatus
	status     typex.DeviceState
	RuleEngine typex.RuleX
	serialPort serial.Port  // 现阶段暂时支持串口
	tcpConn    *net.TCPConn // rawtcp 以后支持
	udpConn    *net.UDPConn // rawudp 以后支持
	mainConfig _CustomProtocolConfig
	locker     sync.Locker
	errorCount int // 记录最大容错数，默认5次，出错超过5此就重启
}

func NewCustomProtocolDevice(e typex.RuleX) typex.XDevice {
	mdev := new(CustomProtocolDevice)
	mdev.RuleEngine = e
	mdev.locker = &sync.Mutex{}
	mdev.mainConfig = _CustomProtocolConfig{}
	mdev.status = typex.DEV_DOWN
	mdev.errorCount = 0
	return mdev

}

//  初始化
func (mdev *CustomProtocolDevice) Init(devId string, configMap map[string]interface{}) error {
	mdev.PointId = devId
	if err := utils.BindSourceConfig(configMap, &mdev.mainConfig); err != nil {
		return err
	}
	if !contains([]string{"N", "E", "O"}, mdev.mainConfig.UartConfig.Parity) {
		return errors.New("parity value only one of 'N','O','E'")
	}
	if !contains([]string{`rawtcp`, `rawudp`, `rs485rawserial`, `rs485rawtcp`},
		mdev.mainConfig.CommonConfig.Transport) {
		return errors.New("parity value only one of 'rawtcp','rawudp','rs485rawserial','rs485rawserial'")
	}
	// parse hex format
	for _, v := range mdev.mainConfig.DeviceConfig {
		if _, err := hex.DecodeString(v.ProtocolArg.In); err != nil {
			errMsg := fmt.Sprintf("hex.DecodeString(ProtocolArg.In) failed:%s", v.ProtocolArg.In)
			glogger.GLogger.Error(errMsg)
			return fmt.Errorf(errMsg)
		}
		if _, err := hex.DecodeString(v.ProtocolArg.Out); err != nil {
			errMsg := fmt.Sprintf("hex.DecodeString(ProtocolArg.Out) failed:%s", v.ProtocolArg.Out)
			glogger.GLogger.Error(errMsg)
			return fmt.Errorf(errMsg)
		}

	}
	return nil
}

// 启动
func (mdev *CustomProtocolDevice) Start(cctx typex.CCTX) error {
	mdev.Ctx = cctx.Ctx
	mdev.CancelCTX = cctx.CancelCTX
	// 现阶段暂时只支持RS485串口, 以后有需求再支持TCP、UDP
	if mdev.mainConfig.CommonConfig.Transport == "rs485rawserial" {
		config := serial.Config{
			Address:  mdev.mainConfig.UartConfig.Uart,
			BaudRate: mdev.mainConfig.UartConfig.BaudRate,
			DataBits: mdev.mainConfig.UartConfig.DataBits,
			Parity:   mdev.mainConfig.UartConfig.Parity,
			StopBits: mdev.mainConfig.UartConfig.StopBits,
			Timeout:  time.Duration(mdev.mainConfig.CommonConfig.Timeout) * time.Second,
		}
		serialPort, err := serial.Open(&config)
		if err != nil {
			glogger.GLogger.Error("serialPort start failed:", err)
			return err
		}
		mdev.serialPort = serialPort
		mdev.status = typex.DEV_UP
	}

	return fmt.Errorf("unsupported transport:%s", mdev.mainConfig.CommonConfig.Transport)
}

// 从设备里面读数据出来
func (mdev *CustomProtocolDevice) OnRead(cmd int, data []byte) (int, error) {

	return 0, nil
}

// 把数据写入设备
// 根据第二个参数来找配置进去的自定义协议, 必须进来一个可识别的指令
// 其中cmd常为0,为无意义参数
func (mdev *CustomProtocolDevice) OnWrite(_ int, data []byte) (int, error) {
	pp, ok := mdev.mainConfig.DeviceConfig[string(data)]
	if ok {
		hexs, err0 := hex.DecodeString(pp.ProtocolArg.In)
		if err0 != nil {
			glogger.GLogger.Error(err0)
			mdev.errorCount++
			return 0, err0
		}
		mdev.locker.Lock()
		// Send
		if _, err1 := mdev.serialPort.Write(hexs); err1 != nil {
			glogger.GLogger.Error(err1)
			mdev.errorCount++
			return 0, err1
		}
		// 同步等待60毫秒
		time.Sleep(60 * time.Microsecond)
		result := [50]byte{}
		n, err2 := mdev.serialPort.Read(result[:])
		if err2 != nil {
			glogger.GLogger.Error(err2)
			mdev.errorCount++
			return 0, err2
		}
		mdev.locker.Unlock()
		// 判断返回值, 把返回值给加工成大写Hex格式
		copy(data, []byte(fmt.Sprintf("%X", result[:n]))) // 把结果返回
		return n, nil
	}
	return 0, errors.New("unknown command:" + string(data))
}

// 设备当前状态
func (mdev *CustomProtocolDevice) Status() typex.DeviceState {
	if mdev.errorCount >= 5 {
		mdev.status = typex.DEV_DOWN
	}
	return mdev.status
}

// 停止设备
func (mdev *CustomProtocolDevice) Stop() {
	mdev.status = typex.DEV_STOP
	mdev.CancelCTX()

}

// 设备属性，是一系列属性描述
func (mdev *CustomProtocolDevice) Property() []typex.DeviceProperty {
	return []typex.DeviceProperty{}
}

// 真实设备
func (mdev *CustomProtocolDevice) Details() *typex.Device {
	return mdev.RuleEngine.GetDevice(mdev.PointId)
}

// 状态
func (mdev *CustomProtocolDevice) SetState(status typex.DeviceState) {
	mdev.status = status
}

// 驱动
func (mdev *CustomProtocolDevice) Driver() typex.XExternalDriver {
	return nil
}
func (mdev *CustomProtocolDevice) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}
