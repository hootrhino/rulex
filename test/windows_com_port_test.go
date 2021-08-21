package test

import (
	"fmt"
	"testing"

	"github.com/ngaut/log"
	"go.bug.st/serial"
)

func TestComPort(t *testing.T) {

	mode := &serial.Mode{
		BaudRate: 115200,
	}
	port, err := serial.Open("COM6", mode)
	if err != nil {
		log.Fatal(err)
	}
	if err := port.SetMode(mode); err != nil {
		log.Fatal(err)
	}
	buff := make([]byte, 100)
	for {
		n, err := port.Read(buff)
		if err != nil {
			log.Fatal(err)
			break
		}
		if n == 0 {
			fmt.Println("\nEOF")
			break
		}
		fmt.Printf("%v", string(buff[:n]))
	}
}
