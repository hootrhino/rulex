package test

import (
	"log"
	"testing"
	"time"

	goserial "github.com/wwhai/goserial"
)

func Test_goserial(t *testing.T) {
	port, err := goserial.Open(&goserial.Config{
		Address:  "COM10",
		BaudRate: 4800,
		DataBits: 8,
		Parity:   "O",
		StopBits: 1,
		Timeout:  30 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()
	_, err = port.Write([]byte{0x01, 0x03, 0x00, 0x00, 0x00, 0x01, 0x84, 0x0A})
	if err != nil {
		log.Fatal(err)
	}
	bytes := []byte{}
	time.Sleep(1 * time.Second)
	n, err := port.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}
	t.Log(n, bytes)
}
