package driver

import (
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"

	"github.com/robinson/gos7"
)

/*
*
* 西门子s1200驱动
*
 */
type siemens_s1200_driver struct {
	state      typex.DriverState
	s7client   gos7.Client
	device     *typex.Device
	RuleEngine typex.RuleX
	dbs        []common.S1200Block // PLC 的DB块
}

func NewS1200Driver(d *typex.Device,
	e typex.RuleX,
	s7client gos7.Client,
	dbs []common.S1200Block) typex.XExternalDriver {
	return &siemens_s1200_driver{
		state:      typex.DRIVER_STOP,
		device:     d,
		RuleEngine: e,
		s7client:   s7client,
		dbs:        dbs,
	}
}

func (s1200 *siemens_s1200_driver) Test() error {
	_, err := s1200.s7client.GetCPUInfo()
	if err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	return nil
}

// 读取配置
func (s1200 *siemens_s1200_driver) Init(_ map[string]string) error {
	return nil
}

func (s1200 *siemens_s1200_driver) Work() error {
	return nil
}

func (s1200 *siemens_s1200_driver) State() typex.DriverState {
	_, err := s1200.s7client.PLCGetStatus()
	if err != nil {
		glogger.GLogger.Error(err)
		return typex.DRIVER_STOP
	}
	return typex.DRIVER_UP
}

// 字节格式:[dbNumber1, start1, size1, dbNumber2, start2, size2]
// 读: db --> dbNumber, start, size, buffer[]
var rData = [common.T_2KB]byte{} // 一次最大接受2KB数据

func (s1200 *siemens_s1200_driver) Read(cmd []byte, data []byte) (int, error) {
	values := []common.S1200Block{}
	for _, db := range s1200.dbs {
		//DB 4字节
		if db.Type == "DB" {
			// 00.00.00.01 | 00.00.00.02 | 00.00.00.03 | 00.00.00.04
			if err := s1200.s7client.AGReadDB(db.Address, db.Start, db.Size, rData[:]); err != nil {
				return 0, err
			}
			count := db.Size
			if db.Size*2 > 2000 {
				count = 2000
			}
			values = append(values, common.S1200Block{
				Tag:     db.Tag,
				Address: db.Address,
				Type:    db.Type,
				Start:   db.Start,
				Size:    db.Size,
				Value:   hex.EncodeToString(rData[:count]),
			})
		}
		//
		if db.Type == "MB" {
			// 00.00.00.01 | 00.00.00.02 | 00.00.00.03 | 00.00.00.04
			if err := s1200.s7client.AGReadMB(db.Start, db.Size, rData[:]); err != nil {
				return 0, err
			}
			count := db.Size
			if db.Size*2 > 2000 {
				count = 2000
			}
			values = append(values, common.S1200Block{
				Tag:     db.Tag,
				Type:    db.Type,
				Address: db.Address,
				Start:   db.Start,
				Size:    db.Size,
				Value:   hex.EncodeToString(rData[:count]),
			})
		}
		if db.Frequency < 100 {
			db.Frequency = 100 // 不能太快
		}
		time.Sleep(time.Duration(db.Frequency) * time.Millisecond)
	}
	bytes, _ := json.Marshal(values)
	copy(data, bytes)
	return len(bytes), nil
}

// db.Address:int, db.Start:int, db.Size:int, rData[]
// [
//
//	{
//	    "tag":"V",
//	    "address":1,
//	    "start":1,
//	    "size":1,
//	    "value":"AAECAwQ="
//	}
//
// ]
func (s1200 *siemens_s1200_driver) Write(cmd []byte, data []byte) (int, error) {
	blocks := []common.S1200Block{}
	if err := json.Unmarshal(data, &blocks); err != nil {
		return 0, err
	}
	//
	for _, block := range blocks {
		hexV, _ := hex.DecodeString(block.Value)
		if err := s1200.s7client.AGWriteDB(
			block.Address,
			block.Start,
			block.Size,
			hexV,
		); err != nil {
			return 0, err
		}
	}
	return 0, nil
}

func (s1200 *siemens_s1200_driver) DriverDetail() typex.DriverDetail {
	return typex.DriverDetail{
		Name:        "SIEMENS_s1200",
		Type:        "TCP",
		Description: "SIEMENS s1200 系列 PLC 驱动",
	}
}

func (s1200 *siemens_s1200_driver) Stop() error {
	return nil
}
