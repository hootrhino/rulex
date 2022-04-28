package test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/goburrow/modbus"
)

var keys = [8]uint16{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}

func TestRTUClientAdvancedUsage(t *testing.T) {
	handler := modbus.NewRTUClientHandler("COM4")
	handler.BaudRate = 9600
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.SlaveId = 1
	handler.Logger = log.New(os.Stdout, "rtu: ", log.LstdFlags)
	err := handler.Connect()
	if err != nil {
		t.Fatal(err)
	}
	defer handler.Close()
	client := modbus.NewClient(handler)
	for _, key := range keys {
		t.Log("------------------> ", key)
		for i := 0; i < 20; i++ {
			client.WriteSingleCoil(key, 0xFF00)
			st1, _ := client.ReadCoils(key, 1)
			t.Log("ReadCoils==> ", st1)
			time.Sleep(500 * time.Microsecond)
			client.WriteSingleCoil(key, 0)
			st2, _ := client.ReadCoils(key, 1)
			t.Log("ReadCoils==> ", st2)
			time.Sleep(500 * time.Microsecond)

		}

	}

}
