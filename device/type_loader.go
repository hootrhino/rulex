package device

import (
	"github.com/i4de/rulex/common"
	"github.com/i4de/rulex/core"
	"github.com/i4de/rulex/typex"
)

var DM typex.DeviceRegistry

/*
*
* 加载系统内支持的设备类型
*
 */
func LoadDt() {
	DM = core.NewDeviceTypeManager()
	DM.Register(typex.TSS200V02, core.GenInConfig(typex.COAP, "About TSS200V02", common.ModBusConfig{}))
	DM.Register(typex.RTU485_THER, core.GenInConfig(typex.COAP, "About RTU485_THER", common.ModBusConfig{}))
	DM.Register(typex.YK08_RELAY, core.GenInConfig(typex.COAP, "About YK08_RELAY", common.ModBusConfig{}))
	DM.Register(typex.S1200PLC, core.GenInConfig(typex.COAP, "About S1200PLC", common.ModBusConfig{}))
	DM.Register(typex.GENERIC_MODBUS, core.GenInConfig(typex.COAP, "About GENERIC_MODBUS", common.ModBusConfig{}))
	DM.Register(typex.GENERIC_UART, core.GenInConfig(typex.COAP, "About GENERIC_UART", common.GenericUartConfig{}))
	DM.Register(typex.GENERIC_SNMP, core.GenInConfig(typex.COAP, "About GENERIC_SNMP", common.GenericSnmpConfig{}))
	DM.Register(typex.USER_G776, core.GenInConfig(typex.COAP, "About USER_G776", common.ModBusConfig{}))
	DM.Register(typex.ICMP_SENDER, core.GenInConfig(typex.COAP, "About ICMP_SENDER", common.HostConfig{}))
}
