package source

import (
	"github.com/i4de/rulex/common"
	"github.com/i4de/rulex/core"
	"github.com/i4de/rulex/typex"
)

var SM typex.SourceRegistry = core.NewSourceTypeManager()

/*
*
* 给前端返回资源类型，这里是个蹩脚的设计
* 以实现功能为主，后续某个重构版本会做的优雅点
*
 */

func LoadSt() {
	SM = core.NewSourceTypeManager()
	SM.Register(typex.COAP, typex.GenInConfig(typex.COAP, "About COAP", common.HostConfig{}))
	SM.Register(typex.GRPC, typex.GenInConfig(typex.GRPC, "About GRPC", common.GrpcConfig{}))
	SM.Register(typex.HTTP, typex.GenInConfig(typex.HTTP, "About HTTP", common.HostConfig{}))
	SM.Register(typex.RULEX_UDP, typex.GenInConfig(typex.RULEX_UDP, "About RULEX_UDP", common.RULEXUdpConfig{}))
	SM.Register(typex.NATS_SERVER, typex.GenInConfig(typex.NATS_SERVER, "About NATS_SERVER", common.NatsConfig{}))
	SM.Register(typex.MQTT, typex.GenInConfig(typex.MQTT, "About MQTT", common.NatsConfig{}))
}
