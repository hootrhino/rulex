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
