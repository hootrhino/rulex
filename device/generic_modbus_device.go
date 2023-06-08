package device

import (
	"context"
	"errors"
	"fmt"
	golog "log"
	"runtime"
	"sync"
	"time"

	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/core"
	"github.com/hootrhino/rulex/driver"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
	modbus "github.com/wwhai/gomodbus"
)

// 这是个通用Modbus采集器, 主要用来在通用场景下采集数据，因此需要配合规则引擎来使用
//
// Modbus 采集到的数据如下, LUA 脚本可做解析, 示例脚本可参照 generic_modbus_parse.lua
//
//	{
//	    "d1":{
//	        "tag":"d1",
//	        "function":3,
//	        "slaverId":1,
//	        "address":0,
//	        "quantity":2,
//	        "value":"..."
//	    },
//	    "d2":{
//	        "tag":"d2",
//	        "function":3,
//	        "slaverId":2,
//	        "address":0,
//	        "quantity":2,
//	        "value":"..."
//	    }
//	}
type _GMODCommonConfig struct {
	Mode        string `json:"mode" title:"工作模式" info:"RTU/TCP"`
	Timeout     int    `json:"timeout" validate:"required" title:"连接超时"`
	AutoRequest bool   `json:"autoRequest" title:"启动轮询"`
	Frequency   int64  `json:"frequency" validate:"required" title:"采集频率"`
}
type _GMODHostConfig struct {
	Host string `json:"host" title:"服务地址"`
	Port int    `json:"port" title:"服务端口"`
}

type _GMODConfig struct {
	CommonConfig _GMODCommonConfig       `json:"commonConfig" validate:"required"`
	RtuConfig    common.CommonUartConfig `json:"rtuConfig" validate:"required"`
	TcpConfig    _GMODHostConfig         `json:"tcpConfig" validate:"required"`
	Registers    []common.RegisterRW     `json:"registers" validate:"required" title:"寄存器配置"`
}
type generic_modbus_device struct {
	typex.XStatus ``
	status        typex.DeviceState
	RuleEngine    typex.RuleX
	driver        typex.XExternalDriver
	rtuHandler    *modbus.RTUClientHandler
	tcpHandler    *modbus.TCPClientHandler
	mainConfig    _GMODConfig
	locker        sync.Locker
	retryTimes    int
}

/*
*
* 温湿度传感器
*
 */
func NewGenericModbusDevice(e typex.RuleX) typex.XDevice {
	mdev := new(generic_modbus_device)
	mdev.RuleEngine = e
	mdev.locker = &sync.Mutex{}
	mdev.mainConfig = _GMODConfig{
		CommonConfig: _GMODCommonConfig{},
		TcpConfig:    _GMODHostConfig{Host: "127.0.0.1", Port: 502},
		RtuConfig:    common.CommonUartConfig{},
	}
	mdev.Busy = false
	mdev.status = typex.DEV_DOWN
	return mdev
}

//  初始化
func (mdev *generic_modbus_device) Init(devId string, configMap map[string]interface{}) error {
	mdev.PointId = devId
	if err := utils.BindSourceConfig(configMap, &mdev.mainConfig); err != nil {
		return err
	}
	// 超时大于30秒无意义
	if mdev.mainConfig.CommonConfig.Timeout > 30000 {
		return errors.New("'timeout' must less than 30 second")
	}
	// 频率不能太快
	if mdev.mainConfig.CommonConfig.Frequency < 50 {
		return errors.New("'frequency' must grate than 50 millisecond")

	}
	// 检查Tag有没有重复
	tags := []string{}
	for _, register := range mdev.mainConfig.Registers {
		tags = append(tags, register.Tag)
	}
	if utils.IsListDuplicated(tags) {
		return errors.New("tag duplicated")
	}
	if !utils.SContains([]string{"RTU", "TCP", "rtu", "tcp"}, mdev.mainConfig.CommonConfig.Mode) {
		return errors.New("unsupported mode, only can be one of 'TCP' or 'RTU'")
	}
	mdev.retryTimes = 0

	return nil
}

// 启动
func (mdev *generic_modbus_device) Start(cctx typex.CCTX) error {
	mdev.Ctx = cctx.Ctx
	mdev.CancelCTX = cctx.CancelCTX

	if mdev.mainConfig.CommonConfig.Mode == "RTU" {
		mdev.rtuHandler = modbus.NewRTUClientHandler(mdev.mainConfig.RtuConfig.Uart)
		mdev.rtuHandler.BaudRate = mdev.mainConfig.RtuConfig.BaudRate
		mdev.rtuHandler.DataBits = mdev.mainConfig.RtuConfig.DataBits
		mdev.rtuHandler.Parity = mdev.mainConfig.RtuConfig.Parity
		mdev.rtuHandler.StopBits = mdev.mainConfig.RtuConfig.StopBits
		// timeout 最大不能超过20, 不然无意义
		mdev.rtuHandler.Timeout = time.Duration(mdev.mainConfig.CommonConfig.Timeout) * time.Millisecond
		if core.GlobalConfig.AppDebugMode {
			mdev.rtuHandler.Logger = golog.New(glogger.GLogger.Writer(),
				"Modbus: ", golog.LstdFlags)
		}

		if err := mdev.rtuHandler.Connect(); err != nil {
			return err
		}
		client := modbus.NewClient(mdev.rtuHandler)
		mdev.driver = driver.NewModBusRtuDriver(mdev.Details(),
			mdev.RuleEngine, mdev.mainConfig.Registers, mdev.rtuHandler, client)
	}
	if mdev.mainConfig.CommonConfig.Mode == "TCP" {
		mdev.tcpHandler = modbus.NewTCPClientHandler(
			fmt.Sprintf("%s:%v", mdev.mainConfig.TcpConfig.Host, mdev.mainConfig.TcpConfig.Port),
		)
		if core.GlobalConfig.AppDebugMode {
			mdev.tcpHandler.Logger = golog.New(glogger.GLogger.Writer(), "Modbus: ", golog.LstdFlags)
		}

		if err := mdev.tcpHandler.Connect(); err != nil {
			return err
		}
		client := modbus.NewClient(mdev.tcpHandler)
		runtime.SetFinalizer(client, func(c modbus.Client) {
			println("runtime.SetFinalizer ===================>")

		})
		mdev.driver = driver.NewModBusTCPDriver(mdev.Details(),
			mdev.RuleEngine, mdev.mainConfig.Registers, mdev.tcpHandler, client)
	}
	//---------------------------------------------------------------------------------
	// Start
	//---------------------------------------------------------------------------------
	if !mdev.mainConfig.CommonConfig.AutoRequest {
		mdev.status = typex.DEV_UP
		return nil
	}
	go func(ctx context.Context, Driver typex.XExternalDriver) {

		mdev.status = typex.DEV_UP
		ticker := time.NewTicker(time.Duration(mdev.mainConfig.CommonConfig.Frequency) * time.Millisecond)
		buffer := make([]byte, common.T_64KB)
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			default:
				{
				}
			}
			if mdev.Busy {
				glogger.GLogger.Warn("Modbus device is busing now")
				continue
			}

			mdev.Busy = true
			n, err := Driver.Read([]byte{}, buffer)
			mdev.Busy = false
			if err != nil {
				glogger.GLogger.Error(err)
				mdev.retryTimes++
			} else {
				mdev.RuleEngine.WorkDevice(mdev.Details(), string(buffer[:n]))
			}
			<-ticker.C
		}

	}(mdev.Ctx, mdev.driver)
	return nil
}

// 从设备里面读数据出来
func (mdev *generic_modbus_device) OnRead(cmd []byte, data []byte) (int, error) {

	n, err := mdev.driver.Read(cmd, data)
	if err != nil {
		glogger.GLogger.Error(err)
		mdev.retryTimes++
	}
	return n, err
}

// 把数据写入设备
func (mdev *generic_modbus_device) OnWrite(cmd []byte, data []byte) (int, error) {
	if mdev.Busy {
		return 0, fmt.Errorf("device busing now")
	}
	return mdev.driver.Write(cmd, data)
}

// 设备当前状态
func (mdev *generic_modbus_device) Status() typex.DeviceState {
	if mdev.retryTimes > 3 {
		return typex.DEV_DOWN
	}
	return typex.DEV_UP
}

// 停止设备
func (mdev *generic_modbus_device) Stop() {
	if mdev.CancelCTX != nil {
		mdev.CancelCTX()
	}
	mdev.status = typex.DEV_DOWN
}

// 设备属性，是一系列属性描述
func (mdev *generic_modbus_device) Property() []typex.DeviceProperty {
	return []typex.DeviceProperty{}
}

// 真实设备
func (mdev *generic_modbus_device) Details() *typex.Device {
	return mdev.RuleEngine.GetDevice(mdev.PointId)
}

// 状态
func (mdev *generic_modbus_device) SetState(status typex.DeviceState) {
	mdev.status = status
}

// 驱动
func (mdev *generic_modbus_device) Driver() typex.XExternalDriver {
	return mdev.driver
}
func (mdev *generic_modbus_device) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}
func (mdev *generic_modbus_device) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	return []byte{}, nil
}
