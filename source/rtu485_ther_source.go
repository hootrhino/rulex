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
	SlaverIds []byte     `json:"slaverIds" validate:"required" title:"TCP端口" info:""`
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
	clients   []modbus.Client
	clientIds []byte
	drivers   []typex.XExternalDriver
}

func NewRtu485THerSource(e typex.RuleX) typex.XSource {
	m := rtu485THerSource{}
	m.RuleEngine = e
	m.drivers = make([]typex.XExternalDriver, 0)
	m.clients = make([]modbus.Client, 0)
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

	// handler.Logger = golog.New(os.Stdout, "485THerSource: ", log.LstdFlags)
	// 串口配置固定写法
	handler := modbus.NewRTUClientHandler(rtuConfig.Uart)
	handler.BaudRate = 4800
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.Timeout = time.Duration(*mainConfig.Timeout) * time.Second
	for i, sId := range mainConfig.SlaverIds {
		handler.SlaveId = sId
		if err := handler.Connect(); err != nil {
			return err
		}
		m.clients = append(m.clients, modbus.NewClient(handler))
		driver := driver.NewRtu485_THer_Driver(m.Details(), m.RuleEngine, m.clients[i])
		m.drivers = append(m.drivers, driver)
		m.clientIds = append(m.clientIds, sId)
	}
	//---------------------------------------------------------------------------------
	// Start
	//---------------------------------------------------------------------------------

	for i, driver := range m.drivers {
		ticker := time.NewTicker(time.Duration(3) * time.Second)
		// log.Debug("start driver:", i)
		go func(ctx context.Context, idx byte, rtuDriver typex.XExternalDriver) {
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
							log.Error("uart read error, SlaverId: ", m.clientIds[idx], " error:", err, ", may should check uart config if baud rate is correct?")
							_sourceState = typex.DOWN
							return
						}
						b, _ := json.Marshal(_data{
							SlaverId:    m.clientIds[idx],
							Humidity:    float32(binary.BigEndian.Uint16(buffer[0:2])) * 0.1,
							Temperature: float32(binary.BigEndian.Uint16(buffer[2:n])) * 0.1,
						})
						m.RuleEngine.Work(m.RuleEngine.GetInEnd(m.PointId), string(b))
					}

				}
			}
		}(m.Ctx, byte(i), driver)
	}

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
	for _, d := range m.drivers {
		d.Stop()
	}
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
	if m.drivers[0] != nil {
		return m.drivers[0]
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
