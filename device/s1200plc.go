package device

import (
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
)

type s1200plc struct {
	typex.XStatus
	status     typex.DeviceState
	RuleEngine typex.RuleX
	driver     typex.XExternalDriver
}

/*
*
* 8路继电器
*
 */
func NewS1200plc(deviceId string, e typex.RuleX) typex.XDevice {
	yk8 := new(s1200plc)
	yk8.PointId = deviceId
	yk8.RuleEngine = e
	return yk8
}

//  初始化
func (yk8 *s1200plc) Init(devId string, config map[string]interface{}) error {

	return nil
}

// 启动
func (yk8 *s1200plc) Start(cctx typex.CCTX) error {
	yk8.Ctx = cctx.Ctx
	yk8.CancelCTX = cctx.CancelCTX
	return nil
}

// 从设备里面读数据出来
func (yk8 *s1200plc) OnRead(data []byte) (int, error) {

	n, err := yk8.driver.Read(data)
	if err != nil {
		glogger.GLogger.Error(err)
		yk8.status = typex.DEV_STOP
	}
	return n, err
}

// 把数据写入设备
func (yk8 *s1200plc) OnWrite(b []byte) (int, error) {
	return yk8.driver.Write(b)
}

// 设备当前状态
func (yk8 *s1200plc) Status() typex.DeviceState {
	return typex.DEV_RUNNING
}

// 停止设备
func (yk8 *s1200plc) Stop() {
	yk8.CancelCTX()
}

// 设备属性，是一系列属性描述
func (yk8 *s1200plc) Property() []typex.DeviceProperty {
	return []typex.DeviceProperty{}
}

// 真实设备
func (yk8 *s1200plc) Details() *typex.Device {
	return yk8.RuleEngine.GetDevice(yk8.PointId)
}

// 状态
func (yk8 *s1200plc) SetState(status typex.DeviceState) {
	yk8.status = status

}

// 驱动
func (yk8 *s1200plc) Driver() typex.XExternalDriver {
	return yk8.driver
}
