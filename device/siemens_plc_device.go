package device

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	siemenscache "github.com/hootrhino/rulex/component/intercache/siemens"
	"github.com/hootrhino/rulex/component/iotschema"

	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/component/interdb"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
	"github.com/robinson/gos7"
)

// 点位表
type __SiemensDataPoint struct {
	UUID            string `json:"uuid"`
	DeviceUUID      string `json:"device_uuid"`
	SiemensAddress  string `json:"siemensAddress"` // 西门子的地址字符串
	Tag             string `json:"tag"`
	Alias           string `json:"alias"`
	Frequency       *int64 `json:"frequency"`
	Status          int    `json:"status"`          // 运行时数据
	LastFetchTime   uint64 `json:"lastFetchTime"`   // 运行时数据
	Value           string `json:"value"`           // 运行时数据
	AddressType     string `json:"addressType"`     // // 西门子解析后的地址信息: 寄存器类型: DB I Q
	DataBlockType   string `json:"dataBlockType"`   // // 西门子解析后的地址信息: 数据类型: INT UINT ....
	DataBlockOrder  string `json:"dataOrder"`       //  西门子解析后的地址信息: 数据类型: INT UINT ....
	DataBlockNumber int    `json:"dataBlockNumber"` // // 西门子解析后的地址信息: 数据块号: 100...
	ElementNumber   int    `json:"elementNumber"`   // // 西门子解析后的地址信息: 元素号:1000...
	DataSize        int    `json:"dataSize"`        // // 西门子解析后的地址信息: 位号,0-8，只针对I、Q
	BitNumber       int    `json:"bitNumber"`       // // 西门子解析后的地址信息: 位号,0-8，只针对I、Q
}

// https://cloudvpn.beijerelectronics.com/hc/en-us/articles/4406049761169-Siemens-S7
type S1200CommonConfig struct {
	Host        string `json:"host" validate:"required"`        // 127.0.0.1:502
	Model       string `json:"model" validate:"required"`       // s7-200 s7-1500
	Rack        *int   `json:"rack" validate:"required"`        // 0
	Slot        *int   `json:"slot" validate:"required"`        // 1
	Timeout     *int   `json:"timeout" validate:"required"`     // 5s
	IdleTimeout *int   `json:"idleTimeout" validate:"required"` // 5s
	AutoRequest *bool  `json:"autoRequest" validate:"required"` // false
}
type S1200Config struct {
	CommonConfig S1200CommonConfig `json:"commonConfig" validate:"required"` // 通用配置
}

// https://www.ad.siemens.com.cn/productportal/prods/s7-1200_plc_easy_plus/07-Program/02-basic/01-Data_Type/01-basic.html
type SIEMENS_PLC struct {
	typex.XStatus
	status              typex.DeviceState
	RuleEngine          typex.RuleX
	mainConfig          S1200Config
	client              gos7.Client
	handler             *gos7.TCPClientHandler
	lock                sync.Mutex
	__SiemensDataPoints map[string]*__SiemensDataPoint
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
	s1200.__SiemensDataPoints = map[string]*__SiemensDataPoint{}
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
	var list []__SiemensDataPoint
	errDb := interdb.DB().Table("m_siemens_data_points").
		Where("device_uuid=?", devId).Find(&list).Error
	if errDb != nil {
		return errDb
	}
	// 开始解析地址表
	for _, SiemensDataPoint := range list {
		// 频率不能太快
		if *SiemensDataPoint.Frequency < 50 {
			return errors.New("'frequency' must grate than 50 millisecond")
		}
		//
		AddressInfo, err1 := utils.ParseSiemensDB(SiemensDataPoint.SiemensAddress)
		if err1 != nil {
			return err1
		}
		SiemensDataPoint.DataBlockNumber = AddressInfo.DataBlockNumber
		SiemensDataPoint.ElementNumber = AddressInfo.ElementNumber
		SiemensDataPoint.AddressType = AddressInfo.AddressType
		SiemensDataPoint.BitNumber = AddressInfo.BitNumber
		SiemensDataPoint.DataSize = AddressInfo.DataBlockSize

		// 提前缓冲
		s1200.__SiemensDataPoints[SiemensDataPoint.UUID] = &SiemensDataPoint
		siemenscache.SetValue(s1200.PointId, SiemensDataPoint.UUID, siemenscache.SiemensPoint{
			UUID:          SiemensDataPoint.UUID,
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
	blocks := []__SiemensDataPoint{}
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
func (s1200 *SIEMENS_PLC) Property() []iotschema.IoTSchema {
	return []iotschema.IoTSchema{}
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
	values := []__SiemensDataPoint{}
	for uuid, db := range s1200.__SiemensDataPoints {
		//DB 4字节
		if db.AddressType == "DB" {
			// 00.00.00.01 | 00.00.00.02 | 00.00.00.03 | 00.00.00.04
			// 根据类型解析长度
			if err := s1200.client.AGReadDB(db.DataBlockNumber,
				db.ElementNumber, db.DataSize, rData[:]); err != nil {
				glogger.GLogger.Error(err)
				return 0, err
			}
			ValidData := [4]byte{} // 固定4字节，以后有8自己的时候再支持
			copy(ValidData[:], rData[:db.DataSize])
			Value := ParseSiemensSignedValue(db.DataBlockType, db.DataBlockOrder, ValidData)
			// Value := hex.EncodeToString(rData[:db.DataSize])
			values = append(values, __SiemensDataPoint{
				DeviceUUID:      db.DeviceUUID,
				Tag:             db.Tag,
				Value:           Value,
				SiemensAddress:  db.SiemensAddress,
				AddressType:     db.AddressType,
				DataBlockType:   db.DataBlockType,
				DataBlockNumber: db.DataBlockNumber,
				ElementNumber:   db.ElementNumber,
				BitNumber:       db.BitNumber,
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
*解析西门子的值 有符号
*
 */
func ParseSiemensSignedValue(DataBlockType string, DataBlockOrder string, byteSlice [4]byte) string {
	switch DataBlockType {
	case "I", "Q":
		{
			return fmt.Sprintf("%d", byteSlice[0])
		}
	case "BYTE":
		{
			return fmt.Sprintf("%d", byteSlice[0])
		}
	case "SHORT":
		{
			// AB: 1234
			// BA: 3412
			if DataBlockOrder == "AB" {
				uint16Value := uint16(byteSlice[1]) | uint16(byteSlice[0])<<8
				return fmt.Sprintf("%d", uint16Value)

			}
			if DataBlockOrder == "BA" {
				uint16Value := uint16(byteSlice[0]) | uint16(byteSlice[1])<<8
				return fmt.Sprintf("%d", uint16Value)
			}

		}
	case "INT":
		// ABCD
		if DataBlockOrder == "ABCD" {
			intValue := int32(byteSlice[0]) | int32(byteSlice[1])<<8 |
				int32(byteSlice[2])<<16 | int32(byteSlice[3])<<24
			return fmt.Sprintf("%d", intValue)

		}
		if DataBlockOrder == "CDAB" {
			slice := [4]byte{}
			slice[0], slice[1] = byteSlice[2], byteSlice[3]
			slice[2], slice[3] = byteSlice[0], byteSlice[1]
			intValue := int32(slice[0]) | int32(slice[1])<<8 |
				int32(slice[2])<<16 | int32(slice[3])<<24
			return fmt.Sprintf("%d", intValue)
		}
		// 大端字节序转换为int32
		if DataBlockOrder == "DCBA" {
			intValue := int32(byteSlice[3]) | int32(byteSlice[2])<<8 |
				int32(byteSlice[1])<<16 | int32(byteSlice[0])<<24
			return fmt.Sprintf("%d", intValue)
		}
	case "FLOAT": // 3.14159:DCBA -> 40490FDC
		// ABCD
		if DataBlockOrder == "ABCD" {
			intValue := int32(byteSlice[0]) | int32(byteSlice[1])<<8 |
				int32(byteSlice[2])<<16 | int32(byteSlice[3])<<24
			floatValue := float32(math.Float32frombits(uint32(intValue)))
			return fmt.Sprintf("%f", floatValue)
		}
		if DataBlockOrder == "CDAB" {
			intValue := int32(byteSlice[2]) | int32(byteSlice[3])<<8 |
				int32(byteSlice[0])<<16 | int32(byteSlice[1])<<24
			floatValue := float32(math.Float32frombits(uint32(intValue)))
			return fmt.Sprintf("%f", floatValue)
		}
		// 大端字节序转换为int32
		if DataBlockOrder == "DCBA" {
			intValue := int32(byteSlice[3]) | int32(byteSlice[2])<<8 |
				int32(byteSlice[1])<<16 | int32(byteSlice[0])<<24
			floatValue := float32(math.Float32frombits(uint32(intValue)))
			return fmt.Sprintf("%f", floatValue)
		}
	}
	return ""
}

/*
*
*解析西门子的值 无符号
*
 */
func ParseSiemensUSignedValue(DataBlockType string, DataBlockOrder string, byteSlice [4]byte) string {
	return ""
}
