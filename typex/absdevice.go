// 抽象设备：
// 1.0 以后的大功能：支持抽象设备，抽象设备就是外挂的设备，Rulex本来是个规则引擎，但是1.0之前的版本没有对硬件设备进行抽象支持
// 因此，1.0以后增加对硬件的抽象
// Target Source 描述了数据的流向，抽象设备描述了数据的载体。
// 举例：外挂一个设备，这个设备具备双工控制功能，例如电磁开关等，此时它强调的是设备的物理功能，而数据则不是主体。
// 因此需要抽象出来一个层专门来描述这些设备
package typex

type DeviceState int

const (
	DEV_STOP    DeviceState = 0
	DEV_RUNNING DeviceState = 1
)

type DeviceInfo struct {
	Name         string // 设备名称，例如：灯光开关
	Type         string // 类型,一般是设备-型号，比如 ARDUINO-R3
	ActionScript string // 当收到指令的时候响应脚本
	Description  string // 设备描述信息
}
type DeviceProperty struct {
	Name  string
	Type  string
	Value interface{}
}
type AbstractDevice interface {
	//  初始化
	Init(config map[string]interface{})
	// 启动
	Start(CCTX) error
	// 从设备里面读数据出来
	Read([]byte) (int, error)
	// 把数据写入设备
	Write([]byte) (int, error)
	// 获取设备信息
	Info() DeviceInfo
	// 设备当前状态
	State() DeviceState
	// 设备属性，是一系列属性描述
	Property() []DeviceProperty
}
