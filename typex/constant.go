package typex

// Source State
type SourceState int

const (
	SOURCE_DOWN  SourceState = 0 // 此状态需要重启
	SOURCE_UP    SourceState = 1
	SOURCE_PAUSE SourceState = 2
	SOURCE_STOP  SourceState = 3
)

func (s SourceState) String() string {
	if s == 0 {
		return "DOWN"
	}
	if s == 1 {
		return "UP"
	}
	if s == 2 {
		return "PAUSE"
	}
	if s == 3 {
		return "STOP"
	}
	return "UnKnown State"

}

// Abstract driver interface
type DriverState int

const (
	// STOP 状态一般用来直接停止一个资源，监听器不需要重启
	DRIVER_STOP DriverState = 0
	// UP 工作态
	DRIVER_UP DriverState = 1
	// DOWN 状态是某个资源挂了，属于工作意外，需要重启
	DRIVER_DOWN DriverState = 2
)

// InEndType
type InEndType string

func (i InEndType) String() string {
	return string(i)
}

const (
	MQTT            InEndType = "MQTT"
	HTTP            InEndType = "HTTP"
	COAP            InEndType = "COAP"
	GRPC            InEndType = "GRPC"
	NATS_SERVER     InEndType = "NATS_SERVER"
	RULEX_UDP       InEndType = "RULEX_UDP"
	GENERIC_IOT_HUB InEndType = "GENERIC_IOT_HUB"
	INTERNAL_EVENT  InEndType = "INTERNAL_EVENT" // 内部消息
	GENERIC_MQTT    InEndType = "GENERIC_MQTT"   // 通用MQTT
)

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

/*
*
* 串口校验形式
*
 */
type Parity string

const (
	ODD  Parity = "O" // 奇校验
	EVEN Parity = "E" // 偶校验
	NONE Parity = "N" // 不校验
)
