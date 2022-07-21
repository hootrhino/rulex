package device

import (
	"context"
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

type tss200_v_0_2_sensor struct {
	typex.XStatus
	status     typex.DeviceState
	RuleEngine typex.RuleX
	driver     typex.XExternalDriver
	slaverIds  []byte
	mainConfig common.ModBusConfig
	rtuConfig  common.RTUConfig
}

/*
*
* TSS200环境传感器
*
 */
func NewTS200Sensor(e typex.RuleX) typex.XDevice {
	tss := new(tss200_v_0_2_sensor)

	tss.RuleEngine = e
	return tss
}

//  初始化
func (tss *tss200_v_0_2_sensor) Init(devId string, configMap map[string]interface{}) error {
	tss.PointId = devId
	if err := utils.BindSourceConfig(configMap, &tss.mainConfig); err != nil {
		glogger.GLogger.Error(err)
		return err
	}

	if err := mapstructure.Decode(tss.mainConfig.Config, &tss.rtuConfig); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	return nil
}

// 启动
func (tss *tss200_v_0_2_sensor) Start(cctx typex.CCTX) error {
	tss.Ctx = cctx.Ctx
	tss.CancelCTX = cctx.CancelCTX

	// 串口配置固定写法
	handler := modbus.NewRTUClientHandler(tss.rtuConfig.Uart)
	handler.BaudRate = 4800
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.Timeout = time.Duration(5) * time.Second
	// handler.Logger = golog.New(os.Stdout, "485THerSource: ", glogger.GLogger.LstdFlags)
	if err := handler.Connect(); err != nil {
		return err
	}
	client := modbus.NewClient(handler)
	tss.driver = driver.NewTSS200_v_0_2_Driver(tss.Details(), tss.RuleEngine, client)
	//---------------------------------------------------------------------------------
	// Start
	//---------------------------------------------------------------------------------
	tss.status = typex.DEV_RUNNING
	lock := sync.Mutex{}
	for _, slaverId := range tss.slaverIds {
		go func(ctx context.Context, slaverId byte, rtuDriver typex.XExternalDriver, handler *modbus.RTUClientHandler) {
			ticker := time.NewTicker(time.Duration(5) * time.Second)
			defer ticker.Stop()
			buffer := make([]byte, 128) //128字节数据
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
				lock.Lock()
				handler.SlaveId = slaverId // 配置ID
				n, err := rtuDriver.Read(buffer)
				lock.Unlock()
				if err != nil {
					glogger.GLogger.Error(err)
				} else {
					tss.RuleEngine.WorkDevice(tss.RuleEngine.GetDevice(tss.PointId), string(buffer[:n]))
				}
			}

		}(tss.Ctx, slaverId, tss.driver, handler)
	}
	return nil
}

// 从设备里面读数据出来
func (tss *tss200_v_0_2_sensor) OnRead(data []byte) (int, error) {

	n, err := tss.driver.Read(data)
	if err != nil {
		glogger.GLogger.Error(err)
		tss.status = typex.DEV_STOP
	}
	return n, err
}

// 把数据写入设备
func (tss *tss200_v_0_2_sensor) OnWrite(b []byte) (int, error) {
	return tss.driver.Write(b)
}

// 设备当前状态
func (tss *tss200_v_0_2_sensor) Status() typex.DeviceState {
	return typex.DEV_RUNNING
}

// 停止设备
func (tss *tss200_v_0_2_sensor) Stop() {
	tss.CancelCTX()
}

// 设备属性，是一系列属性描述
func (tss *tss200_v_0_2_sensor) Property() []typex.DeviceProperty {
	return []typex.DeviceProperty{}
}

// 真实设备
func (tss *tss200_v_0_2_sensor) Details() *typex.Device {
	return tss.RuleEngine.GetDevice(tss.PointId)
}

// 状态
func (tss *tss200_v_0_2_sensor) SetState(status typex.DeviceState) {
	tss.status = status

}

// 驱动
func (tss *tss200_v_0_2_sensor) Driver() typex.XExternalDriver {
	return tss.driver
}
