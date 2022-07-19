package driver

import (
	"encoding/binary"
	"encoding/json"

	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"

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
	dbs        []S1200Block // PLC 的DB块
}

func NewS1200Driver(d *typex.Device,
	e typex.RuleX,
	s7client gos7.Client,
	dbs []S1200Block) typex.XExternalDriver {
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
	_, err := s1200.s7client.GetCPUInfo()
	if err != nil {
		glogger.GLogger.Error(err)
		return typex.DRIVER_STOP
	}
	return typex.DRIVER_RUNNING
}

//
// 字节格式:[dbNumber1, start1, size1, dbNumber2, start2, size2]
// 读: db --> dbNumber, start, size, buffer[]
//
func (s1200 *siemens_s1200_driver) Read(data []byte) (int, error) {
	values := []S1200BlockValue{}
	for _, db := range s1200.dbs {
		rData := []byte{}
		if err := s1200.s7client.AGReadDB(db.Address, db.Start, db.Size, rData); err != nil {
			return 0, err
		}
		values = append(values, S1200BlockValue{
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

//
// db.Address:int, db.Start:int, db.Size:int, rData[]
// Example: 给地址为1的DB 起始1 写入4个字节: 1 2 3 4
//          0x00 0x00 0x00 0x01 | 0x00 0x00 0x00 0x01 | 0x00 0x00 0x00 0x04 | 0x04 0x03 0x02 0x01
// data := []byte{
// 	0x00, 0x00, 0x00, 0x01,
// 	0x00, 0x00, 0x00, 0x01,
// 	0x00, 0x00, 0x00, 0x04,
// 	0x04, 0x03, 0x02, 0x01,
// }
//
//
func (s1200 *siemens_s1200_driver) Write(data []byte) (int, error) {
	Address := binary.BigEndian.Uint32(data[0:4])
	Start := binary.BigEndian.Uint32(data[4:8])
	Size := binary.BigEndian.Uint32(data[8:12])
	//
	if err := s1200.s7client.AGWriteDB(int(Address), int(Start), int(Size), data[12:Size]); err != nil {
		return 0, err
	} else {
		return 0, nil
	}
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
