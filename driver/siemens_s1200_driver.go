package driver

import (
	"encoding/json"
	"sync"

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
	lock       sync.Mutex
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
		lock:       sync.Mutex{},
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
func (s1200 *siemens_s1200_driver) Read(cmd []byte, data []byte) (int, error) {
	values := []common.S1200BlockValue{}
	for _, db := range s1200.dbs {
		rData := []byte{}
		s1200.lock.Lock()
		if err := s1200.s7client.AGReadDB(db.Address, db.Start, db.Size, rData); err != nil {
			s1200.lock.Unlock()
			return 0, err
		}
		s1200.lock.Unlock()
		values = append(values, common.S1200BlockValue{
			Tag:     db.Tag,
			Address: db.Address,
			Start:   db.Start,
			Size:    db.Size,
			Value:   rData,
		})

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
	blocks := []common.S1200BlockValue{}
	if err := json.Unmarshal(data, &blocks); err != nil {
		return 0, err
	}
	//
	for _, block := range blocks {
		s1200.lock.Lock()
		if err := s1200.s7client.AGWriteDB(
			block.Address,
			block.Start,
			block.Size,
			block.Value,
		); err != nil {
			s1200.lock.Unlock()
			return 0, err
		}
		s1200.lock.Unlock()
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
