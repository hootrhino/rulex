package test

import (
	"testing"
	"time"

	profinet "github.com/Kowiste/ProfinetServer"
	"github.com/robinson/gos7"
)

func Test_server(t *testing.T) {

	server := profinet.NewServer()
	server.SetOutput([]uint16{11, 22, 33, 44, 55, 66, 77, 88, 99, 100})
	server.SetInput([]uint16{11, 22, 33, 44, 55, 66, 77, 88, 99, 100})
	server.SetDB(10, []uint16{11, 22, 33, 44, 55, 66, 77, 88, 99, 100})
	err := server.Listen("0.0.0.0:1503", 0, 1)
	if err != nil {
		t.Error(err)
		return
	}
	client(t)
	time.Sleep(2 * time.Second)
}
func client(t *testing.T) {
	handler := gos7.NewTCPClientHandler("127.0.0.1:1503", 0, 1)
	err1 := handler.Connect()
	defer handler.Close()
	if err1 != nil {
		t.Error(err1)
		return
	}
	client := gos7.NewClient(handler)
	// buf := make([]byte, 2)
	// buf[0] = 24
	// buf[1] = 34
	// println("SEND ", binary.BigEndian.Uint16(buf))
	// if err := client.AGWriteDB(13, 4, 2, buf); err != nil {
	// 	t.Error(err)
	// }
	buf2 := make([]byte, 10)
	if err := client.AGReadDB(10, 0, 10, buf2); err != nil {
		t.Error(err)
	}
	t.Log("client.AGReadDB =>", buf2)

}
