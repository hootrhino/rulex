package target

import (
	"github.com/i4de/rulex/core"
	"github.com/i4de/rulex/typex"
)

var TM typex.TargetRegistry

/*
*
* 给前端返回资源类型，这里是个蹩脚的设计
* 以实现功能为主，后续某个重构版本会做的优雅点
*
 */

func LoadTt() {
	TM = core.NewTargetTypeManager()
	TM.Register(typex.HTTP_TARGET, &typex.XConfig{})
	TM.Register(typex.MONGO_SINGLE, &typex.XConfig{})
	TM.Register(typex.MQTT_TARGET, &typex.XConfig{})
	TM.Register(typex.NATS_TARGET, &typex.XConfig{})
	TM.Register(typex.TDENGINE_TARGET, &typex.XConfig{})
}
