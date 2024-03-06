package typex

// TargetType
type TargetType string

func (i TargetType) String() string {
	return string(i)
}

/*
*
* 输出资源类型
*
 */
const (
	MONGO_SINGLE  TargetType = "MONGO_SINGLE"
	MONGO_CLUSTER TargetType = "MONGO_CLUSTER"
	REDIS_SINGLE  TargetType = "REDIS_SINGLE"
	FLINK_SINGLE  TargetType = "FLINK_SINGLE"
	MQTT_TARGET   TargetType = "MQTT"
	MYSQL_TARGET  TargetType = "MYSQL"
	PGSQL_TARGET  TargetType = "PGSQL"
	NATS_TARGET   TargetType = "NATS"
	HTTP_TARGET   TargetType = "HTTP"
	// TDENGINE
	TDENGINE_TARGET TargetType = "TDENGINE"
	// GRPC
	GRPC_CODEC_TARGET TargetType = "GRPC_CODEC_TARGET"
	// UDP Server
	UDP_TARGET TargetType = "UDP_TARGET"
	// SQLITE
	SQLITE_TARGET TargetType = "SQLITE_TARGET"
	// USER_G776 DTU
	USER_G776_TARGET TargetType = "USER_G776_TARGET"
	// TCP 透传
	TCP_TRANSPORT TargetType = "TCP_TRANSPORT"
)

// Stream from source and to target
type XTarget interface {
	//
	// 用来初始化传递资源配置
	//
	Init(outEndId string, configMap map[string]interface{}) error
	//
	// 启动资源
	//
	Start(CCTX) error
	//
	// 获取资源状态
	//
	Status() SourceState
	//
	// 获取资源绑定的的详情
	//
	Details() *OutEnd
	//
	// 数据出口
	//
	To(data interface{}) (interface{}, error)
	//
	// 停止资源, 用来释放资源
	//
	Stop()
}
