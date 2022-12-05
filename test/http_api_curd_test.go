/*
*
* 资源增删改查单元测试
*
 */
package test

import (
	"os/exec"
	"testing"

	"github.com/go-playground/assert/v2"
	httpserver "github.com/i4de/rulex/plugin/http_server"
)

/*
*
* 对设备的增删改查
*
 */

// go test -timeout 30s -run ^Test_CURD_Device github.com/i4de/rulex/test -v -count=1
func Test_CURD_Device(t *testing.T) {
	RmUnitTestDbFile(t)
	StartTestServer(t)
	UT_createDevice(t)
}
func UT_createDevice(t *testing.T) {
	// 通过接口创建一个设备
	body := `{"name":"GENERIC_SNMP","type":"GENERIC_SNMP","config":{"timeout":5,"frequency":5,"target":"127.0.0.1","port":161,"transport":"udp","community":"public","version":3,"dataModels":[]},"description":"GENERIC_SNMP"}`
	output, err := exec.Command("curl",
		"-X", "POST", "http://127.0.0.1:2580/api/v1/devices",
		"-H", "'Content-Type: application/json'", "-d", body).Output()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Test_CURD_Device: ", string(output))
	//
	LoadDB()
	mdevice := []httpserver.MDevice{}
	unitTestDB.Raw("SELECT * FROM m_devices").Find(&mdevice)
	assert.Equal(t, 1, len(mdevice))
	assert.Equal(t, mdevice[0].Name, "GENERIC_SNMP")
	assert.Equal(t, mdevice[0].Description, "GENERIC_SNMP")
	assert.Equal(t, mdevice[0].Type, "GENERIC_SNMP")
}
