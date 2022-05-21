package device

import (
	"rulex/typex"

	"github.com/ngaut/log"
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
	log.Info("simpleDevice Init")
	d.PointId = devId
	return nil
}

// 启动
func (d *simpleDevice) Start(_ typex.CCTX) error {
	log.Info("simpleDevice Start")
	return nil
}

// 从设备里面读数据出来
func (d *simpleDevice) Read(_ []byte) (int, error) {
	log.Info("simpleDevice Read")
	return 0, nil
}

// 把数据写入设备
func (d *simpleDevice) Write(_ []byte) (int, error) {
	log.Info("simpleDevice Write")
	return 0, nil
}

// 设备当前状态
func (d *simpleDevice) Status() typex.DeviceState {
	log.Info("simpleDevice State")
	return typex.DEV_RUNNING
}

// 设备属性，是一系列属性描述
func (d *simpleDevice) Property() []typex.DeviceProperty {
	log.Info("simpleDevice Property")
	return []typex.DeviceProperty{}
}

//
func (d *simpleDevice) Details() *typex.Device {
	log.Info("simpleDevice Details")
	return d.RuleEngine.GetDevice(d.PointId)
}
func (d *simpleDevice) SetState(typex.DeviceState) {
	log.Info("simpleDevice SetState")
}
func (d *simpleDevice) Stop() {
	log.Info("simpleDevice Stop")
	d.CancelCTX()
}
