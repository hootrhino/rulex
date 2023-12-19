package device

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"sync"
	"time"

	siemenscache "github.com/hootrhino/rulex/component/intercache/siemens"

	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/component/interdb"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
	"github.com/robinson/gos7"
)

// 点位表
type SiemensDataPoint struct {
	UUID       string `json:"uuid"` // 当UUID为空时新建
	DeviceUuid string `json:"device_uuid"`
	Tag        string `json:"tag,omitempty"`
	Type       string `json:"type,omitempty"`
	Frequency  *int64 `json:"frequency,omitempty"`
	Address    *int   `json:"address,omitempty"`
	Start      *int   `json:"start"`
	Size       *int   `json:"size"`
	Value      string `json:"value"`
}

type S1200CommonConfig struct {
	Host  string `json:"host" validate:"required"`  // 127.0.0.1:502
	Model string `json:"model" validate:"required"` // s7-200 s7-1500
	// https://cloudvpn.beijerelectronics.com/hc/en-us/articles/4406049761169-Siemens-S7
	Rack        *int  `json:"rack" validate:"required"`        // 0
	Slot        *int  `json:"slot" validate:"required"`        // 1
	Timeout     *int  `json:"timeout" validate:"required"`     // 5s
	IdleTimeout *int  `json:"idleTimeout" validate:"required"` // 5s
	AutoRequest *bool `json:"autoRequest" validate:"required"` // false
}
type S1200Config struct {
	CommonConfig S1200CommonConfig `json:"commonConfig" validate:"required"` // 通用配置
}

// https://www.ad.siemens.com.cn/productportal/prods/s7-1200_plc_easy_plus/07-Program/02-basic/01-Data_Type/01-basic.html
type SIEMENS_PLC struct {
	typex.XStatus
	status            typex.DeviceState
	RuleEngine        typex.RuleX
	mainConfig        S1200Config
	client            gos7.Client
	handler           *gos7.TCPClientHandler
	lock              sync.Mutex
	SiemensDataPoints map[string]*SiemensDataPoint
}

/*
*
* 西门子 S1200 系列 PLC
*
 */
func NewSIEMENS_PLC(e typex.RuleX) typex.XDevice {
	s1200 := new(SIEMENS_PLC)
	s1200.RuleEngine = e
	s1200.lock = sync.Mutex{}
	Rack := 0
	Slot := 1
	Timeout := 1000
	IdleTimeout := 3000
	AutoRequest := false
	s1200.mainConfig = S1200Config{
		CommonConfig: S1200CommonConfig{
			Rack:        &Rack,
			Slot:        &Slot,
			Timeout:     &Timeout,
			IdleTimeout: &IdleTimeout,
			AutoRequest: &AutoRequest,
		},
	}
	s1200.SiemensDataPoints = map[string]*SiemensDataPoint{}
	return s1200
}

// 初始化
func (s1200 *SIEMENS_PLC) Init(devId string, configMap map[string]interface{}) error {
	s1200.PointId = devId
	siemenscache.RegisterSlot(s1200.PointId)
	if err := utils.BindSourceConfig(configMap, &s1200.mainConfig); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	// 合并数据库里面的点位表
	// TODO 这里需要优化一下，而不是直接查表这种形式，应该从物模型组件来加载
	// DataSchema = schema.load(uuid)
	// DataSchema.update(k, v)
	var list []SiemensDataPoint
	errDb := interdb.DB().Table("m_siemens_data_points").
		Where("device_uuid=?", devId).Find(&list).Error
	if errDb != nil {
		return errDb
	}
	for _, v := range list {
		// 频率不能太快
		if *v.Frequency < 50 {
			return errors.New("'frequency' must grate than 50 millisecond")
		}
		s1200.SiemensDataPoints[v.UUID] = &v
		siemenscache.SetValue(s1200.PointId, v.UUID, siemenscache.SiemensPoint{
			UUID:          v.UUID,
			Status:        0,
			LastFetchTime: 0,
			Value:         "0",
		})
	}
	return nil
}

// 启动
func (s1200 *SIEMENS_PLC) Start(cctx typex.CCTX) error {
	s1200.Ctx = cctx.Ctx
	s1200.CancelCTX = cctx.CancelCTX
	//
	s1200.handler = gos7.NewTCPClientHandler(
		s1200.mainConfig.CommonConfig.Host,  // 127.0.0.1:1500
		*s1200.mainConfig.CommonConfig.Rack, // 0
		*s1200.mainConfig.CommonConfig.Slot) // 1
	s1200.handler.Timeout = time.Duration(
		*s1200.mainConfig.CommonConfig.Timeout) * time.Millisecond
	s1200.handler.IdleTimeout = time.Duration(
		*s1200.mainConfig.CommonConfig.IdleTimeout) * time.Millisecond
	if err := s1200.handler.Connect(); err != nil {
		return err
	} else {
		s1200.status = typex.DEV_UP
	}

	s1200.client = gos7.NewClient(s1200.handler)
	if !*s1200.mainConfig.CommonConfig.AutoRequest {
		s1200.status = typex.DEV_UP
		return nil
	}
	go func(ctx context.Context) {
		// 数据缓冲区,最大4KB
		dataBuffer := make([]byte, common.T_4KB)
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			default:
				{
				}
			}
			s1200.lock.Lock()
			// CMD 参数无用
			n, err := s1200.Read([]byte(""), dataBuffer)
			s1200.lock.Unlock()
			if err != nil {
				glogger.GLogger.Error(err)
				s1200.status = typex.DEV_DOWN
				return
			}
			// [] {} ""
			if n < 3 {
				continue
			}
			ok, err := s1200.RuleEngine.WorkDevice(
				s1200.RuleEngine.GetDevice(s1200.PointId),
				string(dataBuffer[:n]),
			)
			// glogger.GLogger.Debug(string(dataBuffer[:n]))
			if !ok {
				glogger.GLogger.Error(err)
			}
		}

	}(cctx.Ctx)
	return nil
}

// 从设备里面读数据出来
func (s1200 *SIEMENS_PLC) OnRead(cmd []byte, data []byte) (int, error) {
	return s1200.Read(cmd, data)
}

// 把数据写入设备
//
// db.Address:int, db.Start:int, db.Size:int, rData[]

func (s1200 *SIEMENS_PLC) OnWrite(cmd []byte, data []byte) (int, error) {
	blocks := []SiemensDataPoint{}
	if err := json.Unmarshal(data, &blocks); err != nil {
		return 0, err
	}
	return s1200.Write(cmd, data)
}

// 设备当前状态
func (s1200 *SIEMENS_PLC) Status() typex.DeviceState {
	if s1200.client == nil {
		return typex.DEV_DOWN
	}
	return s1200.status

}

// 停止设备
func (s1200 *SIEMENS_PLC) Stop() {
	s1200.status = typex.DEV_DOWN
	if s1200.CancelCTX != nil {
		s1200.CancelCTX()
	}
	if s1200.handler != nil {
		s1200.handler.Close()
	}
	siemenscache.UnRegisterSlot(s1200.PointId)
}

// 设备属性，是一系列属性描述
func (s1200 *SIEMENS_PLC) Property() []typex.DeviceProperty {
	return []typex.DeviceProperty{}
}

// 真实设备
func (s1200 *SIEMENS_PLC) Details() *typex.Device {
	return s1200.RuleEngine.GetDevice(s1200.PointId)
}

// 状态
func (s1200 *SIEMENS_PLC) SetState(status typex.DeviceState) {
	s1200.status = status
}

// 驱动
func (s1200 *SIEMENS_PLC) Driver() typex.XExternalDriver {
	return nil
}

func (s1200 *SIEMENS_PLC) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}
func (s1200 *SIEMENS_PLC) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	return []byte{}, nil
}
func (s1200 *SIEMENS_PLC) Write(cmd []byte, data []byte) (int, error) {
	return 0, nil
}

// 字节格式:[dbNumber1, start1, size1, dbNumber2, start2, size2]
// 读: db --> dbNumber, start, size, buffer[]
var rData = [common.T_2KB]byte{} // 一次最大接受2KB数据

func (s1200 *SIEMENS_PLC) Read(cmd []byte, data []byte) (int, error) {
	values := []SiemensDataPoint{}
	for uuid, db := range s1200.SiemensDataPoints {
		//DB 4字节
		if db.Type == "DB" {
			// 00.00.00.01 | 00.00.00.02 | 00.00.00.03 | 00.00.00.04
			if err := s1200.client.AGReadDB(*db.Address, *db.Start, *db.Size, rData[:]); err != nil {
				glogger.GLogger.Error(err)
				return 0, err
			}
			count := db.Size
			if *db.Size*2 > 2000 {
				*count = 2000
			}
			Value := hex.EncodeToString(rData[:*count])
			values = append(values, SiemensDataPoint{
				DeviceUuid: db.DeviceUuid,
				Tag:        db.Tag,
				Address:    db.Address,
				Type:       db.Type,
				Start:      db.Start,
				Size:       db.Size,
				Value:      Value,
			})
			siemenscache.SetValue(s1200.PointId, uuid, siemenscache.SiemensPoint{
				UUID:          uuid,
				Status:        0,
				LastFetchTime: uint64(time.Now().UnixMilli()),
				Value:         Value,
			})
		}
		//
		if db.Type == "MB" {
			// 00.00.00.01 | 00.00.00.02 | 00.00.00.03 | 00.00.00.04
			if err := s1200.client.AGReadMB(*db.Start, *db.Size, rData[:]); err != nil {
				glogger.GLogger.Error(err)
				return 0, err
			}
			count := db.Size
			if *db.Size*2 > 2000 {
				*count = 2000
			}
			Value := hex.EncodeToString(rData[:*count])
			values = append(values, SiemensDataPoint{
				Tag:     db.Tag,
				Type:    db.Type,
				Address: db.Address,
				Start:   db.Start,
				Size:    db.Size,
				Value:   Value,
			})
			siemenscache.SetValue(s1200.PointId, uuid, siemenscache.SiemensPoint{
				UUID:          uuid,
				Status:        0,
				LastFetchTime: uint64(time.Now().UnixMilli()),
				Value:         Value,
			})
		}
		if *db.Frequency < 100 {
			*db.Frequency = 100 // 不能太快
		}
		time.Sleep(time.Duration(*db.Frequency) * time.Millisecond)
	}
	bytes, _ := json.Marshal(values)
	copy(data, bytes)
	return len(bytes), nil
}

/*
*
- 符号地址（Symbolic Addressing）：
使用符号名称来表示变量或输入/输出地址。这种方式更加直观和易于理解，适用于高级编程语言和工程师使用。例如，可以使用变量名"MotorSpeed"或输入名"I1"来表示对应的地址。
- 基于字节的地址（Byte-based Addressing）：
使用字节地址和位地址的组合来表示变量或输入/输出地址。字节地址表示内存中的字节偏移，而位地址表示字节中的位偏移。例如，使用地址"DB1.DBX10.3"表示数据块1中偏移为10的字节的第3位。
- 基于字的地址（Word-based Addressing）：
类似于基于字节的地址，但是将地址表示为字（16位）的偏移。例如，使用地址"DB1.DBD20"表示数据块1中偏移为20的字。
- 基于地址区域的地址（Address Area-based Addressing）：
将地址按照不同的区域进行划分，如输入区域（I），输出区域（Q），数据块区域（DB）等。每个区域都有特定的地址范围。例如，使用地址"I10.3"表示输入区域的第10个输入的第3位。
*
*/

// AddressInfo 包含解析后的地址信息
type AddressInfo struct {
	DataBlockNumber int    // 数据块号
	DataType        string // 数据类型
	ElementNumber   int    // 元素号
}

// 解析DB
func ParseDB_D(s string) string {
	return ""
}

// 解析DBX格式
func ParseDB_X(s string) string {
	return ""
}

// 解析I格式
func ParseADDR_I(s string) string {
	return ""
}

// 解析Q格式
func ParseADDR_Q(s string) string {
	return ""
}
