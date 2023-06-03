package test

import (
	"testing"

	"github.com/hootrhino/rulex/utils"
)

// go test -timeout 30s -run ^TestOk github.com/hootrhino/rulex/test -v -count=1
func Test_Win32_GetCpuUsage(t *testing.T) {
	t.Log(utils.GetCpuUsage())
}
func Test_Win32_GetDiskUsage(t *testing.T) {
	t.Log(utils.GetDiskUsage())
}
func Test_Win32_NetInterfaceUsage(t *testing.T) {
	t.Log(utils.NetInterfaceUsage())
}
