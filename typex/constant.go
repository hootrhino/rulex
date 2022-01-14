package typex

// Resource State
type ResourceState int

const (
	DOWN  ResourceState = 0
	UP    ResourceState = 1
	PAUSE ResourceState = 2
)

//
// Abstract driver interface
//
type DriverState int

const (
	STOP    DriverState = 0
	RUNNING DriverState = 1
)

//
// InEndType
//
type InEndType string

func (i InEndType) String() string {
	return string(i)
}

const (
	MQTT        InEndType = "MQTT"
	HTTP        InEndType = "HTTP"
	UDP         InEndType = "UDP"
	COAP        InEndType = "COAP"
	GRPC        InEndType = "GRPC"
	UART_MODULE InEndType = "UART_MODULE"
	//
	// MODBUS_MASTER
	//
	MODBUS_MASTER InEndType = "MODBUS_MASTER"
	//
	// MODBUS_SLAVER
	//
	MODBUS_SLAVER InEndType = "MODBUS_SLAVER"
	//
	// From snmp server provider
	//
	SNMP_SERVER InEndType = "SNMP_SERVER"
	//
	// NATS.IO SERVER
	//
	NATS_SERVER InEndType = "NATS_SERVER"
	//
	// 西门子S7客户端
	//
	SIEMENS_S7 InEndType = "SIEMENS_S7"
	//
	// RULEX UDP 自定义简单协议
	//
	RULEX_UDP InEndType = "RULEX_UDP"
)

//
// TargetType
//
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
	MONGO_SINGLE          TargetType = "MONGO_SINGLE"
	MONGO_CLUSTER         TargetType = "MONGO_CLUSTER"
	REDIS_SINGLE          TargetType = "REDIS_SINGLE"
	FLINK_SINGLE          TargetType = "FLINK_SINGLE"
	MQTT_TARGET           TargetType = "MQTT"
	MQTT_TELEMETRY_TARGET TargetType = "MQTT_TELEMETRY_TARGET"
	MYSQL_TARGET          TargetType = "MYSQL"
	PGSQL_TARGET          TargetType = "PGSQL"
	NATS_TARGET           TargetType = "NATS"
	HTTP_TARGET           TargetType = "HTTP"
)
