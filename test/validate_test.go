package test

import (
	"encoding/json"
	"testing"

	"github.com/go-playground/validator/v10"
)

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
	Function int    `json:"function" validate:"1|2|3|4|"`      // Current version only support read
	Address  uint16 `json:"address" validate:"gte=0,lte=255"`  // Address
	Quantity uint16 `json:"quantity" validate:"gte=0,lte=255"` // Quantity
}
type ModBusConfig struct {
	Ip             string          `json:"ip" validate:"required"`
	Port           int             `json:"port" validate:"gte=1024,lte=65535"`
	Timeout        int             `json:"timeout" validate:"required"`
	SlaverId       byte            `json:"slaverId" validate:"gte=1,lte=255"`
	RegisterParams []RegisterParam `json:"registerParams" validate:"required"`
}

func Test_validator(t *testing.T) {
	mainConfig := ModBusConfig{
		Ip:             "127.0.0.1",
		Port:           502,
		Timeout:        10, // second
		SlaverId:       1,
		RegisterParams: []RegisterParam{{1, 1, 10}},
	}
	b, _ := json.Marshal(mainConfig)
	// {"ip":"127.0.0.1","port":502,"timeout":10,"slaverId":1,"registerParams":[{"function":1,"address":1,"quantity":10}]}
	t.Log(string(b))
	if err := validator.New().Struct(mainConfig); err != nil {
		t.Error(err)
	}
}
