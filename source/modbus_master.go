package source

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"time"

	"github.com/i4de/rulex/core"

	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"

	"github.com/goburrow/modbus"
	"github.com/mitchellh/mapstructure"
)

// 资源状态
var _sourceState typex.SourceState = typex.SOURCE_UP

type modBusConfig struct {
	Mode           string          `json:"mode" title:"工作模式" info:"RTU/TCP"`
	Timeout        int             `json:"timeout" validate:"required" title:"连接超时" info:""`
	SlaverId       byte            `json:"slaverId" validate:"required" title:"TCP端口" info:""`
	Frequency      int64           `json:"frequency" validate:"required" title:"采集频率" info:""`
	Config         interface{}     `json:"config" validate:"required" title:"工作模式" info:""`
	RegisterParams []registerParam `json:"registerParams" validate:"required" title:"寄存器配置" info:""`
}

const (
	//-------------------------------------------
	// 	Code |  Register Type
	//-------|------------------------------------
	// 	1	 |	Read Coil
	// 	2	 |	Read Discrete Input
	// 	3	 |	Read Holding Registers
	// 	4	 |	Read Input Registers
	// 	5	 |	Write Single Coil
	// 	6	 |	Write Single Holding Register
	// 	15	 |	Write Multiple Coils
	// 	16	 |	Write Multiple Holding Registers
	//-------------------------------------------

	READ_COIL                        = 1
	READ_DISCRETE_INPUT              = 2
	READ_HOLDING_REGISTERS           = 3
	READ_INPUT_REGISTERS             = 4
	WRITE_SINGLE_COIL                = 5
	WRITE_SINGLE_HOLDING_REGISTER    = 6
	WRITE_MULTIPLE_COILS             = 15
	WRITE_MULTIPLE_HOLDING_REGISTERS = 16
)

/*
*
* coilParams 1
*
 */
type coilsW struct {
	Address  uint16 `json:"address" validate:"required" title:"寄存器地址" info:""`
	Quantity uint16 `json:"quantity" validate:"required" title:"写入数量" info:""`
	Values   []byte `json:"values" validate:"required" title:"写入的值" info:""`
}

/*
*
* 2
*
 */
type coilW struct {
	Address  uint16 `json:"address" validate:"required" title:"寄存器地址" info:""`
	Quantity uint16 `json:"quantity" validate:"required" title:"写入数量" info:""`
	Value    uint16 `json:"value" validate:"required" title:"写入的值" info:""`
}

/*
*
* registerParams 3
*
 */
type registerW struct {
	Address  uint16 `json:"address" validate:"required" title:"寄存器地址" info:""`
	Quantity uint16 `json:"quantity" validate:"required" title:"写入数量" info:""`
	Value    uint16 `json:"value" validate:"required" title:"写入的值" info:""`
}

/*
*
* 4
*
 */
type registersW struct {
	Address  uint16 `json:"address" validate:"required" title:"寄存器地址" info:""`
	Quantity uint16 `json:"quantity" validate:"required" title:"写入数量" info:""`
	Values   []byte `json:"values" validate:"required" title:"写入的值" info:""`
}

/*
*
* 配置进来准备采集的寄存器参数
*
 */

type registerParam struct {
	Tag      string `json:"tag" validate:"required"`      // Function
	Function int    `json:"function" validate:"required"` // Function
	Address  uint16 `json:"address" validate:"required"`  // Address
	Quantity uint16 `json:"quantity" validate:"required"` // Quantity
}

/*
*
* 采集到的数据
*
 */
type registerData struct {
	Tag      string `json:"tag" validate:"required"`      // Function
	Function int    `json:"function" validate:"required"` // Function
	SlaverId byte   `json:"slaverId" validate:"required"`
	Address  uint16 `json:"address" validate:"required"`  // Address
	Quantity uint16 `json:"quantity" validate:"required"` // Quantity
	Value    string `json:"value" validate:"required"`    // Quantity
}

//
// Uart "/dev/ttyUSB0"
// BaudRate = 115200
// DataBits = 8
// Parity = "N"
// StopBits = 1
// SlaveId = 1
// Timeout = 5 * time.Second
//
type rtuConfig struct {
	Uart     string `json:"uart" validate:"required" title:"串口路径" info:"本地系统的串口路径"`
	BaudRate int    `json:"baudRate" validate:"required" title:"波特率" info:"串口通信波特率"`
	DataBits int    `json:"dataBits" validate:"required" title:"数据位" info:"串口通信数据位"`
	Parity   string `json:"parity" validate:"required" title:"校验位" info:"串口通信分割位"`
	StopBits int    `json:"stopBits" validate:"required" title:"停止位" info:"串口通信停止位"`
}

//
//
//
type tcpConfig struct {
	Ip   string `json:"ip" validate:"required" title:"IP地址" info:""`
	Port int    `json:"port" validate:"required" title:"端口" info:""`
}

//
//
//---------------------------------------------------------------------------

type modbusMasterSource struct {
	typex.XStatus
	client modbus.Client
}

func NewModbusMasterSource(id string, e typex.RuleX) typex.XSource {
	m := modbusMasterSource{}
	m.RuleEngine = e
	return &m
}
func (*modbusMasterSource) Configs() *typex.XConfig {
	return core.GenInConfig(typex.MODBUS_MASTER, "MODBUS_MASTER", modBusConfig{})
}

func (m *modbusMasterSource) Init(inEndId string, cfg map[string]interface{}) error {
	m.PointId = inEndId

	return nil
}
func (m *modbusMasterSource) Start(cctx typex.CCTX) error {
	m.Ctx = cctx.Ctx
	m.CancelCTX = cctx.CancelCTX

	config := m.RuleEngine.GetInEnd(m.PointId).Config
	var mainConfig modBusConfig
	if err := utils.BindSourceConfig(config, &mainConfig); err != nil {
		return err
	}

	if mainConfig.Mode == "TCP" {
		var tcpConfig tcpConfig
		if errs := mapstructure.Decode(mainConfig.Config, &tcpConfig); errs != nil {
			glogger.GLogger.Error(errs)
			return errs
		}
		handler := modbus.NewTCPClientHandler(
			fmt.Sprintf("%s:%v", tcpConfig.Ip, tcpConfig.Port),
		)
		handler.Timeout = time.Duration(mainConfig.Frequency) * time.Second
		handler.SlaveId = mainConfig.SlaverId
		if err := handler.Connect(); err != nil {
			return err
		}
		m.client = modbus.NewClient(handler)
	} else if mainConfig.Mode == "RTU" {
		var rtuConfig rtuConfig
		if errs := mapstructure.Decode(mainConfig.Config, &rtuConfig); errs != nil {
			glogger.GLogger.Error(errs)
			return errs
		}
		handler := modbus.NewRTUClientHandler(rtuConfig.Uart)
		handler.BaudRate = rtuConfig.BaudRate
		//
		// rtuConfig
		//
		handler.DataBits = rtuConfig.DataBits
		handler.Parity = rtuConfig.Parity
		handler.StopBits = rtuConfig.StopBits
		//---------
		handler.SlaveId = mainConfig.SlaverId
		handler.Timeout = time.Duration(mainConfig.Frequency) * time.Second
		//---------
		if err := handler.Connect(); err != nil {
			return err
		}
		m.client = modbus.NewClient(handler)
	} else {
		return errors.New("no supported mode:" + mainConfig.Mode)
	}
	//---------------------------------------------------------------------------------
	// Start
	//---------------------------------------------------------------------------------

	ticker := time.NewTicker(time.Duration(mainConfig.Frequency) * time.Second)
	//
	// 前端传过来个寄存器和地址的配置列表，然后给每个寄存器分配一个协程去读
	//

	for _, rCfg := range mainConfig.RegisterParams {
		glogger.GLogger.Info("Start read register:", rCfg.Address)
		// 每个寄存器配一个协程读数据
		go func(ctx context.Context, rp registerParam) {
			defer ticker.Stop()
			// Modbus data is most often read and written as "registers" which are [16-bit] pieces of data. Most often,
			// the register is either a signed or unsigned 16-bit integer. If a 32-bit integer or floating point is required,
			// these values are actually read as a pair of registers.
			var results []byte
			var err error
			for {
				<-ticker.C
				select {
				case <-ctx.Done():
					{
						_sourceState = typex.SOURCE_DOWN
						return
					}
				default:
					{
						if rp.Function == READ_COIL {
							results, err = m.client.ReadCoils(rp.Address, rp.Quantity)
						}
						if rp.Function == READ_DISCRETE_INPUT {
							results, err = m.client.ReadDiscreteInputs(rp.Address, rp.Quantity)
						}
						if rp.Function == READ_HOLDING_REGISTERS {
							results, err = m.client.ReadHoldingRegisters(rp.Address, rp.Quantity)
						}
						if rp.Function == READ_INPUT_REGISTERS {
							results, err = m.client.ReadInputRegisters(rp.Address, rp.Quantity)
						}
						//
						// error
						//
						if err != nil {
							glogger.GLogger.Error("NewModbusMasterSource ReadData error: ", err)
							_sourceState = typex.SOURCE_DOWN

						} else {
							data := registerData{
								Tag:      rp.Tag,
								Function: rp.Function,
								Address:  rp.Address,
								SlaverId: mainConfig.SlaverId,
								Quantity: rp.Quantity,
								// 默认将Modbus数据编码成十六进制格式
								Value: hex.EncodeToString(results),
							}
							bytes, _ := json.Marshal(data)
							m.RuleEngine.WorkInEnd(m.RuleEngine.GetInEnd(m.PointId), string(bytes))
						}
					}
				}
			}
		}(m.Ctx, rCfg)
	}
	return nil

}

func (m *modbusMasterSource) Details() *typex.InEnd {
	return m.RuleEngine.GetInEnd(m.PointId)
}

func (m *modbusMasterSource) Test(inEndId string) bool {
	return true
}

func (m *modbusMasterSource) Enabled() bool {
	return m.Enable
}

func (m *modbusMasterSource) DataModels() []typex.XDataModel {
	return m.XDataModels
}

func (m *modbusMasterSource) Reload() {

}

func (m *modbusMasterSource) Pause() {

}

func (m *modbusMasterSource) Status() typex.SourceState {
	return _sourceState

}

func (m *modbusMasterSource) Stop() {
	m.CancelCTX()
}

/*
*
* 对于 modbus 资源来说, 任何直接写入的数据都被认为是给寄存器写值
*
 */
type dataM struct {
	Type   int                    `json:"type" validate:"required"`
	Params map[string]interface{} `json:"params" validate:"required"`
}

/*
*
* 只有RTU模式才带驱动
*
 */
func (m *modbusMasterSource) Driver() typex.XExternalDriver {
	return nil
}

//
// 拓扑
//
func (*modbusMasterSource) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}

//
// 来自外面的数据
//
func (*modbusMasterSource) DownStream([]byte) {}

//
// 上行数据
//
func (*modbusMasterSource) UpStream() {}
