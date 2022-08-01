package test

import (
	"testing"

	"github.com/i4de/rulex/device"
	"github.com/i4de/rulex/engine"
)

/*
*
* 生成标准lua库的文档
*
 */
func Test_Gen_rulexlib_doc(t *testing.T) {
	engine.BuildInLuaLibDoc()
}

/*
*
* 生成设备文档
*
 */
func Test_Gen_device_doc(t *testing.T) {
	device.BuildDoc()
}
