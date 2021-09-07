package core

import (
	"rulex/drivers"

	"github.com/ngaut/log"
	"github.com/tarm/serial"
)

type LoraModuleResource struct {
	XStatus
	serialPort *serial.Port
	loraDriver *gpio.ATK_LORA_01Driver
}

func NewLoraModuleResource(inEndId string, e *RuleEngine) *LoraModuleResource {
	s := LoraModuleResource{}
	s.PointId = inEndId
	s.RuleEngine = e
	//
	return &s
}

func (mm *LoraModuleResource) DataModels() *map[string]XDataModel {
	return &map[string]XDataModel{
		"NodeData": {
			Type:      T_JSON,
			Name:      "NodeSendMsg",
			MinLength: 2,
			MaxLength: 1024,
		},
	}
}

func (s *LoraModuleResource) Test(inEndId string) bool {
	return true
}

func (s *LoraModuleResource) Register(inEndId string) error {
	return nil
}

func (s *LoraModuleResource) Start() error {
	config := s.RuleEngine.GetInEnd(s.PointId).Config
	name := (*config)["name"]
	baud := (*config)["baud"]
	//readTimeout := (*config)["readTimeout"]
	//size := (*config)["size"]
	//parity := (*config)["parity"]
	//stopbits := (*config)["stopbits"]

	serialPort, err := serial.OpenPort(&serial.Config{
		Name:   name.(string),
		Baud:   int(baud.(float64)),
		Parity: 'N',
		//ReadTimeout: time.Duration(readTimeout.(int)),
		//Size:        size.(byte),
		//StopBits: serial.StopBits(stopbits.(float64)),
	})
	if err != nil {
		log.Error("LoraModuleResource start failed:", err)
		return err
	} else {
		s.serialPort = serialPort
		s.loraDriver = gpio.NewATK_LORA_01Driver(serialPort)
		log.Info("LoraModuleResource start success.")
		return nil
	}
}

func (s *LoraModuleResource) Enabled() bool {
	return true
}

func (s *LoraModuleResource) Reload() {
}

func (s *LoraModuleResource) Pause() {

}

func (s *LoraModuleResource) Status() State {
	return UP
}

func (s *LoraModuleResource) Stop() {
	s.serialPort.Close()
}
