package test

import (
	"runtime"
	"testing"

	"github.com/hootrhino/rulex/ossupport"
	"go.bug.st/serial"
)

func Test_GetPortsList(t *testing.T) {
	t.Log(GetUartList())
}
func GetUartList() []string {
	var ports []string
	if runtime.GOOS == "windows" {
		ports, _ = serial.GetPortsList()
	} else {
		ports, _ = ossupport.GetPortsListUnix()
	}
	return ports
}
