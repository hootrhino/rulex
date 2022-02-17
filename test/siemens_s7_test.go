package test

import (
	"encoding/json"
	"testing"
	"time"

	profinet "github.com/Kowiste/ProfinetServer"
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
func Test_gen_config(t *testing.T) {
	type stateAddress struct {
		Address int `json:"address"` // 地址
		Start   int `json:"start"`   // 起始地址
		Size    int `json:"size"`    // 数据长度
	}
	type db struct {
		Tag     string `json:"tag"`     // 数据tag
		Address int    `json:"address"` // 地址
		Start   int    `json:"start"`   // 起始地址
		Size    int    `json:"size"`    // 数据长度
	}
	type siemensS7config struct {
		Host         string       `json:"host" validate:"required" title:"IP地址" info:""`          // 127.0.0.1
		Rack         int          `json:"rack" validate:"required" title:"架号" info:""`            // 0
		Slot         int          `json:"slot" validate:"required" title:"槽号" info:""`            // 1
		Timeout      int          `json:"timeout" validate:"required" title:"连接超时时间" info:""`     // 5s
		IdleTimeout  int          `json:"idleTimeout" validate:"required" title:"心跳超时时间" info:""` // 5s
		Frequency    int64        `json:"frequency" validate:"required" title:"采集频率" info:""`     // 5s
		StateAddress stateAddress `json:"stateAddress" validate:"required" title:"状态地址" info:""`  // 5s
		Dbs          []db         `json:"dbs" validate:"required" title:"采集配置" info:""`           // Db
	}
	c := siemensS7config{
		Host:        "",
		Rack:        0,
		Slot:        1,
		Timeout:     5,
		IdleTimeout: 5,
		Frequency:   5,
		Dbs: []db{
			{
				Tag:     "Votage",
				Address: 0,
				Start:   1,
				Size:    1,
			},
		},
	}
	b, _ := json.MarshalIndent(c, "", " ")
	t.Log(string(b))

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
