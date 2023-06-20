package test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/hootrhino/rulex/glogger"
	modbus "github.com/wwhai/gomodbus"
)

var keys = [8]uint16{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}

/*
*
*  继电器测试
*
 */
func TestRTU_relay(t *testing.T) {
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
			time.Sleep(500 * time.Millisecond)
			client.WriteSingleCoil(key, 0)
			st2, _ := client.ReadCoils(key, 1)
			t.Log("ReadCoils==> ", st2)
			time.Sleep(500 * time.Millisecond)

		}

	}

}

/*
*
* 温湿度传感器测试
*
 */
func TestRTU485_THer_Usage(t *testing.T) {
	handler := modbus.NewRTUClientHandler("COM6")
	handler.BaudRate = 4800
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
	results, err := client.ReadHoldingRegisters(0x00, 2)
	if err != nil {
		glogger.GLogger.Fatal(err)
	}
	glogger.GLogger.Println(results)

}

/*
*
* 空调面板测试
*
 */
func TestRTU485_Air_Usage(t *testing.T) {
	handler := modbus.NewRTUClientHandler("COM4")
	handler.BaudRate = 4800
	handler.DataBits = 8
	handler.Parity = "O"
	handler.StopBits = 1
	handler.SlaveId = 1
	handler.Logger = log.New(os.Stdout, "rtu: ", log.LstdFlags)

	err := handler.Connect()
	if err != nil {
		t.Fatal(err)
	}
	defer handler.Close()
	client := modbus.NewClient(handler)
	results, err := client.ReadHoldingRegisters(0x09, 1)
	glogger.GLogger.Println(results)
	if err != nil {
		glogger.GLogger.Fatal(err)
	}
	glogger.GLogger.Println(results)

}
