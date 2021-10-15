package resource

import (
	"context"
	"encoding/json"
	"fmt"

	"rulex/typex"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/ngaut/log"

	"github.com/goburrow/modbus"
)

type RegisterParam struct {
	// 	Code |  Register Type
	//-------------------------------------------
	// 	1		Read Coil
	// 	2		Read Discrete Input
	// 	3		Read Holding Registers
	// 	4		Read Input Registers
	// 	5		Write Single Coil
	// 	6		Write Single Holding Register
	// 	15		Write Multiple Coils
	// 	16		Write Multiple Holding Registers
	Function int    `json:"function" validate:"1|2|3|4|"`               // Current version only support read
	Address  uint16 `json:"address" validate:"required,gte=0,lte=255"`  // Address
	Quantity uint16 `json:"quantity" validate:"required,gte=0,lte=255"` // Quantity
}
type ModBusConfig struct {
	Ip             string          `json:"ip" validate:"required"`
	Port           int             `json:"port" validate:"required,gte=1,lte=65535"`
	Timeout        int             `json:"timeout" validate:"required,gte=1,lte=60"`
	SlaverId       byte            `json:"slaverId" validate:"required,gte=1,lte=255"`
	Frequency      int64           `json:"frequency" validate:"required,gte=1,lte=10000"`
	RegisterParams []RegisterParam `json:"registerParams" validate:"required"`
}
type ModbusTcpMasterResource struct {
	typex.XStatus
	client  modbus.Client
	canWork bool
	cxt     context.Context
}

func NewModbusTcpMasterResource(id string, e typex.RuleX) typex.XResource {
	m := ModbusTcpMasterResource{}
	m.RuleEngine = e
	m.canWork = false
	m.cxt = context.Background()
	return &m
}

func (m *ModbusTcpMasterResource) Register(inEndId string) error {
	m.PointId = inEndId
	return nil
}

func (m *ModbusTcpMasterResource) Start() error {

	config := m.RuleEngine.GetInEnd(m.PointId).Config
	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	var mainConfig ModBusConfig
	if err := json.Unmarshal(configBytes, &mainConfig); err != nil {
		return err
	}
	if err := validator.New().Struct(mainConfig); err != nil {
		return err
	}

	handler := modbus.NewTCPClientHandler(
		fmt.Sprintf("%s:%v", mainConfig.Ip, mainConfig.Port),
	)
	handler.Timeout = time.Duration(mainConfig.Frequency) * time.Second
	handler.SlaveId = mainConfig.SlaverId
	if err := handler.Connect(); err != nil {
		return err
	}
	m.client = modbus.NewClient(handler)
	m.canWork = true
	for _, rCfg := range mainConfig.RegisterParams {
		log.Info("Start read register:", rCfg.Address)

		go func(ctx context.Context, rp RegisterParam) {
			// Modbus data is most often read and written as "registers" which are [16-bit] pieces of data. Most often,
			// the register is either a signed or unsigned 16-bit integer. If a 32-bit integer or floating point is required,
			// these values are actually read as a pair of registers.
			var results []byte
			var err error
			ticker := time.NewTicker(time.Duration(mainConfig.Frequency) * time.Second)
			for {
				<-ticker.C
				select {
				case <-ctx.Done():
					{
						return
					}
				default:
					{

						if rp.Function == 1 {
							results, err = m.client.ReadCoils(rp.Address, rp.Quantity)
						}
						if rp.Function == 2 {
							results, err = m.client.ReadDiscreteInputs(rp.Address, rp.Quantity)
						}
						if rp.Function == 3 {
							results, err = m.client.ReadHoldingRegisters(rp.Address, rp.Quantity)
						}
						if rp.Function == 4 {
							results, err = m.client.ReadInputRegisters(rp.Address, rp.Quantity)
						}
						//
						// error
						//
						if err != nil {
							m.canWork = false
							log.Error("NewModbusTcpMasterResource ReadData error: ", err)
						} else {
							if err0 := m.RuleEngine.PushQueue(typex.QueueData{
								In:   m.Details(),
								Out:  nil,
								E:    m.RuleEngine,
								Data: string(results),
							}); err0 != nil {
								log.Error("NewModbusTcpMasterResource PushQueue error: ", err0)
							}
						}

					}
				}

			}

		}(m.cxt, rCfg)
	}
	return nil

}

func (m *ModbusTcpMasterResource) Details() *typex.InEnd {
	return m.RuleEngine.GetInEnd(m.PointId)
}

func (m *ModbusTcpMasterResource) Test(inEndId string) bool {
	return true
}

func (m *ModbusTcpMasterResource) Enabled() bool {
	return m.Enable
}

func (m *ModbusTcpMasterResource) DataModels() *map[string]typex.XDataModel {
	return &map[string]typex.XDataModel{}
}

func (m *ModbusTcpMasterResource) Reload() {

}

func (m *ModbusTcpMasterResource) Pause() {

}

func (m *ModbusTcpMasterResource) Status() typex.ResourceState {
	if m.canWork {
		return typex.UP
	} else {
		return typex.DOWN
	}
}

func (m *ModbusTcpMasterResource) Stop() {

	m.cxt.Done()
}
