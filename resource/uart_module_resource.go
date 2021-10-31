package resource

import (
	"rulex/driver"
	"rulex/typex"
	"rulex/utils"
	"time"

	"github.com/goburrow/serial"
	"github.com/ngaut/log"
)

type UartModuleResource struct {
	typex.XStatus
	loraDriver typex.XExternalDriver
}
type UartConfig struct {
	Address  string `json:"address" validate:"required"`
	BaudRate int    `json:"baudRate" validate:"required"`
	DataBits int    `json:"dataBits" validate:"required"`
	StopBits int    `json:"stopBits" validate:"required"`
	Parity   string `json:"parity" validate:"required"`
	Timeout  int64  `json:"timeout" validate:"required"`
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
	m.loraDriver.Write([]byte(data))
	return nil
}
func (s *UartModuleResource) Register(inEndId string) error {
	s.PointId = inEndId
	return nil
}

func (s *UartModuleResource) Start() error {
	config := s.RuleEngine.GetInEnd(s.PointId).Config
	mainConfig := UartConfig{}
	if err := utils.BindResourceConfig(config, &mainConfig); err != nil {
		return err
	}
	serialPort, err := serial.Open(&serial.Config{
		Address:  mainConfig.Address,
		BaudRate: mainConfig.BaudRate, //115200
		DataBits: mainConfig.DataBits, //8
		StopBits: mainConfig.StopBits, //1
		Parity:   "N",                 //'N'
		Timeout:  time.Duration(mainConfig.Timeout) * time.Second,
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
