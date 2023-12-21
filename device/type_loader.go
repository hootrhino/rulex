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
	DM.Register(typex.SIEMENS_PLC, &typex.XConfig{})
	DM.Register(typex.GENERIC_MODBUS, &typex.XConfig{})
	DM.Register(typex.GENERIC_MODBUS_POINT_EXCEL, &typex.XConfig{})
	DM.Register(typex.GENERIC_UART, &typex.XConfig{})
	DM.Register(typex.GENERIC_SNMP, &typex.XConfig{})
	DM.Register(typex.ICMP_SENDER, &typex.XConfig{})
}
