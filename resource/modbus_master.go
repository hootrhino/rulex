package resource

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"rulex/driver"
	"rulex/typex"
	"rulex/utils"
	"time"

	"github.com/ngaut/log"

	"github.com/goburrow/modbus"
)

type ModBusConfig struct {
	Mode           string          `json:"mode"`
	Timeout        int             `json:"timeout" validate:"required,gte=1,lte=60"`
	SlaverId       byte            `json:"slaverId" validate:"required,gte=1,lte=255"`
	Frequency      int64           `json:"frequency" validate:"required,gte=1,lte=10000"`
	RtuConfig      RtuConfig       `json:"rtuConfig" validate:"required"`
	TcpConfig      TcpConfig       `json:"tcpConfig" validate:"required"`
	RegisterParams []registerParam `json:"registerParams" validate:"required"`
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

type modBUSWriteParams struct {
	Function int    `json:"function" validate:"required"`
	Address  uint16 `json:"address" validate:"required"`
	Quantity uint16 `json:"quantity" validate:"required"`
	Value    []byte `json:"value" validate:"required"`
	Values   []byte `json:"values" validate:"required"`
}

type registerParam struct {
	Function int    `json:"function" validate:"required"` // Function
	Address  uint16 `json:"address" validate:"required"`  // Address
	Quantity uint16 `json:"quantity" validate:"required"` // Quantity
}
type registerData struct {
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
	Uart     string `json:"uart" validate:"required"`
	BaudRate int    `json:"baudRate" validate:"required"`
}

//
//
//
type TcpConfig struct {
	Ip   string `json:"ip" validate:"required"`
	Port int    `json:"port" validate:"required,gte=1,lte=65535"`
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
func (*ModbusMasterResource) Configs() []typex.XConfig {
	return []typex.XConfig{}
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
		handler := modbus.NewTCPClientHandler(
			fmt.Sprintf("%s:%v", mainConfig.TcpConfig.Ip, mainConfig.TcpConfig.Port),
		)
		handler.Timeout = time.Duration(mainConfig.Frequency) * time.Second
		handler.SlaveId = mainConfig.SlaverId
		if err := handler.Connect(); err != nil {
			return err
		}
		m.client = modbus.NewClient(handler)
	} else if mainConfig.Mode == "RTU" {
		handler := modbus.NewRTUClientHandler(mainConfig.RtuConfig.Uart)
		handler.BaudRate = mainConfig.RtuConfig.BaudRate
		// Use default uart config
		handler.DataBits = 8
		handler.Parity = "N"
		handler.StopBits = 1
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
	for _, rCfg := range mainConfig.RegisterParams {
		log.Info("Start read register:", rCfg.Address)

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
* 对于 modbus 资源来说，任何直接写入的数据都被认为是给寄存器写值
*
 */
func (m *ModbusMasterResource) OnStreamApproached(data string) error {
	var p modBUSWriteParams
	var errs error = nil
	if errs := utils.TransformConfig([]byte(data), p); errs != nil {
		log.Error(errs)
		return errs
	}
	if p.Function == WRITE_SINGLE_COIL {
		_, errs = m.client.WriteSingleCoil(1, 1)
	}
	if p.Function == WRITE_SINGLE_HOLDING_REGISTER {
		_, errs = m.client.WriteSingleRegister(1, 1)
	}
	if p.Function == WRITE_MULTIPLE_COILS {
		_, errs = m.client.WriteMultipleCoils(p.Address, p.Quantity, p.Value)
	}
	if p.Function == WRITE_MULTIPLE_HOLDING_REGISTERS {
		_, errs = m.client.WriteMultipleRegisters(p.Address, p.Quantity, p.Values)
	}
	if errs != nil {
		log.Error(errs)
	}
	return errs
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
