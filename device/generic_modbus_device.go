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

type generic_modbus_device struct {
	typex.XStatus
	status     typex.DeviceState
	RuleEngine typex.RuleX
	driver     typex.XExternalDriver
	handler    *modbus.RTUClientHandler
	slaverIds  []byte
	mainConfig common.ModBusConfig
	rtuConfig  common.RTUConfig
}

/*
*
* 温湿度传感器
*
 */
func NewGenericModbusDevice(deviceId string, e typex.RuleX) typex.XDevice {
	mdev := new(generic_modbus_device)
	mdev.PointId = deviceId
	mdev.RuleEngine = e
	return mdev
}

//  初始化
func (mdev *generic_modbus_device) Init(devId string, configMap map[string]interface{}) error {
	if err := utils.BindSourceConfig(configMap, &mdev.mainConfig); err != nil {
		return err
	}
	if errs := mapstructure.Decode(mdev.mainConfig.Config, &mdev.rtuConfig); errs != nil {
		glogger.GLogger.Error(errs)
		return errs
	}
	return nil
}

// 启动
func (mdev *generic_modbus_device) Start(cctx typex.CCTX) error {
	mdev.Ctx = cctx.Ctx
	mdev.CancelCTX = cctx.CancelCTX

	// 串口配置固定写法
	mdev.handler = modbus.NewRTUClientHandler(mdev.rtuConfig.Uart)
	mdev.handler.BaudRate = mdev.rtuConfig.BaudRate
	mdev.handler.DataBits = 8
	mdev.handler.Parity = "N"
	mdev.handler.StopBits = 1
	mdev.handler.Timeout = time.Duration(5) * time.Second
	if __debug {
		mdev.handler.Logger = golog.New(os.Stdout, "485mdevSource: ", golog.LstdFlags)
	}
	if err := mdev.handler.Connect(); err != nil {
		return err
	}
	client := modbus.NewClient(mdev.handler)
	mdev.driver = driver.NewModBusRtuDriver(mdev.Details(), mdev.RuleEngine, nil, client)
	mdev.slaverIds = append(mdev.slaverIds, mdev.mainConfig.SlaverIds...)
	//---------------------------------------------------------------------------------
	// Start
	//---------------------------------------------------------------------------------
	lock := sync.Mutex{}
	mdev.status = typex.DEV_RUNNING
	for _, slaverId := range mdev.slaverIds {
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
						mdev.status = typex.DEV_STOP
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
					Device := mdev.RuleEngine.GetDevice(mdev.PointId)
					sdata := __sensor_data{}
					json.Unmarshal(buffer[:n], &sdata)
					bytes, _ := json.Marshal(map[string]interface{}{
						"slaveId": handler.SlaveId,
						"data":    sdata,
					})
					mdev.RuleEngine.WorkDevice(Device, string(bytes))
				}
			}

		}(mdev.Ctx, slaverId, mdev.driver, mdev.handler)
	}
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
	if mdev.handler != nil {
		mdev.handler.Close()
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
