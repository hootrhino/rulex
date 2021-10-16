package test

import (
	"encoding/json"
	"testing"

	"github.com/go-playground/validator/v10"
)

type ModBusConfig struct {
	Mode           string          `json:"mode"`
	Timeout        int             `json:"timeout" validate:"required,gte=1,lte=60"`
	SlaverId       byte            `json:"slaverId" validate:"required,gte=1,lte=255"`
	Frequency      int64           `json:"frequency" validate:"required,gte=1,lte=10000"`
	RtuConfig      RtuConfig       `json:"rtuConfig" validate:"required"`
	TcpConfig      TcpConfig       `json:"tcpConfig" validate:"required"`
	RegisterParams []RegisterParam `json:"registerParams" validate:"required"`
}

type RegisterParam struct {
	// 	Code |  Register Type
	//-------------------------------------------
	// 	1		Read Coil
	// 	2		Read Discrete Input
	// 	3		Read Holding Registers
	// 	4		Read Input Registers
	// 	5		Write Single Coil
	// 	6		Write Single Holding Register
	// 	15		Write Multiple Coils
	// 	16		Write Multiple Holding Registers
	Function int    `json:"function" validate:"1|2|3|4|"`               // Current version only support read
	Address  uint16 `json:"address" validate:"required,gte=0,lte=255"`  // Address
	Quantity uint16 `json:"quantity" validate:"required,gte=0,lte=255"` // Quantity
}

//
// Uart "/dev/ttyUSB0"
// BaudRate = 115200
// DataBits = 8
// Parity = "N"
// StopBits = 1
// SlaveId = 1
// Timeout = 5 * time.Second
//
type RtuConfig struct {
	Uart     string `json:"uart" validate:"required"`
	BaudRate int    `json:"baudRate" validate:"required"`
}

//
//
//
type TcpConfig struct {
	Ip   string `json:"ip" validate:"required"`
	Port int    `json:"port" validate:"required,gte=1,lte=65535"`
}

func Test_validator(t *testing.T) {
	mainConfig := ModBusConfig{
		Frequency: 3, // second
		Mode:      "TCP",
		Timeout:   10, // second
		SlaverId:  1,
		TcpConfig: TcpConfig{
			Ip:   "127.0.0.1",
			Port: 502,
		},
		RtuConfig: RtuConfig{
			Uart:     "TCP",
			BaudRate: 115200,
		},

		RegisterParams: []RegisterParam{{1, 1, 10}},
	}
	b, _ := json.Marshal(mainConfig)
	// {"ip":"127.0.0.1","port":502,"timeout":10,"slaverId":1,"registerParams":[{"function":1,"address":1,"quantity":10}]}
	t.Log(string(b))
	if err := validator.New().Struct(mainConfig); err != nil {
		t.Error(err)
	}
}
