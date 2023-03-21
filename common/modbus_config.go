package common

type ModBusConfig struct {
	Mode     string `json:"mode" title:"工作模式" info:"RTU/TCP"`
	Timeout  int    `json:"timeout" validate:"required" title:"连接超时" info:""`
	WaitTime int    `json:"waitTime" title:"轮询等待时间, 默认100ms" info:""`
	// Weather allow AutoRequest?
	AutoRequest bool `json:"autoRequest" title:"启动轮询" info:""`
	// Request Frequency, default 5 second
	Frequency int64        `json:"frequency" validate:"required" title:"采集频率" info:""`
	Config    interface{}  `json:"config" validate:"required" title:"工作模式" info:""`
	Registers []RegisterRW `json:"registers" validate:"required" title:"寄存器配置" info:""`
}

const (
	READ_COIL                        = 1  //  Read Coil
	READ_DISCRETE_INPUT              = 2  //  Read Discrete Input
	READ_HOLDING_REGISTERS           = 3  //  Read Holding Registers
	READ_INPUT_REGISTERS             = 4  //  Read Input Registers
	WRITE_SINGLE_COIL                = 5  //  Write Single Coil
	WRITE_SINGLE_HOLDING_REGISTER    = 6  //  Write Single Holding Register
	WRITE_MULTIPLE_COILS             = 15 //  Write Multiple Coils
	WRITE_MULTIPLE_HOLDING_REGISTERS = 16 //  Write Multiple Holding Registers
)

/*
*
* coilParams 1
*
 */
type Coils struct {
	Address  uint16 `json:"address" validate:"required" title:"寄存器地址" info:""`
	Quantity uint16 `json:"quantity" validate:"required" title:"写入数量" info:""`
	Values   []byte `json:"values" validate:"required" title:"写入的值" info:""` // 如果是单个写 取 Values[0]
}

/*
*
* 4
*
 */
type Registers struct {
	Address  uint16 `json:"address" validate:"required" title:"寄存器地址" info:""`
	Quantity uint16 `json:"quantity" validate:"required" title:"写入数量" info:""`
	Values   []byte `json:"values" validate:"required" title:"写入的值" info:""` // 如果是单个写 取 Values[0]
}

/*
*
* 采集到的数据
*
 */
type RegisterRW struct {
	Tag      string `json:"tag" validate:"required" title:"数据Tag" info:""`         // Function
	Function int    `json:"function" validate:"required" title:"Modbus功能" info:""` // Function
	SlaverId byte   `json:"slaverId" validate:"required" title:"从机ID" info:""`     // 从机ID
	Address  uint16 `json:"address" validate:"required" title:"地址" info:""`        // Address
	Quantity uint16 `json:"quantity" validate:"required" title:"数量" info:""`       // Quantity
	Value    string `json:"value" title:"值" info:"本地系统的串口路径"`                      // Value
}

//
// Uart "/dev/ttyUSB0"
// BaudRate = 115200
// DataBits = 8
// Parity = "N"
// StopBits = 1
// SlaveId = 1
// Timeout = 5 * time.Second
//
type RTUConfig struct {
	Uart     string `json:"uart" validate:"required" title:"串口路径" info:"本地系统的串口路径"`
	BaudRate int    `json:"baudRate" validate:"required" title:"波特率" info:"串口通信波特率"`
	DataBits int    `json:"dataBits" validate:"required" title:"数据位" info:"串口通信数据位"`
	Parity   string `json:"parity" validate:"required" title:"奇偶校验" info:"奇偶校验"`
	StopBits int    `json:"stopBits" validate:"required" title:"停止位" info:"串口通信停止位"`
}
