// 抽象设备：
// 1.0 以后的大功能：支持抽象设备，抽象设备就是外挂的设备，Rulex本来是个规则引擎，但是1.0之前的版本没有对硬件设备进行抽象支持
// 因此，1.0以后增加对硬件的抽象
// Target Source 描述了数据的流向，抽象设备描述了数据的载体。
// 举例：外挂一个设备，这个设备具备双工控制功能，例如电磁开关等，此时它强调的是设备的物理功能，而数据则不是主体。
// 因此需要抽象出来一个层专门来描述这些设备
package typex

type DeviceType string

func (d DeviceType) String() string {
	return string(d)

}

// 支持的设备类型
const (
	SIEMENS_PLC                DeviceType = "SIEMENS_PLC"                // SIEMENS-S71200
	GENERIC_MODBUS             DeviceType = "GENERIC_MODBUS"             // 通用Modbus
	GENERIC_MODBUS_POINT_EXCEL DeviceType = "GENERIC_MODBUS_POINT_EXCEL" // 通用Modbus通过Excel表配置点位
	GENERIC_UART               DeviceType = "GENERIC_UART"               // 通用串口
	GENERIC_SNMP               DeviceType = "GENERIC_SNMP"               // SNMP 支持
	ICMP_SENDER                DeviceType = "ICMP_SENDER"                // ICMP_SENDER
	GENERIC_PROTOCOL           DeviceType = "GENERIC_PROTOCOL"           // 通用自定义协议处理器
	GENERIC_OPCUA              DeviceType = "GENERIC_OPCUA"              // 通用OPCUA
	GENERIC_CAMERA             DeviceType = "GENERIC_CAMERA"             // 通用摄像头
	GENERIC_AIS_RECEIVER       DeviceType = "GENERIC_AIS_RECEIVER"       // 通用AIS
	GENERIC_BACNET_IP          DeviceType = "GENERIC_BACNET_IP"          // 通用BacnetIP
	RHINOPI_IR                 DeviceType = "RHINOPI_IR"                 // 大犀牛PI的红外线接收器
	GENERIC_HTTP_DEVICE        DeviceType = "GENERIC_HTTP_DEVICE"        // GENERIC_HTTP
	HNC8                       DeviceType = "HNC8"                       // 华中数控机床
	KDN                        DeviceType = "KDN"                        // 凯帝恩控机床
)

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
	// 新特性, 适用于自定义协议读写
	OnCtrl(cmd []byte, args []byte) ([]byte, error)
	// 设备当前状态
	Status() DeviceState
	// 停止设备, 在这里释放资源,一般是先置状态为STOP,然后CancelContext()
	Stop()
	//
	// 0.5.2 新增 Reload() error
	//
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

/*
*
* 子设备网络拓扑[2023-04-17新增]
*
 */
type DeviceTopology struct {
	Id       string                 // 子设备的ID
	Name     string                 // 子设备名
	LinkType int                    // 物理连接方式: 0-ETH 1-WIFI 3-BLE 4 LORA 5 OTHER
	State    int                    // 状态: 0-Down 1-Working
	Info     map[string]interface{} // 子设备的一些额外信息
}
