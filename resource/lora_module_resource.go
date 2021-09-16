package resource

import (
	"rulex/drivers"
	"rulex/typex"

	"github.com/ngaut/log"
	"github.com/tarm/serial"
)

type LoraModuleResource struct {
	typex.XStatus
	loraDriver *drivers.ATK_LORA_01Driver
}

func NewLoraModuleResource(inEndId string, e typex.RuleX) *LoraModuleResource {
	s := LoraModuleResource{}
	s.PointId = inEndId
	s.RuleEngine = e
	//
	return &s
}

func (mm *LoraModuleResource) DataModels() *map[string]typex.XDataModel {
	return &map[string]typex.XDataModel{
		"NodeData": {
			Type:      typex.T_JSON,
			Name:      "NodeSendMsg",
			MinLength: 2,
			MaxLength: 1024,
		},
	}
}

func (s *LoraModuleResource) Test(inEndId string) bool {
	if _, err := s.loraDriver.Test(); err != nil {
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
func (s *LoraModuleResource) Details() *typex.InEnd {
	return s.RuleEngine.GetInEnd(s.PointId)
}

func (s *LoraModuleResource) Status() typex.ResourceState {
	if s.loraDriver != nil {
		if _, err := s.loraDriver.Test(); err != nil {
			log.Error(err)
			return typex.DOWN
		} else {
			return typex.UP
		}
	}
	return typex.DOWN
}

func (s *LoraModuleResource) Stop() {
	if s.loraDriver != nil {
		s.loraDriver.Stop()
	}
}
