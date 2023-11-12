package device

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"reflect"
	"strings"
	"time"

	"github.com/adrianmo/go-nmea"
	aislib "github.com/hootrhino/go-ais"
	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/component/hwportmanager"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
	"github.com/jinzhu/copier"
	serial "github.com/wwhai/goserial"
)

var __AisCodec = aislib.CodecNew(false, false, false)

// --------------------------------------------------------------------------------------------------
// 把AIS包里面的几个结构体拿出来了，主要是适配JSON格式, 下面这些结构体和AIS包里面的完全一样
// --------------------------------------------------------------------------------------------------
type _AISCommonConfig struct {
	Mode     string `json:"mode" title:"工作模式" info:"UART/TCP"`
	ParseAis bool   `json:"parseAis"`
	GwSN     string `json:"gwsn"`
}
type _AISDeviceMasterConfig struct {
	CommonConfig _AISCommonConfig  `json:"commonConfig"`
	HostConfig   common.HostConfig `json:"hostConfig"`
	PortUuid     string            `json:"portUuid"`
}
type AISDeviceMaster struct {
	typex.XStatus
	status       typex.DeviceState
	mainConfig   _AISDeviceMasterConfig
	RuleEngine   typex.RuleX
	tcpListener  net.Listener // TCP 接收端
	serialPort   serial.Port
	hwPortConfig hwportmanager.UartConfig

	// session
	DevicesSessionMap map[string]*__AISDeviceSession
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
		HostConfig: common.HostConfig{
			Host:    "127.0.0.1",
			Port:    6005,
			Timeout: 3000,
		},
		CommonConfig: _AISCommonConfig{
			Mode:     "TCP",
			ParseAis: false,
			GwSN:     "HR0001",
		},
	}
	aism.DevicesSessionMap = map[string]*__AISDeviceSession{}
	return aism
}

//  初始化
func (aism *AISDeviceMaster) Init(devId string, configMap map[string]interface{}) error {
	aism.PointId = devId
	if err := utils.BindSourceConfig(configMap, &aism.mainConfig); err != nil {
		return err
	}
	if !utils.SContains([]string{"UART", "TCP"}, aism.mainConfig.CommonConfig.Mode) {
		return errors.New("unsupported mode, only can be one of 'TCP' or 'RTU'")
	}

	if aism.mainConfig.CommonConfig.Mode == "UART" {
		hwPort, err := hwportmanager.GetHwPort(aism.mainConfig.PortUuid)
		if err != nil {
			return err
		}
		if hwPort.Busy {
			return fmt.Errorf("UART is busying now, Occupied By:%s", hwPort.OccupyBy)
		}
		switch tCfg := hwPort.Config.(type) {
		case hwportmanager.UartConfig:
			{
				aism.hwPortConfig = tCfg
			}
		default:
			{
				return fmt.Errorf("invalid config:%s", hwPort.Config)
			}
		}
	}

	return nil
}

// 启动
func (aism *AISDeviceMaster) Start(cctx typex.CCTX) error {
	aism.Ctx = cctx.Ctx
	aism.CancelCTX = cctx.CancelCTX
	if aism.mainConfig.CommonConfig.Mode == "TCP" {
		//
		listener, err := net.Listen("tcp",
			fmt.Sprintf("%s:%v", aism.mainConfig.HostConfig.Host, aism.mainConfig.HostConfig.Port))
		if err != nil {
			return err
		}
		aism.tcpListener = listener
		glogger.GLogger.Infof("AIS TCP server started on TCP://%s:%v",
			aism.mainConfig.HostConfig.Host, aism.mainConfig.HostConfig.Port)
		go aism.handleTcpConnect(listener)
		aism.status = typex.DEV_UP
		return nil
	}
	// 串口收发卡
	if aism.mainConfig.CommonConfig.Mode == "UART" {
		config := serial.Config{
			Address:  aism.hwPortConfig.Uart,
			BaudRate: aism.hwPortConfig.BaudRate,
			DataBits: aism.hwPortConfig.DataBits,
			Parity:   aism.hwPortConfig.Parity,
			StopBits: aism.hwPortConfig.StopBits,
			Timeout:  time.Duration(aism.hwPortConfig.Timeout) * time.Millisecond,
		}
		var err error
		aism.serialPort, err = serial.Open(&config)
		if err != nil {
			glogger.GLogger.Error("serial port start failed err:", err, ", config:", config)
			return err
		}
		go func() {
			buffer := [4096]byte{}
			defer aism.serialPort.Close()
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
				offset := 0
				endl1 := false
				endl2 := false
				ok := false
				oneByte := [1]byte{}
				readyStatus := false // 超时也是就绪状态
				ctx1, cancel1 := context.WithTimeout(aism.Ctx,
					time.Duration(aism.hwPortConfig.Timeout)*time.Millisecond)
				defer cancel1()
				for {
					select {
					// 控制时间防止死机
					case <-ctx1.Done():
						{
							if !readyStatus {
								glogger.GLogger.Warnf("serialPort %s Read timeout", aism.hwPortConfig.Uart)
							}
							break
						}
					default:
						{
						}
					}
					_, err := aism.serialPort.Read(oneByte[:])
					if err != nil {
						if strings.Contains(err.Error(), "timeout") {
							readyStatus = true
							continue
						}
						readyStatus = false
						aism.status = typex.DEV_DOWN
						glogger.GLogger.Errorf("serialPort %s Read error", aism.hwPortConfig.Uart)
						return
					}
					if oneByte[0] == '\r' {
						endl1 = true
						continue
					}
					if oneByte[0] == '\n' {
						endl2 = true
						ok = true
					}
					if endl1 && endl2 {
						break
					} else {
						buffer[offset] = oneByte[0]
						offset++
					}
				}
				// 可能AIS报文传输失败了
				if !ok {
					glogger.GLogger.Info("serialPort Read may occurred error:", err)
					continue
				}
				rawAiSString := string(buffer[:offset])
				if err != nil {
					glogger.GLogger.Error(err)
					aism.status = typex.DEV_DOWN
					return
				}
				// 这段是个兼容代码，现阶段适配了一款AIS USB 串口接收器，以后会自己做
				{
					if strings.HasPrefix("NONE", rawAiSString) {
						glogger.GLogger.Info("AIS33VRx Receiver Heart Beat Packet")
						continue
					}
					if strings.HasPrefix("AIS33VRx", rawAiSString) {
						glogger.GLogger.Info(rawAiSString)
						continue
					}
					if strings.HasPrefix("AIS Ch 1", rawAiSString) {
						glogger.GLogger.Info(rawAiSString)
						continue
					}
					if strings.HasPrefix("AIS Ch 2", rawAiSString) {
						glogger.GLogger.Info(rawAiSString)
						continue
					}
				}

				// 如果不需要解析,直接原文透传
				if !aism.mainConfig.CommonConfig.ParseAis {
					// {
					//     "ais_receiver_device":"%s",
					//     "gwsn":"%s"
					//     "ais_data":"%s"
					// }
					ds := `{"ais_receiver_device":"%s","gwsn":"%s","ais_data":"%s"}`
					lens := len(rawAiSString)
					if lens > 2 {
						aism.RuleEngine.WorkDevice(aism.Details(),
							fmt.Sprintf(ds, aism.mainConfig.CommonConfig.GwSN,
								aism.mainConfig.CommonConfig.GwSN, rawAiSString), // \r\n
						)
					}
				}
			}
		}()
		hwportmanager.SetInterfaceBusy(aism.mainConfig.PortUuid, hwportmanager.HwPortOccupy{
			UUID: aism.PointId,
			Type: "DEVICE",
			Name: aism.Details().Name,
		})
		aism.status = typex.DEV_UP
		return nil
	}
	aism.status = typex.DEV_DOWN
	return fmt.Errorf("invalid work mode:%s", aism.mainConfig.CommonConfig.Mode)
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
	if aism.CancelCTX != nil {
		aism.CancelCTX()
	}
	if aism.tcpListener != nil {
		aism.tcpListener.Close()
	}
	// release serial port
	if aism.mainConfig.CommonConfig.Mode == "UART" {
		if aism.serialPort != nil {
			aism.serialPort.Close()
		}
		hwportmanager.FreeInterfaceBusy(aism.mainConfig.PortUuid)
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
func (aism *AISDeviceMaster) handleTcpConnect(listener net.Listener) {
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
		go aism.handleTcpAuth(ctx, cancel, &__AISDeviceSession{
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
type __AISDeviceSession struct {
	SN        string   // 注册包里的序列号, 必须是:SN-$AA-$BB-$CC-$DD
	Ip        string   // 注册包里的序列号
	Transport net.Conn // TCP连接
}

func (aism *AISDeviceMaster) handleTcpAuth(ctx context.Context,
	cancel context.CancelFunc, session *__AISDeviceSession) {
	// 5秒内读一个SN
	session.Transport.SetDeadline(time.Now().Add(5 * time.Second))
	reader := bufio.NewReader(session.Transport)
	registerPkt, err := reader.ReadString('$')
	session.Transport.SetDeadline(time.Time{})
	//
	if err != nil {
		glogger.GLogger.Error(session.Transport.RemoteAddr(), err)
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
func (aism *AISDeviceMaster) handleIO(session *__AISDeviceSession) {
	reader := bufio.NewReader(session.Transport)
	for {
		rawAiSString, err := reader.ReadString('\n')
		if err != nil {
			glogger.GLogger.Error(err)
			delete(aism.DevicesSessionMap, session.SN)
			session.Transport.Close()
			aism.status = typex.DEV_DOWN
			return
		}
		// 如果不需要解析,直接原文透传
		if !aism.mainConfig.CommonConfig.ParseAis {
			// {
			//     "ais_receiver_device":"%s",
			//     "ais_data":"%s"
			//     "gwsn":"%s"
			// }
			ds := `{"ais_receiver_device":"%s","gwsn":"%s","ais_data":"%s"}`
			aism.RuleEngine.WorkDevice(aism.Details(), fmt.Sprintf(ds, session.SN,
				aism.mainConfig.CommonConfig.GwSN, rawAiSString))
			continue
		}
		// 可能会收到心跳包: !HRT710,Q,003,0*06
		if strings.HasPrefix(rawAiSString, "!HRT") {
			glogger.GLogger.Debug("Heart beat from:", session.SN, session.Transport.RemoteAddr())
			continue
		}
		if strings.HasPrefix(rawAiSString, "!DYA") {
			glogger.GLogger.Debug("DYA Message from:", session.SN, session.Transport.RemoteAddr())
			continue
		}
		sentence, err := nmea.Parse(rawAiSString)
		if err != nil {
			glogger.GLogger.Error(err, rawAiSString)
			continue
		}
		// glogger.GLogger.Info("Received data:", sentence.DataType(), sentence)
		if sentence.DataType() == nmea.TypeRMC {
			rmc1 := sentence.(nmea.RMC)
			rmc := RMC{}
			copier.Copy(&rmc, &rmc1)
			data := rmc.String()
			glogger.GLogger.Debug("Received RMC data:", data)
			if data != "" {
				aism.RuleEngine.WorkDevice(aism.Details(), data)
			}
		}
		if sentence.DataType() == nmea.TypeGNS {
			gns1 := sentence.(nmea.GNS)
			gns := GNS{}
			copier.Copy(&gns, &gns1)
			data := gns.String()
			glogger.GLogger.Debug("Received GNS data:", data)
			if data != "" {
				aism.RuleEngine.WorkDevice(aism.Details(), data)
			}
		}
		if sentence.DataType() == nmea.TypeVDM {
			vdmo1 := sentence.(nmea.VDMVDO)
			vdmo := VDMVDO{}
			copier.Copy(&vdmo, &vdmo1)
			data := vdmo.PayloadInfo()
			glogger.GLogger.Debug("Received VDM data:", data)
			if data != "" {
				aism.RuleEngine.WorkDevice(aism.Details(), data)
			}
		}
		if sentence.DataType() == nmea.TypeVDO {
			vdmo1 := sentence.(nmea.VDMVDO)
			vdmo := VDMVDO{}
			copier.Copy(&vdmo, &vdmo1)
			data := vdmo.PayloadInfo()
			glogger.GLogger.Debug("Received VDO data:", data)
			if data != "" {
				aism.RuleEngine.WorkDevice(aism.Details(), data)
			}
		}

	}

}

//--------------------------------------------------------------------------------------------------
// AIS 结构, 下面这些结构是从nema包里面拿过来的，删除了一些无用字段，主要为了方便JSON编码操作
//--------------------------------------------------------------------------------------------------

type BaseSentence struct {
	Talker string `json:"talker"` // The talker id (e.g GP)
	Type   string `json:"type"`   // The data type (e.g GSA)
}

// Prefix returns the talker and type of message
func (s BaseSentence) Prefix() string {
	return s.Talker + s.Type
}

// DataType returns the type of the message
func (s BaseSentence) DataType() string {
	return s.Type
}

// TalkerID returns the talker of the message
func (s BaseSentence) TalkerID() string {
	return s.Talker
}

type TagBlock struct {
	Time         int64  `json:"time"`          // TypeUnixTime unix timestamp (unit is likely to be s, but might be ms, YMMV), parameter: -c
	RelativeTime int64  `json:"relative_time"` // TypeRelativeTime relative time, parameter: -r
	Destination  string `json:"destination"`   // TypeDestinationID destination identification 15 char max, parameter: -d
	Grouping     string `json:"grouping"`      // TypeGrouping sentence grouping, parameter: -g
	LineCount    int64  `json:"line_count"`    // TypeLineCount line count, parameter: -n
	Source       string `json:"source"`        // TypeSourceID source identification 15 char max, parameter: -s
	Text         string `json:"text"`          // TypeTextString valid character string, parameter -t
}
type RMC struct {
	BaseSentence `json:"base"` // base
	Time         Time          `json:"time"`       // Time Stamp
	Validity     string        `json:"validity"`   // validity - A-ok, V-invalid
	Latitude     float64       `json:"latitude"`   // Latitude
	Longitude    float64       `json:"longitude"`  // Longitude
	Speed        float64       `json:"speed"`      // Speed in knots
	Course       float64       `json:"course"`     // True course
	Date         Date          `json:"date"`       // Date
	Variation    float64       `json:"variation"`  // Magnetic variation
	FFAMode      string        `json:"ffa_mode"`   // FAA mode indicator (filled in NMEA 2.3 and later)
	NavStatus    string        `json:"nav_status"` // Nav Status (NMEA 4.1 and later)
}

func (s RMC) String() string {
	bytes, err := json.Marshal(s)
	if err != nil {
		return ""
	}
	return string(bytes)
}

/*
*
* AIS消息结构体
*
 */
type VDMVDO struct {
	BaseSentence   `json:"base"`
	NumFragments   int64         `json:"numFragments"`
	FragmentNumber int64         `json:"fragmentNumber"`
	MessageID      int64         `json:"messageId"`
	Channel        string        `json:"channel"`
	Payload        []byte        `json:"-"`
	MessageContent aislib.Packet `json:"messageContent"`
}
type __PositionReport struct {
	MessageID uint8   `json:"message_id,omitempty"`
	UserID    uint32  `json:"user_id,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
	Timestamp uint8   `json:"timestamp,omitempty"`
}

// StaticDataReportA is the A part of message 24
type __StaticDataReportA struct {
	Valid bool   `json:"valid,omitempty"`
	Name  string `json:"name,omitempty"`
}

// StaticDataReportB is the B part of message 24
type __StaticDataReportB struct {
	Valid          bool   `json:"valid,omitempty"`
	ShipType       uint8  `json:"ship_type,omitempty"`
	VendorIDName   string `json:"vendor_id_name,omitempty"`
	VenderIDModel  uint8  `json:"vender_id_model,omitempty"`
	VenderIDSerial uint32 `json:"vender_id_serial,omitempty"`
	CallSign       string `json:"call_sign,omitempty"`
	FixType        uint8  `json:"fix_type,omitempty"`
	Spare          uint8  `json:"spare,omitempty"`
}

type __StaticDataReport struct {
	MessageID  uint8               `json:"message_id,omitempty"`
	UserID     uint32              `json:"user_id,omitempty"`
	Valid      bool                `json:"valid,omitempty"`
	Reserved   uint8               `json:"reserved,omitempty"`
	PartNumber bool                `json:"part_number,omitempty"`
	ReportA    __StaticDataReportA `json:"report_a,omitempty"`
	ReportB    __StaticDataReportB `json:"report_b,omitempty"`
}

func (v VDMVDO) PayloadInfo() string {
	__AisCodec.DropSpace = true
	pkt := __AisCodec.DecodePacket(v.Payload)
	// aislib.StandardClassBPositionReport
	var _Type reflect.Type
	if _Type = reflect.TypeOf(pkt); _Type == nil {
		return ""
	}
	// 上报位置
	if _Type.Name() == "StandardClassBPositionReport" {
		spr := pkt.(aislib.StandardClassBPositionReport)
		pos := __PositionReport{}
		copier.Copy(&pos, &spr)
		bytes, _ := json.Marshal(pos)
		return string(bytes)
	}
	// "StaticDataReport"
	if _Type.Name() == "StaticDataReport" {
		spr := pkt.(aislib.StaticDataReport)
		data := __StaticDataReport{}
		copier.Copy(&data, &spr)
		bytes, _ := json.Marshal(data)
		return string(bytes)
	}
	return ""
}
func (s VDMVDO) String() string {
	bytes, _ := json.Marshal(s)
	return string(bytes)
}

type Time struct {
	Valid       bool `json:"valid"`
	Hour        int  `json:"hour"`
	Minute      int  `json:"minute"`
	Second      int  `json:"second"`
	Millisecond int  `json:"millisecond"`
}

// String representation of Time
func (t Time) String() string {
	seconds := float64(t.Second) + float64(t.Millisecond)/1000
	return fmt.Sprintf("%02d:%02d:%07.4f", t.Hour, t.Minute, seconds)
}

// Date type
type Date struct {
	Valid bool `json:"valid"`
	DD    int  `json:"dd"`
	MM    int  `json:"mm"`
	YY    int  `json:"yy"`
}

// String representation of date
func (d Date) String() string {
	return fmt.Sprintf("%02d/%02d/%02d", d.DD, d.MM, d.YY)
}

type GNS struct {
	BaseSentence
	Time      Time // UTC of position
	Latitude  float64
	Longitude float64
	// FAA mode indicator for each satellite navigation system (constellation) supported by device.
	//
	// May be up to six characters (according to GPSD).
	// '1' - GPS
	// '2' - GLONASS
	// '3' - Galileo
	// '4' - BDS
	// '5' - QZSS
	// '6' - NavIC (IRNSS)
	Mode       []string
	SVs        int64   // Total number of satellites in use, 00-99
	HDOP       float64 // Horizontal Dilution of Precision
	Altitude   float64 // Antenna altitude, meters, re:mean-sea-level(geoid).
	Separation float64 // Geoidal separation meters
	Age        float64 // Age of differential data
	Station    int64   // Differential reference station ID
	NavStatus  string  // Navigation status (NMEA 4.1+). See NavStats* (`NavStatusAutonomous` etc) constants for possible values.
}

func (s GNS) String() string {
	bytes, _ := json.Marshal(s)
	return string(bytes)
}
