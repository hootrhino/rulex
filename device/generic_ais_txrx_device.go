package device

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
)

// --------------------------------------------------------------------------------------------------
// 把AIS包里面的几个结构体拿出来了，主要是适配JSON格式, 下面这些结构体和AIS包里面的完全一样
// --------------------------------------------------------------------------------------------------
/*
*
* 公共结构
*
 */
type _AISDeviceBaseSentence struct {
	Talker   string             `json:"talker,omitempty"`   // The talker id (e.g GP)
	Type     string             `json:"type,omitempty"`     // The data type (e.g GSA)
	Fields   []string           `json:"fields,omitempty"`   // Array of fields
	Checksum string             `json:"checksum,omitempty"` // The Checksum
	Raw      string             `json:"raw,omitempty"`      // The raw NMEA sentence received
	TagBlock _AISDeviceTagBlock `json:"tagBlock,omitempty"` // NMEA tagblock
}

/*
*
* Tag
*
 */
type _AISDeviceTagBlock struct {
	Time         int64  `json:"time"`         // TypeUnixTime unix timestamp (unit is likely to be s, but might be ms, YMMV), parameter: -c
	RelativeTime int64  `json:"relativeTime"` // TypeRelativeTime relative time, parameter: -r
	Destination  string `json:"destination"`  // TypeDestinationID destination identification 15 char max, parameter: -d
	Grouping     string `json:"grouping"`     // TypeGrouping sentence grouping, parameter: -g
	LineCount    int64  `json:"lineCount"`    // TypeLineCount line count, parameter: -n
	Source       string `json:"source"`       // TypeSourceID source identification 15 char max, parameter: -s
	Text         string `json:"text"`         // TypeTextString valid character string, parameter -t
}

/*
*
* AIS包
*
 */
type _AISDevicePacket struct {
	_AISDeviceBaseSentence
	NumFragments   int64  `json:"numFragments"`
	FragmentNumber int64  `json:"fragmentNumber"`
	MessageID      int64  `json:"messageID"`
	Channel        string `json:"channel"`
	Payload        []byte `json:"payload"`
}

// --------------------------------------------------------------------------------------------------
type _AISDeviceConfig struct {
	Mode string `json:"mode"` // TCP UDP UART
	Host string `json:"host" validate:"required"`
	Port int    `json:"port" validate:"required"`
}
type AISDevice struct {
	typex.XStatus
	status      typex.DeviceState
	mainConfig  _AISDeviceConfig
	RuleEngine  typex.RuleX
	tcpListener net.Listener // TCP 接收端
}

/*
*
* AIS 数据解析服务器
*
 */
func NewAISDevice(e typex.RuleX) typex.XDevice {
	aisd := new(AISDevice)
	aisd.RuleEngine = e
	aisd.mainConfig = _AISDeviceConfig{
		Mode: "TCP",
		Host: "0.0.0.0",
		Port: 2600,
	}
	return aisd
}

//  初始化
func (aisd *AISDevice) Init(devId string, configMap map[string]interface{}) error {
	aisd.PointId = devId
	if err := utils.BindSourceConfig(configMap, &aisd.mainConfig); err != nil {
		return err
	}

	return nil
}

// 启动
func (aisd *AISDevice) Start(cctx typex.CCTX) error {
	aisd.Ctx = cctx.Ctx
	aisd.CancelCTX = cctx.CancelCTX
	//
	listener, err := net.Listen("tcp",
		fmt.Sprintf("%s:%v", aisd.mainConfig.Host, aisd.mainConfig.Port))
	if err != nil {
		return err
	}
	aisd.tcpListener = listener
	glogger.GLogger.Infof("AIS TCP server started on TCP://%s:%v",
		aisd.mainConfig.Host, aisd.mainConfig.Port)
	go aisd.handleConnect(listener)
	return nil
}

// 从设备里面读数据出来
func (aisd *AISDevice) OnRead(cmd []byte, data []byte) (int, error) {
	return 0, nil
}

// 把数据写入设备
func (aisd *AISDevice) OnWrite(cmd []byte, _ []byte) (int, error) {
	return 0, nil
}

// 设备当前状态
func (aisd *AISDevice) Status() typex.DeviceState {
	return aisd.status
}

// 停止设备
func (aisd *AISDevice) Stop() {
	aisd.status = typex.DEV_DOWN
	aisd.CancelCTX()
	if aisd.tcpListener != nil {
		aisd.tcpListener.Close()
	}
}

// 设备属性，是一系列属性描述
func (aisd *AISDevice) Property() []typex.DeviceProperty {
	return []typex.DeviceProperty{}
}

// 真实设备
func (aisd *AISDevice) Details() *typex.Device {
	return aisd.RuleEngine.GetDevice(aisd.PointId)
}

// 状态
func (aisd *AISDevice) SetState(status typex.DeviceState) {
	aisd.status = status

}

// 驱动
func (aisd *AISDevice) Driver() typex.XExternalDriver {
	return nil
}

func (aisd *AISDevice) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}

/*
*
* OnCtrl 接口可以用来向外广播数据
*
 */
func (aisd *AISDevice) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	return []byte{}, nil
}

//--------------------------------------------------------------------------------------------------
// 内部
//--------------------------------------------------------------------------------------------------
/*
*
* 处理连接
*
 */
func (aisd *AISDevice) handleConnect(listener net.Listener) {
	for {
		select {
		case <-aisd.Ctx.Done():
			{
				return
			}
		default:
			{
			}
		}
		tcpcon, err := listener.Accept()
		if err != nil {
			glogger.GLogger.Error(err)
			continue
		}
		go aisd.handleIO(tcpcon)
	}

}

/*
*
* 数据处理
*
 */
func (aisd *AISDevice) handleIO(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	s, err := reader.ReadString('\n')
	if err != nil {
		glogger.GLogger.Error(err)
		return
	}
	if strings.HasSuffix(s, "\r\n") {
		glogger.GLogger.Debug("Received data:", s)
	}
}
