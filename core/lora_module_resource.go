package core

import (
	"github.com/ngaut/log"
	"github.com/tarm/serial"
	"rulex/drivers"
)

type LoraModuleResource struct {
	XStatus
	loraDriver *drivers.ATK_LORA_01Driver
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
	if err, _ := s.loraDriver.Test(); err != nil {
		log.Error(err)
		return false
	} else {
		return true
	}
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
		Name:        name.(string),
		Baud:        int(baud.(float64)),
		Parity:      'N',
		ReadTimeout: 0,
		Size:        0,
		StopBits:    1,
	})
	if err != nil {
		log.Error("LoraModuleResource start failed:", err)
		return err
	} else {
		s.loraDriver = drivers.NewATK_LORA_01Driver(serialPort)
		s.loraDriver.Init()
		s.loraDriver.Work()
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
func (s *LoraModuleResource) Details() *inEnd {
	return s.RuleEngine.GetInEnd(s.PointId)
}

func (s *LoraModuleResource) Status() ResourceState {
	if s.loraDriver != nil {
		if err, _ := s.loraDriver.Test(); err != nil {
			log.Error(err)
			return DOWN
		} else {
			return UP
		}
	}
	return DOWN
}

func (s *LoraModuleResource) Stop() {
	if s.loraDriver != nil {
		s.loraDriver.Stop()
	}
}
