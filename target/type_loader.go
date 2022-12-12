package target

import (
	"github.com/i4de/rulex/common"
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
	TM.Register(typex.HTTP_TARGET, typex.GenOutConfig(typex.HTTP_TARGET, "About HTTP_TARGET", common.HTTPConfig{}))
	TM.Register(typex.MONGO_SINGLE, typex.GenOutConfig(typex.MONGO_SINGLE, "About MONGO_SINGLE", common.MongoConfig{}))
	TM.Register(typex.MQTT_TARGET, typex.GenOutConfig(typex.MQTT_TARGET, "About MQTT_TARGET", common.MqttConfig{}))
	TM.Register(typex.NATS_TARGET, typex.GenOutConfig(typex.NATS_TARGET, "About NATS_TARGET", common.NatsConfig{}))
	TM.Register(typex.TDENGINE_TARGET, typex.GenOutConfig(typex.TDENGINE_TARGET, "About TDENGINE_TARGET", common.TDEngineConfig{}))
}
