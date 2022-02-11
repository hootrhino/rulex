package test

import (
	"testing"

	"github.com/ngaut/log"
	"go.bug.st/serial"
)

func Test_GetPortsList(t *testing.T) {
	t.Log(GetUartList())
}
func GetUartList() []string {
	r := []string{}
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Error(err)
		return r
	}
	r = append(r, ports...)
	return r
}
