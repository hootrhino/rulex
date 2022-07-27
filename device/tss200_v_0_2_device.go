package device

import (
	"context"
	golog "log"
	"os"
	"sync"
	"time"

	"github.com/i4de/rulex/common"
	"github.com/i4de/rulex/core"
	"github.com/i4de/rulex/driver"
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"

	"github.com/goburrow/modbus"
	"github.com/mitchellh/mapstructure"
)

type tss200V2 struct {
	typex.XStatus
	status     typex.DeviceState
	RuleEngine typex.RuleX
	driver     typex.XExternalDriver
	rtuHandler *modbus.RTUClientHandler
	mainConfig common.ModBusConfig
	rtuConfig  common.RTUConfig
	locker     sync.Locker
}

/*
*
* TSS200环境传感器
*
 */
func NewTS200Sensor(e typex.RuleX) typex.XDevice {
	tss := new(tss200V2)
	tss.locker = &sync.Mutex{}
	tss.mainConfig = common.ModBusConfig{}
	tss.rtuConfig = common.RTUConfig{}
	tss.RuleEngine = e
	return tss
}

//  初始化
func (tss *tss200V2) Init(devId string, configMap map[string]interface{}) error {
	tss.PointId = devId
	if err := utils.BindSourceConfig(configMap, &tss.mainConfig); err != nil {
		glogger.GLogger.Error(err)
		return err
	}

	if errs := mapstructure.Decode(tss.mainConfig.Config, &tss.rtuConfig); errs != nil {
		glogger.GLogger.Error(errs)
		return errs
	}
	return nil
}

// 启动
func (tss *tss200V2) Start(cctx typex.CCTX) error {
	tss.Ctx = cctx.Ctx
	tss.CancelCTX = cctx.CancelCTX

	// 串口配置固定写法
	// 下面的参数是传感器固定写法
	tss.rtuHandler = modbus.NewRTUClientHandler(tss.rtuConfig.Uart)
	tss.rtuHandler.BaudRate = tss.rtuConfig.BaudRate
	tss.rtuHandler.DataBits = tss.rtuConfig.DataBits
	tss.rtuHandler.Parity = tss.rtuConfig.Parity
	tss.rtuHandler.StopBits = tss.rtuConfig.StopBits
	tss.rtuHandler.Timeout = time.Duration(tss.mainConfig.Frequency) * time.Second
	if core.GlobalConfig.AppDebugMode {
		tss.rtuHandler.Logger = golog.New(os.Stdout, "TSS200-DEVICE: ", golog.LstdFlags)
	}

	if err := tss.rtuHandler.Connect(); err != nil {
		return err
	}
	//---------------------------------------------------------------------------------
	// Start
	//---------------------------------------------------------------------------------
	client := modbus.NewClient(tss.rtuHandler)
	tss.driver = driver.NewTSS200Driver(tss.Details(),
		tss.RuleEngine, tss.mainConfig.Registers, tss.rtuHandler, client)
	tss.status = typex.DEV_RUNNING

	go func(ctx context.Context, Driver typex.XExternalDriver) {
		ticker := time.NewTicker(time.Duration(tss.mainConfig.Frequency) * time.Second)
		defer ticker.Stop()
		buffer := make([]byte, common.T_64KB)
		for {
			<-ticker.C
			select {
			case <-ctx.Done():
				{
					tss.status = typex.DEV_STOP
					return
				}
			default:
				{
				}
			}
			tss.locker.Lock()
			n, err := Driver.Read(buffer)
			tss.locker.Unlock()
			if err != nil {
				glogger.GLogger.Error(err)
			} else {
				tss.RuleEngine.WorkDevice(tss.Details(), string(buffer[:n]))
			}
		}

	}(tss.Ctx, tss.driver)
	return nil
}

// 从设备里面读数据出来
func (tss *tss200V2) OnRead(data []byte) (int, error) {

	n, err := tss.driver.Read(data)
	if err != nil {
		glogger.GLogger.Error(err)
		tss.status = typex.DEV_STOP
	}
	return n, err
}

// 把数据写入设备
func (tss *tss200V2) OnWrite(b []byte) (int, error) {
	return tss.driver.Write(b)
}

// 设备当前状态
func (tss *tss200V2) Status() typex.DeviceState {
	return typex.DEV_RUNNING
}

// 停止设备
func (tss *tss200V2) Stop() {
	if tss.rtuHandler != nil {
		tss.rtuHandler.Close()
	}
	tss.CancelCTX()
}

// 设备属性，是一系列属性描述
func (tss *tss200V2) Property() []typex.DeviceProperty {
	return []typex.DeviceProperty{}
}

// 真实设备
func (tss *tss200V2) Details() *typex.Device {
	return tss.RuleEngine.GetDevice(tss.PointId)
}

// 状态
func (tss *tss200V2) SetState(status typex.DeviceState) {
	tss.status = status

}

// 驱动
func (tss *tss200V2) Driver() typex.XExternalDriver {
	return tss.driver
}
