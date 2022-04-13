package source

import (
	"rulex/core"
	"rulex/driver"
	"rulex/typex"
	"rulex/utils"
	"time"

	"github.com/goburrow/serial"
	"github.com/ngaut/log"
)

type uartModuleSource struct {
	typex.XStatus
	uartDriver typex.XExternalDriver
}
type uartConfig struct {
	Address    string `json:"address" validate:"required" title:"串口地址" info:""`
	BaudRate   int    `json:"baudRate" validate:"required" title:"波特率" info:""`
	DataBits   int    `json:"dataBits" validate:"required" title:"数据位" info:""`
	StopBits   int    `json:"stopBits" validate:"required" title:"停止位" info:""`
	Parity     string `json:"parity" validate:"required" title:"分割大小" info:""`
	Timeout    *int64 `json:"timeout" validate:"required" title:"超时时间" info:""`
	BufferSize *int   `json:"bufferSize" validate:"required" title:"缓冲区大小" info:""`
}

func NewUartModuleSource(inEndId string, e typex.RuleX) typex.XSource {
	s := uartModuleSource{}
	s.PointId = inEndId
	s.RuleEngine = e
	//
	return &s
}
func (*uartModuleSource) Configs() *typex.XConfig {
	return core.GenInConfig(typex.UART_MODULE, "UART_MODULE", uartConfig{})
}

func (mm *uartModuleSource) DataModels() []typex.XDataModel {
	return []typex.XDataModel{}
}

func (u *uartModuleSource) Test(inEndId string) bool {
	return true
}
func (m *uartModuleSource) OnStreamApproached(data string) error {
	m.uartDriver.Write([]byte(data))
	return nil
}

func (u *uartModuleSource) Init(inEndId string, cfg map[string]interface{}) error {
	u.PointId = inEndId
	return nil
}
func (u *uartModuleSource) Start(cctx typex.CCTX) error {
	u.Ctx = cctx.Ctx
	u.CancelCTX = cctx.CancelCTX

	config := u.RuleEngine.GetInEnd(u.PointId).Config
	mainConfig := uartConfig{}
	if err := utils.BindSourceConfig(config, &mainConfig); err != nil {
		return err
	}
	uartDriver, err := driver.NewUartDriver(u.Ctx,
		serial.Config{
			Address:  mainConfig.Address,                               // 串口名
			BaudRate: mainConfig.BaudRate,                              // 115200
			DataBits: mainConfig.DataBits,                              // 8
			StopBits: mainConfig.StopBits,                              // 1
			Parity:   mainConfig.Parity,                                //'N'
			Timeout:  time.Duration(*mainConfig.Timeout) * time.Second, // 超时时间
		}, u.Details(), u.RuleEngine, *mainConfig.BufferSize, nil)
	if err != nil {
		return err
	}
	u.uartDriver = uartDriver
	return nil

}

func (u *uartModuleSource) Enabled() bool {
	return true
}

func (u *uartModuleSource) Reload() {
}

func (u *uartModuleSource) Pause() {

}
func (u *uartModuleSource) Details() *typex.InEnd {
	return u.RuleEngine.GetInEnd(u.PointId)
}

func (u *uartModuleSource) Status() typex.SourceState {
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

func (u *uartModuleSource) Stop() {
	if u.uartDriver != nil {
		u.uartDriver.Stop()
		u.uartDriver = nil
	}
	u.CancelCTX()

}
func (u *uartModuleSource) Driver() typex.XExternalDriver {
	return u.uartDriver
}

//
// 拓扑
//
func (*uartModuleSource) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}
