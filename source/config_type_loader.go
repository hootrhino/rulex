package source

import (
	"rulex/core"
	"rulex/typex"
)

var RM typex.SourceRegistry = core.NewSourceTypeManager()

/*
*
* 给前端返回资源类型，这里是个蹩脚的设计
* 以实现功能为主，后续某个重构版本会做的优雅点
*
 */

func LoadRt() {
	RM.Register(typex.COAP, core.GenInConfig(typex.COAP, "About COAP", coAPConfig{}))
	RM.Register(typex.GRPC, core.GenInConfig(typex.GRPC, "About GRPC", grpcConfig{}))
	RM.Register(typex.HTTP, core.GenInConfig(typex.HTTP, "About HTTP", httpConfig{}))
	RM.Register(typex.MODBUS_MASTER, core.GenInConfig(typex.MODBUS_MASTER, "About MODBUS_MASTER", modBusConfig{}))
	RM.Register(typex.MQTT, core.GenInConfig(typex.MQTT, "About MQTT", natsConfig{}))
	RM.Register(typex.NATS_SERVER, core.GenInConfig(typex.NATS_SERVER, "About NATS_SERVER", snmpConfig{}))
	RM.Register(typex.SNMP_SERVER, core.GenInConfig(typex.SNMP_SERVER, "About SNMP_SERVER", siemensS7config{}))
	RM.Register(typex.UART_MODULE, core.GenInConfig(typex.UART_MODULE, "About UART_MODULE", uartConfig{}))
	RM.Register(typex.RULEX_UDP, core.GenInConfig(typex.RULEX_UDP, "About RULEX_UDP", udpConfig{}))
}
