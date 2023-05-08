package source

import (
	"github.com/hootrhino/rulex/core"
	"github.com/hootrhino/rulex/typex"
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
	SM.Register(typex.COAP, &typex.XConfig{})
	SM.Register(typex.GRPC, &typex.XConfig{})
	SM.Register(typex.HTTP, &typex.XConfig{})
	SM.Register(typex.RULEX_UDP, &typex.XConfig{})
	SM.Register(typex.NATS_SERVER, &typex.XConfig{})
	SM.Register(typex.MQTT, &typex.XConfig{})
}
