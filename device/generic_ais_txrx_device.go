package device

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"time"

	"github.com/adrianmo/go-nmea"
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
type _AISDeviceMasterBaseSentence struct {
	Talker   string                   `json:"talker,omitempty"`   // The talker id (e.g GP)
	Type     string                   `json:"type,omitempty"`     // The data type (e.g GSA)
	Fields   []string                 `json:"fields,omitempty"`   // Array of fields
	Checksum string                   `json:"checksum,omitempty"` // The Checksum
	Raw      string                   `json:"raw,omitempty"`      // The raw NMEA sentence received
	TagBlock _AISDeviceMasterTagBlock `json:"tagBlock,omitempty"` // NMEA tagblock
}

/*
*
* Tag
*
 */
type _AISDeviceMasterTagBlock struct {
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
type _AISDeviceMasterPacket struct {
	_AISDeviceMasterBaseSentence
	NumFragments   int64  `json:"numFragments"`
	FragmentNumber int64  `json:"fragmentNumber"`
	MessageID      int64  `json:"messageID"`
	Channel        string `json:"channel"`
	Payload        []byte `json:"payload"`
}

// --------------------------------------------------------------------------------------------------
type _AISDeviceMasterConfig struct {
	Mode string `json:"mode"` // TCP UDP UART
	Host string `json:"host" validate:"required"`
	Port int    `json:"port" validate:"required"`
}
type AISDeviceMaster struct {
	typex.XStatus
	status      typex.DeviceState
	mainConfig  _AISDeviceMasterConfig
	RuleEngine  typex.RuleX
	tcpListener net.Listener // TCP 接收端
	// session
	DevicesSessionMap map[string]*AISDeviceSession
}

/*
*
* AIS 数据解析服务器
*
 */
func NewAISDeviceMaster(e typex.RuleX) typex.XDevice {
	aism := new(AISDeviceMaster)
	aism.RuleEngine = e
	aism.mainConfig = _AISDeviceMasterConfig{
		Mode: "TCP",
		Host: "0.0.0.0",
		Port: 2600,
	}
	aism.DevicesSessionMap = map[string]*AISDeviceSession{}
	return aism
}

//  初始化
func (aism *AISDeviceMaster) Init(devId string, configMap map[string]interface{}) error {
	aism.PointId = devId
	if err := utils.BindSourceConfig(configMap, &aism.mainConfig); err != nil {
		return err
	}

	return nil
}

// 启动
func (aism *AISDeviceMaster) Start(cctx typex.CCTX) error {
	aism.Ctx = cctx.Ctx
	aism.CancelCTX = cctx.CancelCTX
	//
	listener, err := net.Listen("tcp",
		fmt.Sprintf("%s:%v", aism.mainConfig.Host, aism.mainConfig.Port))
	if err != nil {
		return err
	}
	aism.tcpListener = listener
	glogger.GLogger.Infof("AIS TCP server started on TCP://%s:%v",
		aism.mainConfig.Host, aism.mainConfig.Port)
	go aism.handleConnect(listener)
	return nil
}

// 从设备里面读数据出来
func (aism *AISDeviceMaster) OnRead(cmd []byte, data []byte) (int, error) {
	return 0, nil
}

// 把数据写入设备
func (aism *AISDeviceMaster) OnWrite(cmd []byte, _ []byte) (int, error) {
	return 0, nil
}

// 设备当前状态
func (aism *AISDeviceMaster) Status() typex.DeviceState {
	return aism.status
}

// 停止设备
func (aism *AISDeviceMaster) Stop() {
	aism.status = typex.DEV_DOWN
	aism.CancelCTX()
	if aism.tcpListener != nil {
		aism.tcpListener.Close()
	}
}

// 设备属性，是一系列属性描述
func (aism *AISDeviceMaster) Property() []typex.DeviceProperty {
	return []typex.DeviceProperty{}
}

// 真实设备
func (aism *AISDeviceMaster) Details() *typex.Device {
	return aism.RuleEngine.GetDevice(aism.PointId)
}

// 状态
func (aism *AISDeviceMaster) SetState(status typex.DeviceState) {
	aism.status = status

}

// 驱动
func (aism *AISDeviceMaster) Driver() typex.XExternalDriver {
	return nil
}

func (aism *AISDeviceMaster) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}

/*
*
* OnCtrl 接口可以用来向外广播数据
*
 */
func (aism *AISDeviceMaster) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
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
func (aism *AISDeviceMaster) handleConnect(listener net.Listener) {
	for {
		select {
		case <-aism.Ctx.Done():
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
		ctx, cancel := context.WithCancel(aism.Ctx)
		go aism.handleAuth(ctx, cancel, &AISDeviceSession{
			Transport: tcpcon,
		})

	}

}

/*
*
* 等待认证: 传感器发送的第一个包必须为ID, 最大不能超过64字节
* 注意：Auth只针对AIS主机，来自AIS的数据只解析不做验证
*
 */
type AISDeviceSession struct {
	SN        string   // 注册包里的序列号, 必须是:SN-$AA-$BB-$CC-$DD
	Ip        string   // 注册包里的序列号
	Transport net.Conn // TCP连接
}

func (aism *AISDeviceMaster) handleAuth(ctx context.Context,
	cancel context.CancelFunc, session *AISDeviceSession) {
	// 5秒内读一个SN
	session.Transport.SetDeadline(time.Now().Add(5 * time.Second))
	reader := bufio.NewReader(session.Transport)
	registerPkt, err := reader.ReadString('$')
	session.Transport.SetDeadline(time.Time{})
	//
	if err != nil {
		glogger.GLogger.Error(err)
		session.Transport.Close()
		return
	}
	// 对SN有要求, 必须不少于4个字符
	if len(registerPkt) < 4 {
		glogger.GLogger.Error("Must have register packet and can not less than 4 character")
		session.Transport.Close()
		return
	}
	sn := registerPkt[:len(registerPkt)-1] // 去除$
	glogger.GLogger.Debug("AIS Device ready to auth:", sn)
	if aism.DevicesSessionMap[sn] != nil {
		glogger.GLogger.Error("SN Already Have Been Registered:", sn)
		session.Transport.Close()
		return
	}
	session.SN = sn
	session.Ip = session.Transport.RemoteAddr().String()
	aism.DevicesSessionMap[sn] = session
	go aism.handleIO(session)

}

/*
*
* 数据处理
*
 */
func (aism *AISDeviceMaster) handleIO(session *AISDeviceSession) {

	reader := bufio.NewReader(session.Transport)
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			glogger.GLogger.Error(err)
			delete(aism.DevicesSessionMap, session.SN)
			session.Transport.Close()
			return
		}
		sentence, err := nmea.Parse(s)
		if err != nil {
			glogger.GLogger.Error(err, s)
			continue
		}
		// glogger.GLogger.Info("Received data:", sentence.DataType(), sentence)
		if sentence.DataType() == nmea.TypeRMC {
			rmc := sentence.(nmea.RMC)
			glogger.GLogger.Info("Received RMC data:", rmc.String())
		}
		if sentence.DataType() == nmea.TypeVDM {
			vdmo := sentence.(nmea.VDMVDO)
			glogger.GLogger.Info("Received VDM data:", vdmo.String())
		}
	}

}
