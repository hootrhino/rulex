package typex

//
// 外挂驱动, 比如串口, PLC等, 驱动可以挂在输入或者输出资源上。
// 典型案例:
// 1. MODBUS TCP模式 ,数据输入后转JSON输出到串口屏幕上
// 2. MODBUS TCP模式外挂了很多继电器,来自云端的 PLC 控制指令先到网关, 然后网关决定推送到哪个外挂
//
type DriverDetail struct {
	UUID        string `json:"uuid" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Description string `json:"description" binding:"required"`
}

//
// 驱动由源(Source)或者设备(Device)启动,驱动的状态(Status)被RULEX获取,或者被源或者设备获取后返回给RULEX
//
type XExternalDriver interface {
	Test() error
	Init(map[string]string) error
	Work() error
	State() DriverState
	Read(cmd []byte, data []byte) (int, error)
	Write(cmd []byte, data []byte) (int, error)
	DriverDetail() DriverDetail
	Stop() error
}
