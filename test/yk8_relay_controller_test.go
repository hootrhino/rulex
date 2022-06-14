package test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/goburrow/modbus"
)

/*
*
*  继电器测试
*
 */
func TestRTU_YK08(t *testing.T) {
	handler := modbus.NewRTUClientHandler("COM6")
	handler.BaudRate = 9600
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.SlaveId = 1
	handler.Logger = log.New(os.Stdout, "rtu: ", log.LstdFlags)

	if err := handler.Connect(); err != nil {
		t.Fatal(err)
	}
	defer handler.Close()
	client := modbus.NewClient(handler)
	client.WriteMultipleCoils(0, 1, []byte{0b00000001})
	time.Sleep(1 * time.Second)

	if results, err := client.ReadCoils(0x00, 1); err != nil {
		t.Fatal(err)
	} else {
		t.Log("===> ", results)
	}

	client.WriteMultipleCoils(0, 1, []byte{0b00000011})
	time.Sleep(1 * time.Second)
	if results, err := client.ReadCoils(0x00, 1); err != nil {
		t.Fatal(err)
	} else {
		t.Log("===> ", results)
	}

	client.WriteMultipleCoils(0, 1, []byte{0b00000111})
	time.Sleep(1 * time.Second)
	if results, err := client.ReadCoils(0x00, 1); err != nil {
		t.Fatal(err)
	} else {
		t.Log("===> ", results)
	}

	client.WriteMultipleCoils(0, 1, []byte{0b00001111})
	time.Sleep(1 * time.Second)
	//
	if results, err := client.ReadCoils(0x00, 1); err != nil {
		t.Fatal(err)
	} else {
		t.Log("===> ", results)
	}

	client.WriteMultipleCoils(0, 1, []byte{0b00011111})
	time.Sleep(1 * time.Second)
	//
	if results, err := client.ReadCoils(0x00, 1); err != nil {
		t.Fatal(err)
	} else {
		t.Log("===> ", results)
	}

	client.WriteMultipleCoils(0, 1, []byte{0b00111111})
	time.Sleep(1 * time.Second)
	//
	if results, err := client.ReadCoils(0x00, 1); err != nil {
		t.Fatal(err)
	} else {
		t.Log("===> ", results)
	}

	client.WriteMultipleCoils(0, 1, []byte{0b01111111})
	time.Sleep(1 * time.Second)
	//
	if results, err := client.ReadCoils(0x00, 1); err != nil {
		t.Fatal(err)
	} else {
		t.Log("===> ", results)
	}

	client.WriteMultipleCoils(0, 1, []byte{0b11111111})
	time.Sleep(1 * time.Second)
	//
	if results, err := client.ReadCoils(0x00, 1); err != nil {
		t.Fatal(err)
	} else {
		t.Log("===> ", results)
	}

}

func TestRTU_YK081(t *testing.T) {
	handler := modbus.NewRTUClientHandler("COM6")
	handler.BaudRate = 9600
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.SlaveId = 1
	handler.Logger = log.New(os.Stdout, "rtu: ", log.LstdFlags)
	if err := handler.Connect(); err != nil {
		t.Error(err)
		return
	}
	defer handler.Close()
	client := modbus.NewClient(handler)

	if results, err := client.ReadCoils(0x00, 0x08); err != nil {
		t.Error(err)
		return
	} else {
		t.Log("===> ", results)
	}
}
