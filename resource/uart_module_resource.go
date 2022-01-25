package resource

import (
	"rulex/core"
	"rulex/driver"
	"rulex/typex"
	"rulex/utils"
	"time"

	"github.com/goburrow/serial"
	"github.com/ngaut/log"
)

type uartModuleResource struct {
	typex.XStatus
	uartDriver typex.XExternalDriver
}
type uartConfig struct {
	Address  string `json:"address" validate:"required" title:"串口地址" info:""`
	BaudRate int    `json:"baudRate" validate:"required" title:"波特率" info:""`
	DataBits int    `json:"dataBits" validate:"required" title:"数据位" info:""`
	StopBits int    `json:"stopBits" validate:"required" title:"停止位" info:""`
	Parity   string `json:"parity" validate:"required" title:"分割大小" info:""`
	Timeout  *int64 `json:"timeout" validate:"required" title:"超时时间" info:""`
}

func NewUartModuleResource(inEndId string, e typex.RuleX) typex.XResource {
	s := uartModuleResource{}
	s.PointId = inEndId
	s.RuleEngine = e
	//
	return &s
}
func (*uartModuleResource) Configs() *typex.XConfig {
	return core.GenInConfig(typex.UART_MODULE, "UART_MODULE", uartConfig{})
}

func (mm *uartModuleResource) DataModels() []typex.XDataModel {
	return []typex.XDataModel{}
}

func (u *uartModuleResource) Test(inEndId string) bool {
	return true
}
func (m *uartModuleResource) OnStreamApproached(data string) error {
	m.uartDriver.Write([]byte(data))
	return nil
}
func (u *uartModuleResource) Register(inEndId string) error {
	u.PointId = inEndId
	return nil
}

func (u *uartModuleResource) Start() error {
	config := u.RuleEngine.GetInEnd(u.PointId).Config
	mainConfig := uartConfig{}
	if err := utils.BindResourceConfig(config, &mainConfig); err != nil {
		return err
	}
	driver, err := driver.NewUartDriver(serial.Config{
		Address:  mainConfig.Address,  // 串口名
		BaudRate: mainConfig.BaudRate, // 115200
		DataBits: mainConfig.DataBits, // 8
		StopBits: mainConfig.StopBits, // 1
		Parity:   mainConfig.Parity,   //'N'
		Timeout:  time.Duration(*mainConfig.Timeout) * time.Second,
	}, u.Details(), u.RuleEngine, nil)
	if err != nil {
		return err
	}
	u.uartDriver = driver
	return nil

}

func (u *uartModuleResource) Enabled() bool {
	return true
}

func (u *uartModuleResource) Reload() {
}

func (u *uartModuleResource) Pause() {

}
func (u *uartModuleResource) Details() *typex.InEnd {
	return u.RuleEngine.GetInEnd(u.PointId)
}

func (u *uartModuleResource) Status() typex.ResourceState {
	if u.uartDriver != nil {
		if err := u.uartDriver.Test(); err != nil {
			log.Error(err)
			return typex.DOWN
		} else {
			return typex.UP
		}
	}
	return typex.DOWN
}

func (u *uartModuleResource) Stop() {
	if u.uartDriver != nil {
		u.uartDriver.Stop()
		u.uartDriver = nil
	}

}
func (u *uartModuleResource) Driver() typex.XExternalDriver {
	return u.uartDriver
}

//
// 拓扑
//
func (*uartModuleResource) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}
