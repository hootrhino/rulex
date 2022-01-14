package resource

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"rulex/core"
	"rulex/driver"
	"rulex/typex"
	"rulex/utils"
	"time"

	"github.com/ngaut/log"

	"github.com/goburrow/modbus"
	"github.com/mitchellh/mapstructure"
)

type ModBusConfig struct {
	Mode           string          `json:"mode" title:"工作模式" info:"可以在 RTU/TCP 两个模式之间切换"`
	Timeout        int             `json:"timeout" validate:"required" title:"连接超时" info:""`
	SlaverId       byte            `json:"slaverId" validate:"required" title:"TCP端口" info:""`
	Frequency      int64           `json:"frequency" validate:"required" title:"采集频率" info:""`
	Config         interface{}     `json:"config" validate:"required" title:"工作模式配置" info:""`
	RegisterParams []registerParam `json:"registerParams" validate:"required" title:"寄存器配置" info:""`
}

const (
	//-------------------------------------------
	// 	Code |  Register Type
	//-------|------------------------------------
	// 	1	 |	Read Coil
	// 	2	 |	Read Discrete Input
	// 	3	 |	Read Holding Registers
	// 	4	 |	Read Input Registers
	// 	5	 |	Write Single Coil
	// 	6	 |	Write Single Holding Register
	// 	15	 |	Write Multiple Coils
	// 	16	 |	Write Multiple Holding Registers
	//-------------------------------------------

	READ_COIL                        = 1
	READ_DISCRETE_INPUT              = 2
	READ_HOLDING_REGISTERS           = 3
	READ_INPUT_REGISTERS             = 4
	WRITE_SINGLE_COIL                = 5
	WRITE_SINGLE_HOLDING_REGISTER    = 6
	WRITE_MULTIPLE_COILS             = 15
	WRITE_MULTIPLE_HOLDING_REGISTERS = 16
)

/*
*
* coilParams 1
*
 */
type coilsW struct {
	Address  uint16 `json:"address" validate:"required" title:"寄存器地址" info:""`
	Quantity uint16 `json:"quantity" validate:"required" title:"写入数量" info:""`
	Values   []byte `json:"values" validate:"required" title:"写入的值" info:""`
}

/*
*
* 2
*
 */
type coilW struct {
	Address  uint16 `json:"address" validate:"required" title:"寄存器地址" info:""`
	Quantity uint16 `json:"quantity" validate:"required" title:"写入数量" info:""`
	Value    uint16 `json:"value" validate:"required" title:"写入的值" info:""`
}

/*
*
* registerParams 3
*
 */
type registerW struct {
	Address  uint16 `json:"address" validate:"required" title:"寄存器地址" info:""`
	Quantity uint16 `json:"quantity" validate:"required" title:"写入数量" info:""`
	Value    uint16 `json:"value" validate:"required" title:"写入的值" info:""`
}

/*
*
* 4
*
 */
type registersW struct {
	Address  uint16 `json:"address" validate:"required" title:"寄存器地址" info:""`
	Quantity uint16 `json:"quantity" validate:"required" title:"写入数量" info:""`
	Values   []byte `json:"values" validate:"required" title:"写入的值" info:""`
}

/*
*
* 配置进来准备采集的寄存器参数
*
 */

type registerParam struct {
	Tag      string `json:"tag" validate:"required"`      // Function
	Function int    `json:"function" validate:"required"` // Function
	Address  uint16 `json:"address" validate:"required"`  // Address
	Quantity uint16 `json:"quantity" validate:"required"` // Quantity
}

/*
*
* 采集到的数据
*
 */
type registerData struct {
	Tag      string `json:"tag" validate:"required"`      // Function
	Function int    `json:"function" validate:"required"` // Function
	Address  uint16 `json:"address" validate:"required"`  // Address
	Quantity uint16 `json:"quantity" validate:"required"` // Quantity
	Value    string `json:"value" validate:"required"`    // Quantity
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
type RtuConfig struct {
	Uart     string `json:"uart" validate:"required" title:"串口路径" info:"本地系统的串口路径"`
	BaudRate int    `json:"baudRate" validate:"required" title:"波特率" info:"串口通信波特率"`
	DataBits int    `json:"dataBits" validate:"required" title:"数据位" info:"串口通信数据位"`
	Parity   string `json:"parity" validate:"required" title:"分割位" info:"串口通信分割位"`
	StopBits int    `json:"stopBits" validate:"required" title:"停止位" info:"串口通信停止位"`
}

//
//
//
type TcpConfig struct {
	Ip   string `json:"ip" validate:"required" title:"IP地址" info:""`
	Port int    `json:"port" validate:"required" title:"端口" info:""`
}

//
//
//---------------------------------------------------------------------------

type ModbusMasterResource struct {
	typex.XStatus
	client    modbus.Client
	cxt       context.Context
	rtuDriver typex.XExternalDriver
}

func NewModbusMasterResource(id string, e typex.RuleX) typex.XResource {
	m := ModbusMasterResource{}
	m.RuleEngine = e
	m.cxt = context.Background()
	return &m
}
func (*ModbusMasterResource) Configs() typex.XConfig {
	config, err := core.RenderConfig("MODBUS_MASTER", "", ModBusConfig{})
	if err != nil {
		log.Error(err)
		return typex.XConfig{}
	} else {
		return config
	}
}

func (m *ModbusMasterResource) Register(inEndId string) error {
	m.PointId = inEndId
	return nil
}

func (m *ModbusMasterResource) Start() error {

	config := m.RuleEngine.GetInEnd(m.PointId).Config
	var mainConfig ModBusConfig
	if err := utils.BindResourceConfig(config, &mainConfig); err != nil {
		return err
	}

	if mainConfig.Mode == "TCP" {
		var tcpConfig TcpConfig
		if errs := mapstructure.Decode(mainConfig.Config, &tcpConfig); errs != nil {
			log.Error(errs)
			return errs
		}
		handler := modbus.NewTCPClientHandler(
			fmt.Sprintf("%s:%v", tcpConfig.Ip, tcpConfig.Port),
		)
		handler.Timeout = time.Duration(mainConfig.Frequency) * time.Second
		handler.SlaveId = mainConfig.SlaverId
		if err := handler.Connect(); err != nil {
			return err
		}
		m.client = modbus.NewClient(handler)
	} else if mainConfig.Mode == "RTU" {
		var rtuConfig RtuConfig
		if errs := mapstructure.Decode(mainConfig.Config, &rtuConfig); errs != nil {
			log.Error(errs)
			return errs
		}
		handler := modbus.NewRTUClientHandler(rtuConfig.Uart)
		handler.BaudRate = rtuConfig.BaudRate
		//
		// rtuConfig
		//
		handler.DataBits = rtuConfig.DataBits
		handler.Parity = rtuConfig.Parity
		handler.StopBits = rtuConfig.StopBits
		//---------
		handler.SlaveId = mainConfig.SlaverId
		handler.Timeout = time.Duration(mainConfig.Frequency) * time.Second
		//---------
		if err := handler.Connect(); err != nil {
			return err
		}
		m.client = modbus.NewClient(handler)
		m.rtuDriver = driver.NewModBusRtuDriver(m.Details(), m.RuleEngine, m.client)
	} else {
		return errors.New("no supported mode:" + mainConfig.Mode)
	}
	//---------------------------------------------------------------------------------
	// Start
	//---------------------------------------------------------------------------------

	ticker := time.NewTicker(time.Duration(mainConfig.Frequency) * time.Second)
	//
	// 前端传过来个寄存器和地址的配置列表，然后给每个寄存器分配一个协程去读
	//
	for _, rCfg := range mainConfig.RegisterParams {
		log.Info("Start read register:", rCfg.Address)
		// 每个寄存器配一个协程读数据
		go func(ctx context.Context, rp registerParam) {
			defer ticker.Stop()
			// Modbus data is most often read and written as "registers" which are [16-bit] pieces of data. Most often,
			// the register is either a signed or unsigned 16-bit integer. If a 32-bit integer or floating point is required,
			// these values are actually read as a pair of registers.
			var results []byte
			var err error
			for {
				<-ticker.C
				select {
				case <-ctx.Done():
					{
						return
					}
				default:
					{
						if rp.Function == READ_COIL {
							results, err = m.client.ReadCoils(rp.Address, rp.Quantity)
						}
						if rp.Function == READ_DISCRETE_INPUT {
							results, err = m.client.ReadDiscreteInputs(rp.Address, rp.Quantity)
						}
						if rp.Function == READ_HOLDING_REGISTERS {
							results, err = m.client.ReadHoldingRegisters(rp.Address, rp.Quantity)
						}
						if rp.Function == READ_INPUT_REGISTERS {
							results, err = m.client.ReadInputRegisters(rp.Address, rp.Quantity)
						}
						//
						// error
						//
						if err != nil {
							log.Error("NewModbusMasterResource ReadData error: ", err)
						} else {
							data := registerData{
								Tag:      rp.Tag,
								Function: rp.Function,
								Address:  rp.Address,
								Quantity: rp.Quantity,
								Value:    string(results),
							}
							bytes, _ := json.Marshal(data)
							m.RuleEngine.Work(m.RuleEngine.GetInEnd(m.PointId), string(bytes))
						}
					}
				}
			}
		}(m.cxt, rCfg)
	}
	return nil

}

func (m *ModbusMasterResource) Details() *typex.InEnd {
	return m.RuleEngine.GetInEnd(m.PointId)
}

func (m *ModbusMasterResource) Test(inEndId string) bool {
	return true
}

func (m *ModbusMasterResource) Enabled() bool {
	return m.Enable
}

func (m *ModbusMasterResource) DataModels() []typex.XDataModel {
	return []typex.XDataModel{}
}

func (m *ModbusMasterResource) Reload() {

}

func (m *ModbusMasterResource) Pause() {

}

func (m *ModbusMasterResource) Status() typex.ResourceState {
	return typex.UP

}

func (m *ModbusMasterResource) Stop() {
	m.cxt.Done()
}

/*
*
* 对于 modbus 资源来说, 任何直接写入的数据都被认为是给寄存器写值
*
 */
type dataM struct {
	Type   int                    `json:"type" validate:"required"`
	Params map[string]interface{} `json:"params" validate:"required"`
}

/*
*
* 写入值
*
 */
func (m *ModbusMasterResource) OnStreamApproached(data string) error {
	var dm dataM
	if err := json.Unmarshal([]byte(data), &dm); err != nil {
		log.Error(err)
		return err
	}
	/*
	*
	* 线圈
	*
	 */
	if dm.Type == WRITE_SINGLE_COIL {
		var c coilW
		if errs := mapstructure.Decode(dm.Params, &c); errs != nil {
			log.Error(errs)
			return errs
		}
		if _, errs := m.client.WriteSingleCoil(c.Address, c.Value); errs != nil {
			log.Error(errs)
			return errs
		}
	}
	if dm.Type == WRITE_MULTIPLE_COILS {
		var cs coilsW
		if errs := mapstructure.Decode(dm.Params, &cs); errs != nil {
			log.Error(errs)
			return errs
		}
		if _, errs := m.client.WriteMultipleCoils(cs.Address, cs.Quantity, cs.Values); errs != nil {
			log.Error(errs)
			return errs
		}
	}
	/*
	*
	* 寄存器
	*
	 */
	if dm.Type == WRITE_SINGLE_HOLDING_REGISTER {
		var r registerW
		if errs := mapstructure.Decode(dm.Params, r); errs != nil {
			log.Error(errs)
			return errs
		}
		if _, errs := m.client.WriteSingleRegister(r.Address, r.Value); errs != nil {
			log.Error(errs)
			return errs
		}
	}

	if dm.Type == WRITE_MULTIPLE_HOLDING_REGISTERS {
		var rs registersW
		if errs := mapstructure.Decode(dm.Params, rs); errs != nil {
			log.Error(errs)
			return errs
		}
		if _, errs := m.client.WriteMultipleRegisters(rs.Address, rs.Quantity, rs.Values); errs != nil {
			log.Error(errs)
			return errs
		}
	}
	return nil
}

/*
*
* 只有RTU模式才带驱动
*
 */
func (m *ModbusMasterResource) Driver() typex.XExternalDriver {
	if m.client != nil {
		return m.rtuDriver
	} else {
		return nil
	}
}

//
// 拓扑
//
func (*ModbusMasterResource) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}
