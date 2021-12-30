package typex

//
// 外挂驱动, 比如串口, PLC等, 驱动可以挂在输入或者输出资源上。
// 典型案例:
// 1. MODBUS TCP模式 ,数据输入后转JSON输出到串口屏幕上
// 2. MODBUS TCP模式外挂了很多继电器,来自云端的 PLC 控制指令先到网关, 然后网关决定推送到哪个外挂
//
type DriverDetail struct {
	Name        string `json:"name" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Description string `json:"description" binding:"required"`
}
type XExternalDriver interface {
	Test() error
	Init() error
	Work() error
	State() DriverState
	SetState(DriverState)
	//---------------------------------------------------
	// 读写接口是给LUA标准库用的, 驱动只管实现读写逻辑即可
	//---------------------------------------------------
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	//---------------------------------------------------
	DriverDetail() *DriverDetail
	Stop() error
}
