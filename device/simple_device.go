package device

import (
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
)

type simpleDevice struct {
	typex.XStatus
}

func NewSimpleDevice(deviceId string, e typex.RuleX) typex.XDevice {
	sd := &simpleDevice{}
	sd.RuleEngine = e
	return sd
}

//  初始化
func (d *simpleDevice) Init(devId string, config map[string]interface{}) error {
	glogger.GLogger.Info("simpleDevice Init")
	d.PointId = devId
	return nil
}

// 启动
func (d *simpleDevice) Start(_ typex.CCTX) error {
	glogger.GLogger.Info("simpleDevice Start")
	return nil
}

// 从设备里面读数据出来
func (d *simpleDevice) OnRead(_ []byte) (int, error) {
	glogger.GLogger.Info("simpleDevice Read")
	return 0, nil
}

// 把数据写入设备
func (d *simpleDevice) OnWrite(_ []byte) (int, error) {
	glogger.GLogger.Info("simpleDevice Write")
	return 0, nil
}

// 设备当前状态
func (d *simpleDevice) Status() typex.DeviceState {
	glogger.GLogger.Info("simpleDevice State")
	return typex.DEV_RUNNING
}

// 设备属性，是一系列属性描述
func (d *simpleDevice) Property() []typex.DeviceProperty {
	glogger.GLogger.Info("simpleDevice Property")
	return []typex.DeviceProperty{}
}

//
func (d *simpleDevice) Details() *typex.Device {
	glogger.GLogger.Info("simpleDevice Details")
	return d.RuleEngine.GetDevice(d.PointId)
}
func (d *simpleDevice) SetState(typex.DeviceState) {
	glogger.GLogger.Info("simpleDevice SetState")
}
func (d *simpleDevice) Stop() {
	glogger.GLogger.Info("simpleDevice Stop")
	d.CancelCTX()
}
func (d *simpleDevice) Driver() typex.XExternalDriver {

	return nil
}
