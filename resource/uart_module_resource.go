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
	Address   string `json:"address" validate:"required"`
	BaudRate  int    `json:"baudRate" validate:"required"`
	DataBits  int    `json:"dataBits" validate:"required"`
	StopBits  int    `json:"stopBits" validate:"required"`
	Parity    string `json:"parity" validate:"required"`
	Separator string `json:"separator" validate:"required"`
	Timeout   *int64 `json:"timeout" validate:"required"`
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

func (u *UartModuleResource) Test(inEndId string) bool {
	if err := u.loraDriver.Test(); err != nil {
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
func (u *UartModuleResource) Register(inEndId string) error {
	u.PointId = inEndId
	return nil
}

func (u *UartModuleResource) Start() error {
	config := u.RuleEngine.GetInEnd(u.PointId).Config
	mainConfig := UartConfig{}
	if err := utils.BindResourceConfig(config, &mainConfig); err != nil {
		return err
	}
	serialPort, err := serial.Open(&serial.Config{
		Address:  mainConfig.Address,
		BaudRate: mainConfig.BaudRate, //115200
		DataBits: mainConfig.DataBits, //8
		StopBits: mainConfig.StopBits, //1
		Parity:   mainConfig.Parity,   //'N'
		Timeout:  time.Duration(*mainConfig.Timeout) * time.Second,
	})
	if err != nil {
		log.Error("UartModuleResource start failed:", err)
		return err
	} else {
		log.Info("Uart port open successfully:", mainConfig.Address)
		u.loraDriver = driver.NewUartDriver(
			serialPort,
			u.Details(),
			u.RuleEngine,
			mainConfig.Separator[0])
		return nil
	}
}

func (u *UartModuleResource) Enabled() bool {
	return true
}

func (u *UartModuleResource) Reload() {
}

func (u *UartModuleResource) Pause() {

}
func (u *UartModuleResource) Details() *typex.InEnd {
	return u.RuleEngine.GetInEnd(u.PointId)
}

func (u *UartModuleResource) Status() typex.ResourceState {
	if u.loraDriver != nil {
		if err := u.loraDriver.Test(); err != nil {
			log.Error(err)
			return typex.DOWN
		} else {
			return typex.UP
		}
	}
	return typex.DOWN
}

func (u *UartModuleResource) Stop() {
	if u.loraDriver != nil {
		u.loraDriver.Stop()
		u.loraDriver = nil
	}

}
func (u *UartModuleResource) Driver() typex.XExternalDriver {
	return u.loraDriver
}
