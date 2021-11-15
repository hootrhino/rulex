package typex

// Resource State
type ResourceState int

const (
	DOWN  ResourceState = 0
	UP    ResourceState = 1
	PAUSE ResourceState = 2
)

//
// Rule type is for property store,
// XResource implements struct type is actually worker
//
type ModelType int

// 'T' means Type
const (
	T_NUMBER  ModelType = 1
	T_STRING  ModelType = 2
	T_BOOLEAN ModelType = 3
	T_JSON    ModelType = 4
	T_BIN     ModelType = 5
	T_RAW     ModelType = 6
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
	MQTT              InEndType = "MQTT"
	HTTP              InEndType = "HTTP"
	UDP               InEndType = "UDP"
	COAP              InEndType = "COAP"
	GRPC              InEndType = "GRPC"
	UART_MODULE       InEndType = "UART_MODULE"
	MODBUS_TCP_MASTER InEndType = "MODBUS_TCP_MASTER"
	MODBUS_RTU_MASTER InEndType = "MODBUS_RTU_MASTER"
	MODBUS_TCP_SLAVER InEndType = "MODBUS_TCP_SLAVER"
	MODBUS_RTU_SLAVER InEndType = "MODBUS_RTU_SLAVER"
	//
	// From snmp server provider
	//
	SNMP_SERVER InEndType = "SNMP_SERVER"
)

//
// TargetType
//
type TargetType string

func (i TargetType) String() string {
	return string(i)
}

const (
	MONGO_SINGLE  TargetType = "MONGO_SINGLE"
	MONGO_CLUSTER TargetType = "MONGO_CLUSTER"
	REDIS_SINGLE  TargetType = "REDIS_SINGLE"
	FLINK_SINGLE  TargetType = "FLINK_SINGLE"
	MQTT_TARGET   TargetType = "MQTT"
	MYSQL_TARGET  TargetType = "MYSQL"
	PGSQL_TARGET  TargetType = "PGSQL"
)

//
//
// 创建资源的时候需要一个通用配置类, XConfig 可认为是接收参数的Form
// from v0.0.2
//
//
type XConfig struct {
	Name        string
	Type        string
	Config      map[string]interface{}
	DataModels  map[string]XDataModel
	Description string
}
