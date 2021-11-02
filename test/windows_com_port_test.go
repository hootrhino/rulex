package test

import (
	"testing"

	"github.com/goburrow/serial"
	"github.com/ngaut/log"
)

func TestComPort(t *testing.T) {
	port, err := serial.Open(&serial.Config{Address: "COM4", BaudRate: 115200})
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()

	buffer := make([]byte, 1)
	_, err = port.Read(buffer)
	for {
		if err != nil {
			log.Fatal(err)
		} else {
			print((string(buffer)))
		}
	}
}
