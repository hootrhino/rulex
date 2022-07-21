package device

import (
	"context"
	"encoding/json"
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

type rtu485_ther struct {
	typex.XStatus
	status     typex.DeviceState
	RuleEngine typex.RuleX
	driver     typex.XExternalDriver
	handler    *modbus.RTUClientHandler
	slaverIds  []byte
	mainConfig common.ModBusConfig
	rtuConfig  common.RTUConfig
}

var __debug bool = false

// Example: 0x02 0x92 0xFF 0x98
type __sensor_data struct {
	TEMP float32 `json:"temp"` //系数: 0.1
	HUM  float32 `json:"hum"`  //系数: 0.1
}

/*
*
* 温湿度传感器
*
 */
func NewRtu485Ther(deviceId string, e typex.RuleX) typex.XDevice {
	ther := new(rtu485_ther)
	ther.PointId = deviceId
	ther.RuleEngine = e
	return ther
}

//  初始化
func (ther *rtu485_ther) Init(devId string, configMap map[string]interface{}) error {
	if err := utils.BindSourceConfig(configMap, &ther.mainConfig); err != nil {
		return err
	}
	if errs := mapstructure.Decode(ther.mainConfig.Config, &ther.rtuConfig); errs != nil {
		glogger.GLogger.Error(errs)
		return errs
	}
	return nil
}

// 启动
func (ther *rtu485_ther) Start(cctx typex.CCTX) error {
	ther.Ctx = cctx.Ctx
	ther.CancelCTX = cctx.CancelCTX
	config := ther.RuleEngine.GetDevice(ther.PointId).Config
	var mainConfig common.ModBusConfig
	if err := utils.BindSourceConfig(config, &mainConfig); err != nil {
		return err
	}
	var rtuConfig common.RTUConfig
	if errs := mapstructure.Decode(mainConfig.Config, &rtuConfig); errs != nil {
		glogger.GLogger.Error(errs)
		return errs
	}

	// 串口配置固定写法
	ther.handler = modbus.NewRTUClientHandler(rtuConfig.Uart)
	ther.handler.BaudRate = 4800
	ther.handler.DataBits = 8
	ther.handler.Parity = "N"
	ther.handler.StopBits = 1
	ther.handler.Timeout = time.Duration(5) * time.Second
	if __debug {
		ther.handler.Logger = golog.New(os.Stdout, "485THerSource: ", golog.LstdFlags)
	}
	if err := ther.handler.Connect(); err != nil {
		return err
	}
	client := modbus.NewClient(ther.handler)
	ther.driver = driver.NewRtu485THerDriver(ther.Details(), ther.RuleEngine, client)
	ther.slaverIds = append(ther.slaverIds, ther.mainConfig.SlaverIds...)
	//---------------------------------------------------------------------------------
	// Start
	//---------------------------------------------------------------------------------
	lock := sync.Mutex{}
	ther.status = typex.DEV_RUNNING
	for _, slaverId := range ther.slaverIds {
		go func(ctx context.Context, slaverId byte,
			rtuDriver typex.XExternalDriver,
			handler *modbus.RTUClientHandler) {
			ticker := time.NewTicker(time.Duration(5) * time.Second)
			defer ticker.Stop()
			// {"SlaveId":1,"Data":"{\"temp\":28.7,\"hum\":66.1}}
			buffer := make([]byte, 64) //32字节数据
			for {
				<-ticker.C
				select {
				case <-ctx.Done():
					{
						ther.status = typex.DEV_STOP
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
					Device := ther.RuleEngine.GetDevice(ther.PointId)
					sdata := __sensor_data{}
					json.Unmarshal(buffer[:n], &sdata)
					bytes, _ := json.Marshal(map[string]interface{}{
						"slaveId": handler.SlaveId,
						"data":    sdata,
					})
					ther.RuleEngine.WorkDevice(Device, string(bytes))
				}
			}

		}(ther.Ctx, slaverId, ther.driver, ther.handler)
	}
	return nil
}

// 从设备里面读数据出来
func (ther *rtu485_ther) OnRead(data []byte) (int, error) {

	n, err := ther.driver.Read(data)
	if err != nil {
		glogger.GLogger.Error(err)
		ther.status = typex.DEV_STOP
	}
	return n, err
}

// 把数据写入设备
func (ther *rtu485_ther) OnWrite(_ []byte) (int, error) {
	return 0, nil
}

// 设备当前状态
func (ther *rtu485_ther) Status() typex.DeviceState {
	return typex.DEV_RUNNING
}

// 停止设备
func (ther *rtu485_ther) Stop() {
	if ther.handler != nil {
		ther.handler.Close()
	}
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
