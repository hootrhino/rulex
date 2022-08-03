package device

import (
	"context"
	"errors"
	"fmt"
	golog "log"
	"os"
	"sync"
	"time"

	modbus "github.com/wwhai/gomodbus"
	"github.com/i4de/rulex/common"
	"github.com/i4de/rulex/core"
	"github.com/i4de/rulex/driver"
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"

	"github.com/mitchellh/mapstructure"
)

//
// 这是个通用Modbus采集器, 主要用来在通用场景下采集数据，因此需要配合规则引擎来使用
//
// Modbus 采集到的数据如下, LUA 脚本可做解析, 示例脚本可参照 generic_modbus_parse.lua
// {
//     "d1":{
//         "tag":"d1",
//         "function":3,
//         "slaverId":1,
//         "address":0,
//         "quantity":2,
//         "value":"..."
//     },
//     "d2":{
//         "tag":"d2",
//         "function":3,
//         "slaverId":2,
//         "address":0,
//         "quantity":2,
//         "value":"..."
//     }
// }

type generic_modbus_device struct {
	typex.XStatus
	status     typex.DeviceState
	RuleEngine typex.RuleX
	driver     typex.XExternalDriver
	rtuHandler *modbus.RTUClientHandler
	tcpHandler *modbus.TCPClientHandler
	mainConfig common.ModBusConfig
	rtuConfig  common.RTUConfig
	tcpConfig  common.HostConfig
	locker     sync.Locker
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
	mdev.mainConfig = common.ModBusConfig{}
	mdev.tcpConfig = common.HostConfig{}
	mdev.rtuConfig = common.RTUConfig{}
	return mdev
}

//  初始化
func (mdev *generic_modbus_device) Init(devId string, configMap map[string]interface{}) error {
	mdev.PointId = devId
	if err := utils.BindSourceConfig(configMap, &mdev.mainConfig); err != nil {
		return err
	}
	// 检查Tag有没有重复
	tags := []string{}
	for _, register := range mdev.mainConfig.Registers {
		tags = append(tags, register.Tag)
	}
	if utils.IsListDuplicated(tags) {
		return errors.New("tag duplicated")
	}
	if !((mdev.mainConfig.Mode == "RTU") || (mdev.mainConfig.Mode == "TCP")) {
		return errors.New("unsupported mode, only can be one of 'TCP' or 'RTU'")
	}
	if mdev.mainConfig.Mode == "TCP" {
		if errs := mapstructure.Decode(mdev.mainConfig.Config, &mdev.tcpConfig); errs != nil {
			glogger.GLogger.Error(errs)
			return errs
		}
	}
	if mdev.mainConfig.Mode == "RTU" {
		if errs := mapstructure.Decode(mdev.mainConfig.Config, &mdev.rtuConfig); errs != nil {
			glogger.GLogger.Error(errs)
			return errs
		}
	}

	return nil
}

// 启动
func (mdev *generic_modbus_device) Start(cctx typex.CCTX) error {
	mdev.Ctx = cctx.Ctx
	mdev.CancelCTX = cctx.CancelCTX

	if mdev.mainConfig.Mode == "RTU" {
		mdev.rtuHandler = modbus.NewRTUClientHandler(mdev.rtuConfig.Uart)
		mdev.rtuHandler.BaudRate = mdev.rtuConfig.BaudRate
		mdev.rtuHandler.DataBits = mdev.rtuConfig.DataBits
		mdev.rtuHandler.Parity = mdev.rtuConfig.Parity
		mdev.rtuHandler.StopBits = mdev.rtuConfig.StopBits
		mdev.rtuHandler.Timeout = time.Duration(mdev.mainConfig.Frequency) * time.Second
		if core.GlobalConfig.AppDebugMode {
			mdev.rtuHandler.Logger = golog.New(os.Stdout, "485mdevSource: ", golog.LstdFlags)
		}

		if err := mdev.rtuHandler.Connect(); err != nil {
			return err
		}
		client := modbus.NewClient(mdev.rtuHandler)
		mdev.driver = driver.NewModBusRtuDriver(mdev.Details(),
			mdev.RuleEngine, mdev.mainConfig.Registers, mdev.rtuHandler, client)
	}
	if mdev.mainConfig.Mode == "TCP" {
		mdev.tcpHandler = modbus.NewTCPClientHandler(
			fmt.Sprintf("%s:%v", mdev.tcpConfig.Host, mdev.tcpConfig.Port),
		)
		if core.GlobalConfig.AppDebugMode {
			mdev.tcpHandler.Logger = golog.New(os.Stdout, "485mdevSource: ", golog.LstdFlags)
		}

		if err := mdev.tcpHandler.Connect(); err != nil {
			return err
		}
		client := modbus.NewClient(mdev.tcpHandler)
		mdev.driver = driver.NewModBusTCPDriver(mdev.Details(),
			mdev.RuleEngine, mdev.mainConfig.Registers, mdev.tcpHandler, client)
	}
	//---------------------------------------------------------------------------------
	// Start
	//---------------------------------------------------------------------------------
	mdev.status = typex.DEV_RUNNING

	go func(ctx context.Context, Driver typex.XExternalDriver) {
		ticker := time.NewTicker(time.Duration(5) * time.Second)
		defer ticker.Stop()
		buffer := make([]byte, common.T_64KB)
		for {
			<-ticker.C
			select {
			case <-ctx.Done():
				{
					mdev.status = typex.DEV_STOP
					return
				}
			default:
				{
				}
			}
			mdev.locker.Lock()
			n, err := Driver.Read(buffer)
			mdev.locker.Unlock()
			if err != nil {
				glogger.GLogger.Error(err)
			} else {
				mdev.RuleEngine.WorkDevice(mdev.Details(), string(buffer[:n]))
			}
		}

	}(mdev.Ctx, mdev.driver)
	return nil
}

// 从设备里面读数据出来
func (mdev *generic_modbus_device) OnRead(data []byte) (int, error) {

	n, err := mdev.driver.Read(data)
	if err != nil {
		glogger.GLogger.Error(err)
		mdev.status = typex.DEV_STOP
	}
	return n, err
}

// 把数据写入设备
func (mdev *generic_modbus_device) OnWrite(_ []byte) (int, error) {
	return 0, nil
}

// 设备当前状态
func (mdev *generic_modbus_device) Status() typex.DeviceState {
	return typex.DEV_RUNNING
}

// 停止设备
func (mdev *generic_modbus_device) Stop() {
	if mdev.tcpHandler != nil {
		mdev.tcpHandler.Close()
	}
	if mdev.rtuHandler != nil {
		mdev.rtuHandler.Close()
	}
	mdev.CancelCTX()
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
