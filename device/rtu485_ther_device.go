package device

import (
	"context"
	"errors"
	golog "log"
	"sync"
	"time"

	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/core"
	"github.com/hootrhino/rulex/driver"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
	"github.com/mitchellh/mapstructure"
	modbus "github.com/wwhai/gomodbus"
)

type rtu485_ther struct {
	typex.XStatus
	status     typex.DeviceState
	RuleEngine typex.RuleX
	driver     typex.XExternalDriver
	rtuHandler *modbus.RTUClientHandler
	mainConfig common.ModBusConfig
	rtuConfig  common.RTUConfig
	locker     sync.Locker
}

// Example: 0x02 0x92 0xFF 0x98
/*
*
* 温湿度传感器
*
 */
func NewRtu485Ther(e typex.RuleX) typex.XDevice {
	ther := new(rtu485_ther)
	ther.RuleEngine = e
	ther.locker = &sync.Mutex{}
	ther.mainConfig = common.ModBusConfig{}
	ther.rtuConfig = common.RTUConfig{}
	return ther
}

//  初始化
func (ther *rtu485_ther) Init(devId string, configMap map[string]interface{}) error {
	ther.PointId = devId
	if err := utils.BindSourceConfig(configMap, &ther.mainConfig); err != nil {
		return err
	}
	if errs := mapstructure.Decode(ther.mainConfig.Config, &ther.rtuConfig); errs != nil {
		glogger.GLogger.Error(errs)
		return errs
	}
	// 检查Tag有没有重复
	tags := []string{}
	for _, register := range ther.mainConfig.Registers {
		tags = append(tags, register.Tag)
	}

	if utils.IsListDuplicated(tags) {
		return errors.New("tag duplicated")
	}
	return nil
}

// 启动
func (ther *rtu485_ther) Start(cctx typex.CCTX) error {
	ther.Ctx = cctx.Ctx
	ther.CancelCTX = cctx.CancelCTX
	//
	// 串口配置固定写法
	ther.rtuHandler = modbus.NewRTUClientHandler(ther.rtuConfig.Uart)
	ther.rtuHandler.BaudRate = ther.rtuConfig.BaudRate
	ther.rtuHandler.DataBits = ther.rtuConfig.DataBits
	ther.rtuHandler.Parity = ther.rtuConfig.Parity
	ther.rtuHandler.StopBits = ther.rtuConfig.StopBits
	ther.rtuHandler.Timeout = time.Duration(ther.mainConfig.Timeout) * time.Second
	if core.GlobalConfig.AppDebugMode {
		ther.rtuHandler.Logger = golog.New(glogger.GLogger.Writer(), "485-TEMP-HUMI-DEVICE: ", golog.LstdFlags)
	}
	if err := ther.rtuHandler.Connect(); err != nil {
		return err
	}
	client := modbus.NewClient(ther.rtuHandler)
	ther.driver = driver.NewRtu485THerDriver(ther.Details(),
		ther.RuleEngine, ther.mainConfig.Registers, ther.rtuHandler, client)
	//---------------------------------------------------------------------------------
	// Start
	//---------------------------------------------------------------------------------
	if !ther.mainConfig.AutoRequest {
		ther.status = typex.DEV_UP
		return nil
	}
	go func(ctx context.Context, Driver typex.XExternalDriver) {
		ticker := time.NewTicker(time.Duration(ther.mainConfig.Frequency) * time.Millisecond)
		buffer := make([]byte, common.T_64KB)
		ther.driver.Read([]byte{}, buffer) //清理缓存
		for {
			<-ticker.C
			select {
			case <-ctx.Done():
				{
					ticker.Stop()
					if ther.rtuHandler != nil {
						ther.rtuHandler.Close()
					}
					return
				}
			default:
				{
				}
			}
			ther.locker.Lock()
			n, err := Driver.Read([]byte{}, buffer)
			ther.locker.Unlock()
			if err != nil {
				glogger.GLogger.Error(err)
			} else {
				ther.RuleEngine.WorkDevice(ther.Details(), string(buffer[:n]))
			}
		}

	}(ther.Ctx, ther.driver)
	ther.status = typex.DEV_UP
	return nil
}

// 从设备里面读数据出来
func (ther *rtu485_ther) OnRead(cmd []byte, data []byte) (int, error) {

	n, err := ther.driver.Read(cmd, data)
	if err != nil {
		glogger.GLogger.Error(err)
		ther.status = typex.DEV_DOWN
	}
	return n, err
}

// 把数据写入设备
func (ther *rtu485_ther) OnWrite(cmd []byte, _ []byte) (int, error) {
	return 0, nil
}

// 设备当前状态
func (ther *rtu485_ther) Status() typex.DeviceState {
	return ther.status
}

// 停止设备
func (ther *rtu485_ther) Stop() {
	ther.status = typex.DEV_DOWN
	ther.CancelCTX()

}

// 设备属性，是一系列属性描述
func (ther *rtu485_ther) Property() []typex.DeviceProperty {
	return []typex.DeviceProperty{}
}

// 真实设备
func (ther *rtu485_ther) Details() *typex.Device {
	return ther.RuleEngine.GetDevice(ther.PointId)
}

// 状态
func (ther *rtu485_ther) SetState(status typex.DeviceState) {
	ther.status = status

}

// 驱动
func (ther *rtu485_ther) Driver() typex.XExternalDriver {
	return ther.driver
}

func (ther *rtu485_ther) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}
func (ther *rtu485_ther) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	return []byte{}, nil
}
