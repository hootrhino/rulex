package device

import (
	"context"
	"errors"
	"fmt"
	golog "log"
	"os"
	"sync"
	"time"

	"github.com/i4de/rulex/common"
	"github.com/i4de/rulex/driver"
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"

	"github.com/goburrow/modbus"
	"github.com/mitchellh/mapstructure"
)

var __debug2 bool = false

type generic_modbus_device struct {
	typex.XStatus
	status     typex.DeviceState
	RuleEngine typex.RuleX
	driver     typex.XExternalDriver
	rtuHandler *modbus.RTUClientHandler
	tcpHandler *modbus.TCPClientHandler
	mainConfig common.ModBusConfig
	rtuConfig  common.RTUConfig
	tcpConfig  common.TCPConfig
}

/*
*
* 温湿度传感器
*
 */
func NewGenericModbusDevice(e typex.RuleX) typex.XDevice {
	mdev := new(generic_modbus_device)
	mdev.RuleEngine = e
	return mdev
}

//  初始化
func (mdev *generic_modbus_device) Init(devId string, configMap map[string]interface{}) error {
	mdev.PointId = devId
	if err := utils.BindSourceConfig(configMap, &mdev.mainConfig); err != nil {
		return err
	}
	if (mdev.mainConfig.Mode != "RTU") || (mdev.mainConfig.Mode == "TCP") {
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
		mdev.rtuHandler.Timeout = time.Duration(5) * time.Second
		if __debug2 {
			mdev.rtuHandler.Logger = golog.New(os.Stdout, "485mdevSource: ", golog.LstdFlags)
		}

		if err := mdev.rtuHandler.Connect(); err != nil {
			return err
		}
	}
	if mdev.mainConfig.Mode == "TCP" {
		mdev.tcpHandler = modbus.NewTCPClientHandler(
			fmt.Sprintf("%s:%v", mdev.tcpConfig.Ip, mdev.tcpConfig.Port),
		)
		if __debug2 {
			mdev.rtuHandler.Logger = golog.New(os.Stdout, "485mdevSource: ", golog.LstdFlags)
		}

		if err := mdev.tcpHandler.Connect(); err != nil {
			return err
		}
	}
	if mdev.mainConfig.Mode == "TCP" {
		client := modbus.NewClient(mdev.tcpHandler)
		mdev.driver = driver.NewModBusTCPDriver(mdev.Details(),
			mdev.RuleEngine, mdev.mainConfig.Registers, mdev.tcpHandler, client)
	}
	if mdev.mainConfig.Mode == "RTU" {
		client := modbus.NewClient(mdev.rtuHandler)
		mdev.driver = driver.NewModBusRtuDriver(mdev.Details(),
			mdev.RuleEngine, mdev.mainConfig.Registers, mdev.rtuHandler, client)
	}
	//---------------------------------------------------------------------------------
	// Start
	//---------------------------------------------------------------------------------
	lock := sync.Mutex{}
	mdev.status = typex.DEV_RUNNING

	go func(ctx context.Context, Driver typex.XExternalDriver) {
		ticker := time.NewTicker(time.Duration(5) * time.Second)
		defer ticker.Stop()
		buffer := make([]byte, common.T_64KB) //32字节数据
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
			lock.Lock()
			n, err := Driver.Read(buffer)
			lock.Unlock()
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
