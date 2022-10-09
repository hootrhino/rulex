package device

import (
	"context"
	"errors"
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
	"github.com/mitchellh/mapstructure"
	modbus "github.com/wwhai/gomodbus"
)

type YK8Controller struct {
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
* 8路继电器
*
 */
func NewYK8Controller(e typex.RuleX) typex.XDevice {
	yk8 := new(YK8Controller)
	yk8.locker = &sync.Mutex{}
	yk8.RuleEngine = e
	return yk8
}

//  初始化
func (yk8 *YK8Controller) Init(devId string, configMap map[string]interface{}) error {
	yk8.PointId = devId
	if err := utils.BindSourceConfig(configMap, &yk8.mainConfig); err != nil {
		glogger.GLogger.Error(err)
		return err
	}

	if errs := mapstructure.Decode(yk8.mainConfig.Config, &yk8.rtuConfig); errs != nil {
		glogger.GLogger.Error(errs)
		return errs
	}
	// 检查Tag有没有重复
	tags := []string{}
	for _, register := range yk8.mainConfig.Registers {
		tags = append(tags, register.Tag)
	}
	if utils.IsListDuplicated(tags) {
		return errors.New("tag duplicated")
	}
	return nil
}

// 启动
func (yk8 *YK8Controller) Start(cctx typex.CCTX) error {
	yk8.Ctx = cctx.Ctx
	yk8.CancelCTX = cctx.CancelCTX

	// 串口配置固定写法
	// 下面的参数是传感器固定写法
	yk8.rtuHandler = modbus.NewRTUClientHandler(yk8.rtuConfig.Uart)
	yk8.rtuHandler.BaudRate = yk8.rtuConfig.BaudRate
	yk8.rtuHandler.DataBits = yk8.rtuConfig.DataBits
	yk8.rtuHandler.Parity = yk8.rtuConfig.Parity
	yk8.rtuHandler.StopBits = yk8.rtuConfig.StopBits
	yk8.rtuHandler.Timeout = time.Duration(yk8.mainConfig.Timeout) * time.Second
	if core.GlobalConfig.AppDebugMode {
		yk8.rtuHandler.Logger = golog.New(os.Stdout, "YK8-DEVICE: ", golog.LstdFlags)
	}

	if err := yk8.rtuHandler.Connect(); err != nil {
		return err
	}
	//---------------------------------------------------------------------------------
	// Start
	//---------------------------------------------------------------------------------
	client := modbus.NewClient(yk8.rtuHandler)
	yk8.driver = driver.NewYK8RelayControllerDriver(yk8.Details(),
		yk8.RuleEngine, yk8.mainConfig.Registers, yk8.rtuHandler, client)
	yk8.status = typex.DEV_UP

	go func(ctx context.Context, Driver typex.XExternalDriver) {
		ticker := time.NewTicker(time.Duration(yk8.mainConfig.Frequency) * time.Second)
		buffer := make([]byte, common.T_64KB)
		yk8.driver.Read(buffer) //清理缓存
		for {
			<-ticker.C
			select {
			case <-ctx.Done():
				{
					yk8.status = typex.DEV_STOP
					ticker.Stop()
					return
				}
			default:
				{
				}
			}
			yk8.locker.Lock()
			n, err := Driver.Read(buffer)
			yk8.locker.Unlock()
			if err != nil {
				glogger.GLogger.Error(err)
			} else {
				yk8.RuleEngine.WorkDevice(yk8.Details(), string(buffer[:n]))
			}
		}

	}(yk8.Ctx, yk8.driver)
	return nil
}

// 从设备里面读数据出来
func (yk8 *YK8Controller) OnRead(data []byte) (int, error) {

	n, err := yk8.driver.Read(data)
	if err != nil {
		glogger.GLogger.Error(err)
		yk8.status = typex.DEV_DOWN
	}
	return n, err
}

// 把数据写入设备
func (yk8 *YK8Controller) OnWrite(b []byte) (int, error) {
	return yk8.driver.Write(b)
}

// 设备当前状态
func (yk8 *YK8Controller) Status() typex.DeviceState {
	return typex.DEV_UP
}

// 停止设备
func (yk8 *YK8Controller) Stop() {
	if yk8.rtuHandler != nil {
		yk8.rtuHandler.Close()
	}
	yk8.CancelCTX()
	yk8.status = typex.DEV_STOP
}

// 设备属性，是一系列属性描述
func (yk8 *YK8Controller) Property() []typex.DeviceProperty {
	return []typex.DeviceProperty{}
}

// 真实设备
func (yk8 *YK8Controller) Details() *typex.Device {
	return yk8.RuleEngine.GetDevice(yk8.PointId)
}

// 状态
func (yk8 *YK8Controller) SetState(status typex.DeviceState) {
	yk8.status = status

}

// 驱动
func (yk8 *YK8Controller) Driver() typex.XExternalDriver {
	return yk8.driver
}

func (yk8 *YK8Controller) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}
