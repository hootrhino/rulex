package resource

import (
	"rulex/driver"
	"rulex/typex"

	"github.com/ngaut/log"
	"github.com/tarm/serial"
)

type UartModuleResource struct {
	typex.XStatus
	loraDriver typex.XExternalDriver
}

func NewUartModuleResource(inEndId string, e typex.RuleX) typex.XResource {
	s := UartModuleResource{}
	s.PointId = inEndId
	s.RuleEngine = e
	//
	return &s
}

func (mm *UartModuleResource) DataModels() []typex.XDataModel {
	return []typex.XDataModel{}
}

func (s *UartModuleResource) Test(inEndId string) bool {
	if err := s.loraDriver.Test(); err != nil {
		log.Error(err)
		return false
	} else {
		return true
	}
}
func (m *UartModuleResource) OnStreamApproached(data string) error {
	return nil
}
func (s *UartModuleResource) Register(inEndId string) error {
	s.PointId = inEndId
	return nil
}

func (s *UartModuleResource) Start() error {
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
		log.Error("UartModuleResource start failed:", err)
		return err
	} else {
		s.loraDriver = driver.NewUartDriver(serialPort, s.Details(), s.RuleEngine)
		err0 := s.loraDriver.Init()
		if err != nil {
			return err0
		}
		err1 := s.loraDriver.Work()
		if err != nil {
			return err1
		}
		log.Info("UartModuleResource start success.")
		return nil
	}
}

func (s *UartModuleResource) Enabled() bool {
	return true
}

func (s *UartModuleResource) Reload() {
}

func (s *UartModuleResource) Pause() {

}
func (s *UartModuleResource) Details() *typex.InEnd {
	return s.RuleEngine.GetInEnd(s.PointId)
}

func (s *UartModuleResource) Status() typex.ResourceState {
	if s.loraDriver != nil {
		if err := s.loraDriver.Test(); err != nil {
			log.Error(err)
			return typex.DOWN
		} else {
			return typex.UP
		}
	}
	return typex.DOWN
}

func (s *UartModuleResource) Stop() {
	if s.loraDriver != nil {
		s.loraDriver.Stop()
	}
}
