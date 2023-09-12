package device

import (
	"context"
	"sync"
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/driver"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
)

type _SNMPCommonConfig struct {
	AutoRequest bool  `json:"autoRequest" title:"启动轮询"`
	Frequency   int64 `json:"frequency" validate:"required" title:"采集频率"`
}

type _GSNMPConfig struct {
	CommonConfig _SNMPCommonConfig        `json:"commonConfig" validate:"required"`
	SNMPConfig   common.GenericSnmpConfig `json:"snmpConfig" validate:"required"`
}

type genericSnmpDevice struct {
	typex.XStatus
	status     typex.DeviceState
	RuleEngine typex.RuleX
	driver     typex.XExternalDriver
	locker     sync.Locker
	mainConfig _GSNMPConfig
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
	sd.mainConfig = _GSNMPConfig{}
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
	client := &gosnmp.GoSNMP{
		Target:             sd.mainConfig.SNMPConfig.Target,
		Port:               sd.mainConfig.SNMPConfig.Port,
		Community:          sd.mainConfig.SNMPConfig.Community,
		Transport:          "udp",
		Version:            1,
		Timeout:            time.Duration(5) * time.Second,
		Retries:            3,
		ExponentialTimeout: true,
		MaxOids:            60,
	}
	err := client.Connect()
	if err != nil {
		glogger.GLogger.Errorf("Connect err: %v", err)
		return err
	}

	sd.driver = driver.NewSnmpDriver(sd.Details(), sd.RuleEngine, client)
	//---------------------------------------------------------------------------------
	// Start
	//---------------------------------------------------------------------------------
	if !sd.mainConfig.CommonConfig.AutoRequest {
		sd.status = typex.DEV_UP
		return nil
	}
	go func(ctx context.Context, Driver typex.XExternalDriver) {
		ticker := time.NewTicker(time.Duration(sd.mainConfig.CommonConfig.Frequency) * time.Millisecond)
		buffer := make([]byte, common.T_64KB)
		sd.driver.Read([]byte{}, buffer) //清理缓存
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
			n, err := Driver.Read([]byte{}, buffer)
			if err != nil {
				glogger.GLogger.Error(err)
			} else {
				s := string(buffer[:n])
				sd.RuleEngine.WorkDevice(sd.Details(), s)
			}
		}

	}(sd.Ctx, sd.driver)
	sd.status = typex.DEV_UP
	return nil
}

// 从设备里面读数据出来
func (sd *genericSnmpDevice) OnRead(cmd []byte, data []byte) (int, error) {

	n, err := sd.driver.Read(cmd, data)
	if err != nil {
		glogger.GLogger.Error(err)
		return n, err
	}
	return n, err
}

// 把数据写入设备
func (sd *genericSnmpDevice) OnWrite(cmd []byte, _ []byte) (int, error) {
	return 0, nil
}

// 设备当前状态
func (sd *genericSnmpDevice) Status() typex.DeviceState {
	return typex.DEV_UP
}

// 停止设备
func (sd *genericSnmpDevice) Stop() {
	sd.status = typex.DEV_STOP
	if sd.CancelCTX != nil {
		sd.CancelCTX()
	}
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

func (sd *genericSnmpDevice) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}

func (sd *genericSnmpDevice) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	return []byte{}, nil
}
