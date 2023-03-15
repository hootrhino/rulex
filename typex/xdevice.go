// 抽象设备：
// 1.0 以后的大功能：支持抽象设备，抽象设备就是外挂的设备，Rulex本来是个规则引擎，但是1.0之前的版本没有对硬件设备进行抽象支持
// 因此，1.0以后增加对硬件的抽象
// Target Source 描述了数据的流向，抽象设备描述了数据的载体。
// 举例：外挂一个设备，这个设备具备双工控制功能，例如电磁开关等，此时它强调的是设备的物理功能，而数据则不是主体。
// 因此需要抽象出来一个层专门来描述这些设备
package typex

import (
	"github.com/i4de/rulex/utils"
)

type DeviceState int

const (
	// 外部停止
	DEV_STOP DeviceState = 0
	// 设备启用
	DEV_UP DeviceState = 1
	// 设备故障
	DEV_DOWN DeviceState = 2
)

type DeviceType string

// 支持的设备类型
const (
	TSS200V02        DeviceType = "TSS200V02"        // Multi params Sensor
	RTU485_THER      DeviceType = "RTU485_THER"      // RS485 Sensor
	YK08_RELAY       DeviceType = "YK08_RELAY"       // YK8 RS485 Relay
	S1200PLC         DeviceType = "S1200PLC"         // SIEMENS-S71200
	GENERIC_MODBUS   DeviceType = "GENERIC_MODBUS"   // 通用Modbus
	GENERIC_UART     DeviceType = "GENERIC_UART"     // 通用串口
	GENERIC_SNMP     DeviceType = "GENERIC_SNMP"     // SNMP 支持
	USER_G776        DeviceType = "USER_G776"        // 有人 G776 4G模组
	ICMP_SENDER      DeviceType = "ICMP_SENDER"      // ICMP_SENDER
	GENERIC_PROTOCOL DeviceType = "GENERIC_PROTOCOL" // 通用自定义协议处理器

)

// 设备元数据
type Device struct {
	UUID         string                 `json:"uuid"`         // UUID
	Name         string                 `json:"name"`         // 设备名称，例如：灯光开关
	Type         DeviceType             `json:"type"`         // 类型,一般是设备-型号，比如 ARDUINO-R3
	ActionScript string                 `json:"actionScript"` // 当收到指令的时候响应脚本
	Description  string                 `json:"description"`  // 设备描述信息
	BindRules    map[string]Rule        `json:"-"`            // 与之关联的规则
	State        DeviceState            `json:"state"`        // 状态
	Config       map[string]interface{} `json:"config"`       // 配置
	Device       XDevice                `json:"-"`            // 实体设备
}

func NewDevice(t DeviceType,
	name string,
	description string,
	actionScript string,
	config map[string]interface{}) *Device {
	return &Device{
		UUID:        utils.DeviceUuid(),
		Name:        name,
		Type:        t,
		State:       DEV_STOP,
		Description: description,
		Config:      config,
		BindRules:   map[string]Rule{},
	}
}

// 设备的属性，是个描述结构
type DeviceProperty struct {
	Name  string
	Type  string
	Value interface{}
}
type DCAModel struct {
	UUID    string      `json:"uuid"`
	Command string      `json:"command"`
	Args    interface{} `json:"args"`
}
type DCAResult struct {
	Error error
	Data  string
}

// 真实工作设备,即具体实现
type XDevice interface {
	// 初始化 通常用来获取设备的配置
	Init(devId string, configMap map[string]interface{}) error
	// 启动, 设备的工作进程
	Start(CCTX) error
	// 从设备里面读数据出来, 第一个参数一般作flag用, 也就是常说的指令类型
	OnRead(cmd []byte, data []byte) (int, error)
	// 把数据写入设备, 第一个参数一般作flag用, 也就是常说的指令类型
	OnWrite(cmd []byte, data []byte) (int, error)
	// 设备当前状态
	Status() DeviceState
	// 停止设备, 在这里释放资源,一般是先置状态为STOP,然后CancelContext()
	Stop()
	// 设备属性，是一系列属性描述
	Property() []DeviceProperty
	// 链接指向真实设备，保存在内存里面，和SQLite里的数据是对应关系
	Details() *Device
	// 状态
	SetState(DeviceState)
	// 驱动接口, 通常用来和硬件交互
	Driver() XExternalDriver
	// 外部调用, 该接口是个高级功能, 准备为了设计分布式部署设备的时候用, 但是相当长时间内都不会开启
	// 默认情况下该接口没有用
	OnDCACall(UUID string, Command string, Args interface{}) DCAResult
}
