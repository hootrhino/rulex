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
 * 对App的增删改查
 *
 */

// go test -timeout 30s -run ^Test_CURD_App github.com/hootrhino/rulex/test -v -count=1
func Test_CURD_App(t *testing.T) {
	RmUnitTestDbFile(t)
	StartTestServer(t)
	uuid := UT_createApp(t)
	UT_updateApp(t, uuid)
	UT_deleteApp(t, uuid)
}
func UT_createApp(t *testing.T) string {
	// 通过接口创建一个App
	body := `{"name": "testlua1","version": "1.0.0","autoStart": true,"description": "hello world"}`
	output, err := exec.Command("curl",
		"-X", "POST", "http://127.0.0.1:2580/api/v1/app",
		"-H", "Content-Type: application/json", "-d", body).Output()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("UT_createApp: ", string(output))
	//
	LoadUnitTestDB()
	mApp := []model.MApp{}
	unitTestDB.Raw("SELECT * FROM m_apps").Find(&mApp)
	assert.Equal(t, 1, len(mApp))
	t.Log(mApp[0].UUID)
	assert.Equal(t, mApp[0].Name, "testlua1")
	assert.Equal(t, mApp[0].Version, "1.0.0")
	assert.Equal(t, mApp[0].AutoStart, true)
	return mApp[0].UUID
}

/*
*
*
*
 */
func UT_updateApp(t *testing.T, uuid string) {
	body := `{"uuid": "%s","name": "testlua11","version": "1.0.1","autoStart": false,"luaSource":"AppNAME='OK1'\nAppVERSION='0.0.3'\nfunction Main()\n\tprint('Hello-World')\nend"}`
	t.Logf(body, uuid)
	output, err := exec.Command("curl",
		"-X", "PUT", "http://127.0.0.1:2580/api/v1/app",
		"-H", "Content-Type: application/json", "-d", fmt.Sprintf(body, uuid)).Output()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("UT_updateApp: ", string(output))
	LoadUnitTestDB()
	mApp := []model.MApp{}
	unitTestDB.Raw("SELECT * FROM m_apps").Find(&mApp)
	assert.Equal(t, 1, len(mApp))
	t.Log(mApp[0].UUID)
	assert.Equal(t, mApp[0].Name, "testlua11")
	assert.Equal(t, mApp[0].Version, "1.0.1")
	assert.Equal(t, mApp[0].AutoStart, false)
}
func UT_deleteApp(t *testing.T, uuid string) {
	// 删除一个App
	output, err := exec.Command("curl", "-X", "DELETE", "http://127.0.0.1:2580/api/v1/app?uuid="+uuid).Output()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("UT_deleteApp: ", string(output))
	//
	LoadUnitTestDB()
	mApp := []model.MApp{}
	unitTestDB.Raw("SELECT * FROM m_apps").Find(&mApp)
	assert.Equal(t, 0, len(mApp))
}
