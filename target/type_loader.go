package target

import (
	"github.com/i4de/rulex/core"
	"github.com/i4de/rulex/typex"
)

var TM typex.TargetRegistry = core.NewTargetTypeManager()

/*
*
* 给前端返回资源类型，这里是个蹩脚的设计
* 以实现功能为主，后续某个重构版本会做的优雅点
*
 */

func LoadTt() {
	TM.Register(typex.HTTP_TARGET, core.GenOutConfig(typex.HTTP_TARGET, "About HTTP_TARGET", httpConfig{}))
	TM.Register(typex.MONGO_SINGLE, core.GenOutConfig(typex.MONGO_SINGLE, "About MONGO_SINGLE", mongoConfig{}))
	TM.Register(typex.MQTT_TARGET, core.GenOutConfig(typex.MQTT_TARGET, "About MQTT_TARGET", mqttConfig{}))
	TM.Register(typex.NATS_TARGET, core.GenOutConfig(typex.NATS_TARGET, "About NATS_TARGET", natsConfig{}))
	TM.Register(typex.TDENGINE_TARGET, core.GenOutConfig(typex.TDENGINE_TARGET, "About TDENGINE_TARGET", tdEngineConfig{}))
}
