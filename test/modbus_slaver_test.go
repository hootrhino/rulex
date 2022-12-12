package test

import (
	"log"
	"os"
	"testing"
	"time"

	modbus "github.com/wwhai/gomodbus"
)

func TestTCPClient(t *testing.T) {
	handler := modbus.NewTCPClientHandler("127.0.0.1:502")
	handler.Timeout = 10 * time.Second
	handler.SlaveId = 0xFF
	handler.Logger = log.New(os.Stdout, "test: ", log.LstdFlags)
	if err := handler.Connect(); err != nil {
		t.Fatal(err)
	}
	defer handler.Close()

	client := modbus.NewClient(handler)
	if results, err := client.ReadCoils(0, 10); err != nil {
		t.Fatal(err)
	} else {
		t.Log(results)
	}
	if results, err := client.ReadInputRegisters(1, 10); err != nil {
		t.Fatal(err)
	} else {
		t.Log(results)
	}
	if results, err := client.ReadHoldingRegisters(2, 10); err != nil {
		t.Fatal(err)
	} else {
		t.Log(results)
	}
	if results, err := client.ReadDiscreteInputs(3, 10); err != nil {
		t.Fatal(err)
	} else {
		t.Log(results)
	}
	time.Sleep(1 * time.Second)
}
