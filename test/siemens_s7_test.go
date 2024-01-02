package test

import (
	"encoding/json"
	"testing"
	"time"

	profinet "github.com/Kowiste/ProfinetServer"
	"github.com/hootrhino/rulex/utils"
	"github.com/robinson/gos7"
)

func Test_server(t *testing.T) {
	server := profinet.NewServer()
	server.SetDB(10, []uint16{11, 22, 33, 44, 55, 66, 77, 88, 99, 100})
	err := server.Listen("0.0.0.0:1800", 0, 1)
	if err != nil {
		t.Error(err)
		return
	}
	client(t)
	time.Sleep(2 * time.Second)
}

func client(t *testing.T) {
	handler := gos7.NewTCPClientHandler("127.0.0.1:1800", 0, 1)

	defer handler.Close()
	if err := handler.Connect(); err != nil {
		t.Error(err)
		return
	}
	client := gos7.NewClient(handler)
	dataBuf := make([]byte, 10)
	if err := client.AGReadDB(10, 0, 10, dataBuf); err != nil {
		t.Error(err)
		return
	}
	t.Log("client.AGReadDB =>", dataBuf)

}

/*
*
*  Error: Address out of range
*
 */
func Test_readDB(t *testing.T) {
	handler := gos7.NewTCPClientHandler("127.0.0.1:1800", 0, 1)

	defer handler.Close()
	if err := handler.Connect(); err != nil {
		t.Error(err)
		return
	}
	client := gos7.NewClient(handler)
	dataBuf := make([]byte, 10)
	if err := client.AGReadDB(10, 0, 10, dataBuf); err != nil {
		t.Error(err)
		return
	}
	t.Log("client.AGReadDB =>", dataBuf)

}
func Test_gen_db_json(t *testing.T) {

	type S1200BlockValue struct {
		Tag     string `json:"tag"`     // 数据tag
		Address int    `json:"address"` // 地址
		Start   int    `json:"start"`   // 起始地址
		Size    int    `json:"size"`    // 数据长度
		Value   []byte `json:"value"`
	}
	blocks := []S1200BlockValue{{
		Tag:     "V",
		Address: 1,
		Start:   1,
		Size:    1,
		Value:   []byte{0, 1, 2, 3, 4},
	}}
	bytes, _ := json.Marshal(blocks)
	t.Log(string(bytes))
	blocks2 := []S1200BlockValue{}
	json.Unmarshal(bytes, &blocks2)
	t.Log(blocks2[0].Value)
}

// go test -timeout 30s -run ^Test_gen_point github.com/hootrhino/rulex/test -v -count=1
func Test_gen_point(t *testing.T) {
	type SiemensDataPoint struct {
		UUID           string `json:"uuid"`
		DeviceUUID     string `json:"device_uuid"`
		SiemensAddress string `json:"siemensAddress"` // 西门子的地址字符串
		Tag            string `json:"tag"`
		Alias          string `json:"alias"`
		DataOrder      string `json:"dataOrder"` // 字节序
		DataType       string `json:"dataType"`
		Frequency      *int64 `json:"frequency"`
		Status         int    `json:"status"`        // 运行时数据
		LastFetchTime  uint64 `json:"lastFetchTime"` // 运行时数据
		Value          string `json:"value"`         // 运行时数据
		// 西门子解析后的地址信息
		AddressType     string `json:"addressType"`     // 寄存器类型: DB I Q
		DataBlockType   string `json:"dataBlockType"`   // 数据类型: D X
		DataBlockNumber *int   `json:"dataBlockNumber"` // 数据块号: 100...
		ElementNumber   *int   `json:"elementNumber"`   // 元素号:1000...
		BitNumber       *int   `json:"bitNumber"`       // 位号，只针对I、Q
	}
	P := SiemensDataPoint{
		UUID:       utils.MakeUUID("S7"),
		DeviceUUID: "123",
		Tag:        "温度",
		Frequency: func() *int64 {
			a := int64(100)
			return &a
		}(),
		SiemensAddress: "DB4900.DBD1000",
		AddressType:    "DB",
		DataBlockType:  "DB",
		DataBlockNumber: func() *int {
			a := int(4900)
			return &a
		}(),
		ElementNumber: func() *int {
			a := int(1000)
			return &a
		}(),
		BitNumber: func() *int {
			a := int(0)
			return &a
		}(),
		DataOrder: "ABCD",
	}
	if bytes, err := json.MarshalIndent(P, "", "    "); err != nil {
		panic(err)
	} else {
		t.Log(string(bytes))
	}

}
