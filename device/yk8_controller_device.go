package device

import (
	"context"
	"sync"
	"time"

	"github.com/i4de/rulex/driver"
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"

	"github.com/goburrow/modbus"
	"github.com/mitchellh/mapstructure"
)

type YK8Controller struct {
	typex.XStatus
	status     typex.DeviceState
	RuleEngine typex.RuleX
	driver     typex.XExternalDriver
	slaverIds  []byte
}

/*
*
* 8路继电器
*
 */
func NewYK8Controller(deviceId string, e typex.RuleX) typex.XDevice {
	yk8 := new(YK8Controller)
	yk8.PointId = deviceId
	yk8.RuleEngine = e
	return yk8
}

//  初始化
func (yk8 *YK8Controller) Init(devId string, config map[string]interface{}) error {

	return nil
}

// 启动
func (yk8 *YK8Controller) Start(cctx typex.CCTX) error {
	yk8.Ctx = cctx.Ctx
	yk8.CancelCTX = cctx.CancelCTX
	config := yk8.RuleEngine.GetDevice(yk8.PointId).Config
	var mainConfig modBusConfig
	if err := utils.BindSourceConfig(config, &mainConfig); err != nil {
		return err
	}
	var rtuConfig rtuConfig
	if errs := mapstructure.Decode(mainConfig.Config, &rtuConfig); errs != nil {
		glogger.GLogger.Error(errs)
		return errs
	}

	// 串口配置固定写法
	handler := modbus.NewRTUClientHandler(rtuConfig.Uart)
	handler.BaudRate = 9600
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.Timeout = time.Duration(*mainConfig.Timeout) * time.Second
	// handler.Logger = golog.New(os.Stdout, "485THerSource: ", glogger.GLogger.LstdFlags)
	if err := handler.Connect(); err != nil {
		return err
	}
	client := modbus.NewClient(handler)
	yk8.driver = driver.NewYK8RelayControllerDriver(yk8.Details(), yk8.RuleEngine, client)
	yk8.slaverIds = append(yk8.slaverIds, mainConfig.SlaverIds...)
	//---------------------------------------------------------------------------------
	// Start
	//---------------------------------------------------------------------------------
	yk8.status = typex.DEV_RUNNING
	lock := sync.Mutex{}
	for _, slaverId := range yk8.slaverIds {
		go func(ctx context.Context, slaverId byte, rtuDriver typex.XExternalDriver, handler *modbus.RTUClientHandler) {
			ticker := time.NewTicker(time.Duration(5) * time.Second)
			defer ticker.Stop()
			buffer := make([]byte, 128) //128字节数据
			for {
				<-ticker.C
				select {
				case <-ctx.Done():
					{
						yk8.status = typex.DEV_STOP
						return
					}
				default:
					{
					}
				}
				lock.Lock()
				handler.SlaveId = slaverId // 配置ID
				n, err := rtuDriver.Read(buffer)
				lock.Unlock()
				if err != nil {
					glogger.GLogger.Error(err)
				} else {
					td := yk8.RuleEngine.GetDevice(yk8.PointId)
					yk8.RuleEngine.WorkDevice(td, string(buffer[:n]))
				}
			}

		}(yk8.Ctx, slaverId, yk8.driver, handler)
	}
	return nil
}

// 从设备里面读数据出来
func (yk8 *YK8Controller) OnRead(data []byte) (int, error) {

	n, err := yk8.driver.Read(data)
	if err != nil {
		glogger.GLogger.Error(err)
		yk8.status = typex.DEV_STOP
	}
	return n, err
}

// 把数据写入设备
func (yk8 *YK8Controller) OnWrite(b []byte) (int, error) {
	return yk8.driver.Write(b)
}

// 设备当前状态
func (yk8 *YK8Controller) Status() typex.DeviceState {
	return typex.DEV_RUNNING
}

// 停止设备
func (yk8 *YK8Controller) Stop() {
	yk8.CancelCTX()
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
