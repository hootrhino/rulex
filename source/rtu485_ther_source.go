package source

import (
	"context"
	"encoding/binary"
	"encoding/json"

	"rulex/core"
	"rulex/driver"
	"rulex/typex"
	"rulex/utils"
	"time"

	"github.com/ngaut/log"

	"github.com/goburrow/modbus"
	"github.com/mitchellh/mapstructure"
)

type _modBusConfig struct {
	Timeout   *int       `json:"timeout" validate:"required" title:"连接超时" info:""`
	SlaverId  byte       `json:"slaverId" validate:"required" title:"TCP端口" info:""`
	Frequency *int64     `json:"frequency" validate:"required" title:"采集频率" info:""`
	Config    _rtuConfig `json:"config" validate:"required" title:"工作模式" info:""`
}

type _rtuConfig struct {
	Uart     string       `json:"uart" validate:"required" title:"串口路径" info:"本地系统的串口路径"`
	BaudRate int          `json:"baudRate" validate:"required" title:"波特率" info:"串口通信波特率"`
	DataBits int          `json:"dataBits" validate:"required" title:"数据位" info:"串口通信数据位"`
	Parity   typex.Parity `json:"parity" validate:"required" title:"校验位" info:"串口通信校验位"`
	StopBits int          `json:"stopBits" validate:"required" title:"停止位" info:"串口通信停止位"`
}
type _data struct {
	SlaverId    byte    `json:"slaverId"`
	Humidity    float32 `json:"humidity"`
	Temperature float32 `json:"temperature"`
}

//
//
//---------------------------------------------------------------------------

type rtu485THerSource struct {
	typex.XStatus
	client    modbus.Client
	rtuDriver typex.XExternalDriver
}

func NewRtu485THerSource(e typex.RuleX) typex.XSource {
	m := rtu485THerSource{}
	m.RuleEngine = e
	return &m
}
func (*rtu485THerSource) Configs() *typex.XConfig {
	return core.GenInConfig(typex.RTU485_THER, "RTU485_THER", _modBusConfig{})
}

func (m *rtu485THerSource) Init(inEndId string, cfg map[string]interface{}) error {
	m.PointId = inEndId
	return nil
}
func (m *rtu485THerSource) Start(cctx typex.CCTX) error {
	m.Ctx = cctx.Ctx
	m.CancelCTX = cctx.CancelCTX

	config := m.RuleEngine.GetInEnd(m.PointId).Config
	var mainConfig _modBusConfig
	if err := utils.BindSourceConfig(config, &mainConfig); err != nil {
		return err
	}
	var rtuConfig rtuConfig
	if errs := mapstructure.Decode(mainConfig.Config, &rtuConfig); errs != nil {
		log.Error(errs)
		return errs
	}
	handler := modbus.NewRTUClientHandler(rtuConfig.Uart)
	// handler.Logger = golog.New(os.Stdout, "485THerSource: ", log.LstdFlags)
	handler.BaudRate = rtuConfig.BaudRate
	handler.DataBits = rtuConfig.DataBits
	handler.Parity = rtuConfig.Parity
	handler.StopBits = rtuConfig.StopBits
	handler.SlaveId = mainConfig.SlaverId
	handler.Timeout = time.Duration(*mainConfig.Timeout) * time.Second
	if err := handler.Connect(); err != nil {
		return err
	}
	m.client = modbus.NewClient(handler)
	m.rtuDriver = driver.NewRtu485_THer_Driver(m.Details(), m.RuleEngine, m.client)
	//---------------------------------------------------------------------------------
	// Start
	//---------------------------------------------------------------------------------

	ticker := time.NewTicker(time.Duration(*mainConfig.Frequency) * time.Second)
	go func(ctx context.Context, rtuDriver typex.XExternalDriver) {
		defer ticker.Stop()
		buffer := make([]byte, 4) //4字节数据
		for {
			<-ticker.C
			select {
			case <-ctx.Done():
				{
					_sourceState = typex.DOWN
					return
				}
			default:
				{
					n, err := rtuDriver.Read(buffer)
					if err != nil {
						log.Error("uart read error: ", err, ", may should check uart config if baud rate is correct?")
						_sourceState = typex.DOWN
						return
					}
					b, _ := json.Marshal(_data{
						SlaverId:    mainConfig.SlaverId,
						Humidity:    float32(binary.BigEndian.Uint16(buffer[0:2])) * 0.1,
						Temperature: float32(binary.BigEndian.Uint16(buffer[2:n])) * 0.1,
					})
					m.RuleEngine.Work(m.RuleEngine.GetInEnd(m.PointId), string(b))
				}

			}
		}
	}(m.Ctx, m.rtuDriver)
	return nil

}

func (m *rtu485THerSource) Details() *typex.InEnd {
	return m.RuleEngine.GetInEnd(m.PointId)
}

func (m *rtu485THerSource) Test(inEndId string) bool {
	return true
}

func (m *rtu485THerSource) Enabled() bool {
	return m.Enable
}

func (m *rtu485THerSource) DataModels() []typex.XDataModel {
	return m.XDataModels
}

func (m *rtu485THerSource) Reload() {

}

func (m *rtu485THerSource) Pause() {

}

func (m *rtu485THerSource) Status() typex.SourceState {
	return _sourceState

}

func (m *rtu485THerSource) Stop() {
	m.CancelCTX()
}

/*
*
* 写入值
*
 */
func (m *rtu485THerSource) OnStreamApproached(data string) error {
	return nil
}

/*
*
* 只有RTU模式才带驱动
*
 */
func (m *rtu485THerSource) Driver() typex.XExternalDriver {
	if m.client != nil {
		return m.rtuDriver
	} else {
		return nil
	}
}

//
// 拓扑
//
func (*rtu485THerSource) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}
