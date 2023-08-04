/*
*
* 资源增删改查单元测试
*
 */
package test

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/hootrhino/rulex/plugin/http_server/model"
)

/*
*
* 对设备的增删改查
*
 */

// go test -timeout 30s -run ^Test_CURD_Device github.com/hootrhino/rulex/test -v -count=1
func Test_CURD_Device(t *testing.T) {
	RmUnitTestDbFile(t)
	StartTestServer(t)
	uuid := UT_createDevice(t)
	UT_updateDevice(t, uuid)
	UT_deleteDevice(t, uuid)
}
func UT_createDevice(t *testing.T) string {
	// 通过接口创建一个设备
	body := `{"name":"GENERIC_SNMP","type":"GENERIC_SNMP","config":{"timeout":5,"frequency":5,"target":"127.0.0.1","port":161,"transport":"udp","community":"public","version":3,"dataModels":[]},"description":"GENERIC_SNMP"}`
	output, err := exec.Command("curl",
		"-X", "POST", "http://127.0.0.1:2580/api/v1/devices",
		"-H", "'Content-Type: application/json'", "-d", body).Output()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("UT_createDevice: ", string(output))
	//
	LoadUnitTestDB()
	mDevice := []model.MDevice{}
	unitTestDB.Raw("SELECT * FROM m_devices").Find(&mDevice)
	assert.Equal(t, 1, len(mDevice))
	t.Log(mDevice[0].UUID)
	assert.Equal(t, mDevice[0].Name, "GENERIC_SNMP")
	assert.Equal(t, mDevice[0].Description, "GENERIC_SNMP")
	assert.Equal(t, mDevice[0].Type, "GENERIC_SNMP")
	return mDevice[0].UUID
}
func UT_updateDevice(t *testing.T, uuid string) {
	// 通过接口创建一个设备
	body := `{"uuid":"%s","name":"GENERIC_SNMP_NEW","type":"GENERIC_SNMP","config":{"timeout":10,"frequency":10,"target":"127.0.0.2","port":161,"transport":"udp","community":"public","version":3,"dataModels":[]},"description":"GENERIC_SNMP_NEW"}`
	output, err := exec.Command("curl",
		"-X", "PUT", "http://127.0.0.1:2580/api/v1/devices",
		"-H", "'Content-Type: application/json'", "-d", fmt.Sprintf(body, uuid)).Output()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("UT_updateDevice: ", string(output))
	LoadUnitTestDB()
	mDevice := []model.MDevice{}
	unitTestDB.Raw("SELECT * FROM m_devices").Find(&mDevice)
	assert.Equal(t, 1, len(mDevice))
	t.Log(mDevice[0].UUID)
	assert.Equal(t, mDevice[0].Name, "GENERIC_SNMP_NEW")
	assert.Equal(t, mDevice[0].Description, "GENERIC_SNMP_NEW")
	assert.Equal(t, mDevice[0].Type, "GENERIC_SNMP")
}
func UT_deleteDevice(t *testing.T, uuid string) {
	// 删除一个设备
	output, err := exec.Command("curl", "-X", "DELETE", "http://127.0.0.1:2580/api/v1/devices?uuid="+uuid).Output()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("UT_deleteDevice: ", string(output))
	//
	LoadUnitTestDB()
	mDevice := []model.MDevice{}
	unitTestDB.Raw("SELECT * FROM m_devices").Find(&mDevice)
	assert.Equal(t, 0, len(mDevice))
}
