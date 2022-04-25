package test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

var (
	// Use mcu pin 10, corresponds to physical pin 19 on the pi
	pin = rpio.Pin(10)
)

func Test_GPIO_BLINKER(t *testing.T) {
	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Unmap gpio memory when done
	defer rpio.Close()

	// Set pin to output mode
	pin.Output()

	// Toggle pin 20 times
	for x := 0; x < 20; x++ {
		pin.Toggle()
		time.Sleep(time.Second / 5)
	}
}
func Test_GPIO_SPI(t *testing.T) {
	if err := rpio.Open(); err != nil {
		panic(err)
	}

	if err := rpio.SpiBegin(rpio.Spi0); err != nil {
		panic(err)
	}

	rpio.SpiChipSelect(0) // Select CE0 slave

	// Send
	rpio.SpiTransmit(0xFF)             // send single byte
	rpio.SpiTransmit(0xDE, 0xAD, 0xBE) // send several bytes

	data := []byte{'H', 'e', 'l', 'l', 'o', 0}
	rpio.SpiTransmit(data...) // send slice of bytes

	// Receive

	received := rpio.SpiReceive(5) // receive 5 bytes, (sends 5 x 0s)
	fmt.Println(received)

	// Send & Receive

	buffer := []byte{0xDE, 0xED, 0xBE, 0xEF}
	rpio.SpiExchange(buffer) // buffer is populated with received data
	fmt.Println(buffer)

	rpio.SpiEnd(rpio.Spi0)
	rpio.Close()
}
