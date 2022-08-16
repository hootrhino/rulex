package device

import (
	"context"
	"sync"
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/i4de/rulex/common"
	"github.com/i4de/rulex/driver"
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"
)

type genericSnmpDevice struct {
	typex.XStatus
	status     typex.DeviceState
	RuleEngine typex.RuleX
	driver     typex.XExternalDriver
	locker     sync.Locker
	mainConfig common.GenericSnmpConfig
}

// Example: 0x02 0x92 0xFF 0x98
/*
*
* 温湿度传感器
*
 */
func NewGenericSnmpDevice(e typex.RuleX) typex.XDevice {
	sd := new(genericSnmpDevice)
	sd.RuleEngine = e
	sd.locker = &sync.Mutex{}
	sd.mainConfig = common.GenericSnmpConfig{}
	return sd
}

//  初始化
func (sd *genericSnmpDevice) Init(devId string, configMap map[string]interface{}) error {
	sd.PointId = devId
	if err := utils.BindSourceConfig(configMap, &sd.mainConfig); err != nil {
		return err
	}
	return nil
}

// 启动
func (sd *genericSnmpDevice) Start(cctx typex.CCTX) error {
	sd.Ctx = cctx.Ctx
	sd.CancelCTX = cctx.CancelCTX
	//
	// 串口配置固定写法
	gosnmp.Default.Target = sd.mainConfig.Target
	gosnmp.Default.Port = sd.mainConfig.Port
	gosnmp.Default.Community = sd.mainConfig.Community
	err := gosnmp.Default.Connect()
	if err != nil {
		glogger.GLogger.Fatalf("Connect() err: %v", err)
	}
	sd.driver = driver.NewSnmpDriver(sd.Details(), sd.RuleEngine, gosnmp.Default)
	//---------------------------------------------------------------------------------
	// Start
	//---------------------------------------------------------------------------------
	sd.status = typex.DEV_UP
	go func(ctx context.Context, Driver typex.XExternalDriver) {
		ticker := time.NewTicker(time.Duration(sd.mainConfig.Frequency) * time.Second)
		buffer := make([]byte, common.T_64KB)
		sd.driver.Read(buffer) //清理缓存
		for {
			<-ticker.C
			select {
			case <-ctx.Done():
				{
					sd.status = typex.DEV_STOP
					ticker.Stop()
					return
				}
			default:
				{
				}
			}
			n, err := Driver.Read(buffer)
			if err != nil {
				glogger.GLogger.Error(err)
			} else {
				sd.RuleEngine.WorkDevice(sd.Details(), string(buffer[:n]))
			}
		}

	}(sd.Ctx, sd.driver)
	return nil
}

// 从设备里面读数据出来
func (sd *genericSnmpDevice) OnRead(data []byte) (int, error) {

	n, err := sd.driver.Read(data)
	if err != nil {
		glogger.GLogger.Error(err)
		sd.status = typex.DEV_DOWN
	}
	return n, err
}

// 把数据写入设备
func (sd *genericSnmpDevice) OnWrite(_ []byte) (int, error) {
	return 0, nil
}

// 设备当前状态
func (sd *genericSnmpDevice) Status() typex.DeviceState {
	return sd.status
}

// 停止设备
func (sd *genericSnmpDevice) Stop() {
	if sd.driver != nil {
		sd.driver.Stop()
	}
	sd.CancelCTX()
}

// 设备属性，是一系列属性描述
func (sd *genericSnmpDevice) Property() []typex.DeviceProperty {
	return []typex.DeviceProperty{}
}

// 真实设备
func (sd *genericSnmpDevice) Details() *typex.Device {
	return sd.RuleEngine.GetDevice(sd.PointId)
}

// 状态
func (sd *genericSnmpDevice) SetState(status typex.DeviceState) {
	sd.status = status

}

// 驱动
func (sd *genericSnmpDevice) Driver() typex.XExternalDriver {
	return sd.driver
}
