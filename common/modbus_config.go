package common

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
	Address  uint16 `json:"address" validate:"required" title:"寄存器地址"`
	Quantity uint16 `json:"quantity" validate:"required" title:"写入数量"`
	Values   []byte `json:"values" validate:"required" title:"写入的值"` // 如果是单个写 取 Values[0]
}

/*
*
* 4
*
 */
type Registers struct {
	Address  uint16 `json:"address" validate:"required" title:"寄存器地址"`
	Quantity uint16 `json:"quantity" validate:"required" title:"写入数量"`
	Values   []byte `json:"values" validate:"required" title:"写入的值"` // 如果是单个写 取 Values[0]
}

/*
*
* 采集到的数据
*
 */
type RegisterRW struct {
	UUID      string  `json:"UUID"`
	Tag       string  `json:"tag" validate:"required" title:"数据Tag"`         // 数据Tag
	Alias     string  `json:"alias" validate:"required" title:"别名"`          // 别名
	Function  int     `json:"function" validate:"required" title:"Modbus功能"` // Function
	SlaverId  byte    `json:"slaverId" validate:"required" title:"从机ID"`     // 从机ID
	Address   uint16  `json:"address" validate:"required" title:"地址"`        // Address
	Frequency int64   `json:"frequency" validate:"required" title:"采集频率"`    // 间隔
	Quantity  uint16  `json:"quantity" validate:"required" title:"数量"`       // Quantity
	Type      string  `json:"type"`                                          // 运行时数据
	Order     string  `json:"order"`                                         // 运行时数据
	Weight    float64 `json:"weight"`
	Value     string  `json:"value,omitempty"` // 运行时数据. Type, Order不同值也不同
}

type RegisterList []*RegisterRW

func (r RegisterList) Len() int {
	return len(r)
}

func (r RegisterList) Less(i, j int) bool {
	if r[i].SlaverId == r[j].SlaverId {
		if r[i].Function == r[j].Function {
			if r[i].Frequency == r[j].Frequency {
				return r[i].Address < r[j].Address
			} else {
				return r[i].Frequency < r[j].Frequency
			}
		} else {
			return r[i].Function < r[j].Function
		}
	} else {
		return r[i].SlaverId < r[j].SlaverId
	}
}

func (r RegisterList) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

/*
*
* 写入的数据
*
 */
type RegisterW struct {
	Function int    `json:"function"` // Function
	SlaverId byte   `json:"slaverId"` // 从机ID
	Address  uint16 `json:"address"`  // Address
	Quantity uint16 `json:"quantity"` // Quantity
	Values   []byte `json:"values"`   // Value
}

type ModBusConfig struct {
	Mode        string       `json:"mode" title:"工作模式" info:"UART/TCP"`
	Timeout     *int         `json:"timeout" validate:"required" title:"连接超时"`
	AutoRequest *bool        `json:"autoRequest" validate:"required"`
	Frequency   *int64       `json:"frequency" validate:"required" title:"采集频率"`
	Config      interface{}  `json:"config" validate:"required" title:"工作模式"`
	Registers   []RegisterRW `json:"registers" validate:"required" title:"寄存器配置"`
}
