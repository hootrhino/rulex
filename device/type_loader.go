package device

import (
	"github.com/hootrhino/rulex/core"
	"github.com/hootrhino/rulex/typex"
)

var DM typex.DeviceRegistry

/*
*
* 加载系统内支持的设备类型
*
 */
func LoadDt() {
	DM = core.NewDeviceTypeManager()
	DM.Register(typex.TSS200V02, &typex.XConfig{})
	DM.Register(typex.RTU485_THER, &typex.XConfig{})
	DM.Register(typex.YK08_RELAY, &typex.XConfig{})
	DM.Register(typex.S1200PLC, &typex.XConfig{})
	DM.Register(typex.GENERIC_MODBUS, &typex.XConfig{})
	DM.Register(typex.GENERIC_MODBUS_POINT_EXCEL, &typex.XConfig{})
	DM.Register(typex.GENERIC_UART, &typex.XConfig{})
	DM.Register(typex.GENERIC_SNMP, &typex.XConfig{})
	DM.Register(typex.USER_G776, &typex.XConfig{})
	DM.Register(typex.ICMP_SENDER, &typex.XConfig{})
}
