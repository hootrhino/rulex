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
	Function int    `json:"function" validate:"1|2|3|4|"`      // Current version only support read
	Address  uint16 `json:"address" validate:"gte=0,lte=255"`  // Address
	Quantity uint16 `json:"quantity" validate:"gte=0,lte=255"` // Quantity
}
type ModBusConfig struct {
	Ip             string          `json:"ip" validate:"required"`
	Port           int             `json:"port" validate:"gte=1024,lte=65535"`
	Timeout        int             `json:"timeout" validate:"gte=1,lte=60"`
	SlaverId       byte            `json:"slaverId" validate:"gte=1,lte=255"`
	RegisterParams []RegisterParam `json:"registerParams" validate:"required"`
}
type ModbusTcpMasterResource struct {
	typex.XStatus
	client modbus.Client
}

func NewModbusTcpMasterResource(e typex.RuleX) typex.XResource {
	m := ModbusTcpMasterResource{}
	m.RuleEngine = e
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
	handler.Timeout = time.Duration(mainConfig.Timeout) * time.Second
	handler.SlaveId = mainConfig.SlaverId
	handler.Logger = log.Logger()
	if err := handler.Connect(); err != nil {
		return err
	}
	m.client = modbus.NewClient(handler)

	for _, rcfg := range mainConfig.RegisterParams {
		log.Info("Start read register:", rcfg.Address)

		go func(ctx context.Context, p RegisterParam) {
			var results []byte
			var err error
			for {
				select {
				case <-ctx.Done():
					{
						return
					}
				default:
					{

						if p.Function == 1 {
							results, err = m.client.ReadCoils(p.Address, p.Quantity)
						}
						if p.Function == 2 {
							results, err = m.client.ReadDiscreteInputs(p.Address, p.Quantity)
						}
						if p.Function == 3 {
							results, err = m.client.ReadHoldingRegisters(p.Address, p.Quantity)
						}
						if p.Function == 4 {
							results, err = m.client.ReadInputRegisters(p.Address, p.Quantity)
						}
						//
						// error
						//
						if err != nil {
							log.Error(err)
						} else {
							log.Info("m.client.ReadCoils", results)
						}
					}
				}

			}

		}(context.Background(), rcfg)
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
	return typex.UP
}

func (m *ModbusTcpMasterResource) Stop() {

}
