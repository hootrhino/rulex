package httpserver

import (
	"fmt"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/i4de/rulex/appstack"
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"
)

/*
*
* 其实这个结构体扮演的角色VO层
*
 */
type web_data_app struct {
	UUID      string `json:"uuid"`      // 名称
	Name      string `json:"name"`      // 名称
	Version   string `json:"version"`   // 版本号
	AutoStart bool   `json:"autoStart"` // 自动启动
	AppState  int    `json:"appState"`  // 状态: 1 运行中, 0 停止
	Filepath  string `json:"filepath"`  // 文件路径, 是相对于main的apps目录
	LuaSource string `json:"luaSource"`
}

// 列表
func Apps(c *gin.Context, hs *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	// 从配置拿App
	if uuid == "" {
		result := []web_data_app{}
		for _, app := range hs.AllApp() {
			web_data := web_data_app{
				UUID:      app.UUID,
				Name:      app.Name,
				Version:   app.Version,
				AutoStart: *app.AutoStart,
				AppState: func() int {
					if a := e.GetApp(app.UUID); a != nil {
						return int(a.AppState)
					}
					return 0
				}(),
				Filepath: app.Filepath,
			}
			result = append(result, web_data)
		}
		c.JSON(200, OkWithData(result))
		return
	}
	// uuid
	appInfo, err1 := hs.GetAppWithUUID(uuid)
	if err1 != nil {
		c.JSON(200, Error400(err1))
		return
	}
	web_data := web_data_app{
		UUID:      appInfo.UUID,
		Name:      appInfo.Name,
		Version:   appInfo.Version,
		AutoStart: *appInfo.AutoStart,
		AppState: func() int {
			if a := e.GetApp(appInfo.UUID); a != nil {
				return int(a.AppState)
			}
			return 0
		}(),
		Filepath: appInfo.Filepath,
		LuaSource: func() string {
			path := "./apps/" + appInfo.UUID + ".lua"
			bytes, err := os.ReadFile(path)
			if err != nil {
				return err.Error()
			}
			return string(bytes)
		}(),
	}
	c.JSON(200, OkWithData(web_data))

}

/*
*
* 直接新建一个文件，文件名为 UUID.lua
*
 */
const semVerRegexExpr = `^(0|[1-9]+[0-9]*)\.(0|[1-9]+[0-9]*)\.(0|[1-9]+[0-9]*)(-(0|[1-9A-Za-z-][0-9A-Za-z-]*)(\.[0-9A-Za-z-]+)*)?(\+[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?$`
const luaTemplate = `
AppNAME = "%s"
AppVERSION = "%s"
AppDESCRIPTION = "%s"
--
-- Main
--
%s
`
const defaultLuaMain = `
function Main(arg)
	print("Hello World:",  applib:Time())
	return 0
end
`

func CreateApp(c *gin.Context, hs *HttpApiServer, e typex.RuleX) {
	type Form struct {
		Name        string `json:"name"`        // 名称
		Version     string `json:"version"`     // 版本号
		AutoStart   bool   `json:"autoStart"`   // 自动启动
		Description string `json:"description"` // 描述文本
	}
	form := Form{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(200, Error400(err))
		return
	}
	match, _ := regexp.Match(semVerRegexExpr, []byte(form.Version))
	if !match {
		c.JSON(200, Error400(fmt.Errorf("version not match server style:%s", form.Version)))
		return
	}
	_, errStat := os.Stat("./apps/")
	if os.IsNotExist(errStat) {
		err := os.Mkdir("./apps/", 0777)
		if err != nil {
			c.JSON(200, Error400(err))
			return
		}
	}
	newUUID := utils.AppUuid()
	// 开始在 ./apps目录下 新建文件
	path := "./apps/" + newUUID + ".lua"
	_, err := os.Create(path)
	if err != nil {
		c.JSON(200, Error400(err))
		return
	}
	err1 := os.WriteFile(path, []byte(fmt.Sprintf(luaTemplate, form.Name,
		form.Version, form.Description, defaultLuaMain)), 0777)
	if err1 != nil {
		c.JSON(200, Error400(err1))
		return
	}
	if err := hs.InsertApp(&MApp{
		UUID:      newUUID,
		Name:      form.Name,
		Version:   form.Version,
		Filepath:  path,
		AutoStart: &form.AutoStart,
	}); err != nil {
		c.JSON(200, Error400(err))
		return
	}
	// 立即加载
	if err := e.LoadApp(typex.NewApplication(
		newUUID, form.Name, form.Version, path)); err != nil {
		glogger.GLogger.Error("App Load failed:", err)
		c.JSON(200, Error400(err))
		return
	}
	// 自启动立即运行
	if form.AutoStart {
		glogger.GLogger.Debugf("App autoStart allowed:%s-%s-%s", newUUID, form.Version, form.Name)
		if err2 := e.StartApp(newUUID); err2 != nil {
			glogger.GLogger.Error("App autoStart failed:", err2)
		}
	}
	c.JSON(200, OkWithData("App create successfully"))
}

/*
*
* Update app
*
 */
func UpdateApp(c *gin.Context, hs *HttpApiServer, e typex.RuleX) {
	type Form struct {
		UUID        string `json:"uuid"`        // uuid
		Name        string `json:"name"`        // 名称
		Version     string `json:"version"`     // 版本号
		AutoStart   bool   `json:"autoStart"`   // 自动启动
		LuaSource   string `json:"luaSource"`   // lua 源码
		Description string `json:"description"` // 描述文本

	}
	form := Form{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(200, Error400(err))
		return
	}
	// 校验版本号
	match, _ := regexp.Match(semVerRegexExpr, []byte(form.Version))
	if !match {
		c.JSON(200, Error400(fmt.Errorf("version not match server style:%s", form.Version)))
		return
	}
	// 校验语法
	if err1 := appstack.ValidateLuaSyntax([]byte(form.LuaSource)); err1 != nil {
		c.JSON(200, Error400(err1))
		return
	}

	if err := hs.UpdateApp(&MApp{
		UUID:      form.UUID,
		Name:      form.Name,
		Version:   form.Version,
		AutoStart: &form.AutoStart,
	}); err != nil {
		c.JSON(200, Error400(err))
		return
	}
	// 最后把文件内容改变了
	path := "./apps/" + form.UUID + ".lua"
	err1 := os.WriteFile(path, []byte(form.LuaSource), 0644)
	if err1 != nil {
		c.JSON(200, Error400(err1))
		return
	}
	//
	if form.AutoStart {
		glogger.GLogger.Debugf("App autoStart allowed:%s-%s-%s", form.UUID, form.Version, form.Name)
		if err2 := e.StartApp(form.UUID); err2 != nil {
			glogger.GLogger.Error("App autoStart failedF:", err2)
		}
	}
	c.JSON(200, OkWithData("App update successfully:"+form.UUID))
}

/*
*
* 启动应用: 用来从数据库里面启动, 有2种情况：
* 1 停止了的, 就需要重启一下
* 2 还未被加载进来的（刚新建），先load后start
 */
func StartApp(c *gin.Context, hs *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	// 检查数据库
	mApp, err := hs.GetAppWithUUID(uuid)
	if err != nil {
		c.JSON(200, Error400(err))
		return
	}
	// 如果内存里面有, 判断状态
	if app := e.GetApp(uuid); app != nil {
		glogger.GLogger.Debug("Already loaded, will try to start:", uuid)
		// 已经启动了就不能再启动
		if app.AppState == 1 {
			c.JSON(200, Error400(fmt.Errorf("app is running now:%s", uuid)))
		}
		if app.AppState == 0 {
			if err := e.StartApp(uuid); err != nil {
				c.JSON(200, Error400(err))
			}
			c.JSON(200, OkWithData("App start successfully:"+uuid))
		}
		return
	}
	// 如果内存里面没有，尝试从配置加载
	glogger.GLogger.Debug("No loaded, will try to load:", uuid)
	if err := e.LoadApp(typex.NewApplication(
		mApp.UUID, mApp.Name, mApp.Version, mApp.Filepath)); err != nil {
		c.JSON(200, Error400(err))
		return
	}
	glogger.GLogger.Debug("App loaded, will try to start:", uuid)
	if err := e.StartApp(uuid); err != nil {
		c.JSON(200, Error400(err))
		return
	}
	c.JSON(200, OkWithData("App start successfully:"+uuid))
}

// 停止, 但是不删除，仅仅是把虚拟机进程给杀死
func StopApp(c *gin.Context, hs *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	if app := e.GetApp(uuid); app != nil {
		if app.AppState == 0 {
			c.JSON(200, Error400(fmt.Errorf("app is stopping now:%s", uuid)))
			return
		}
		if app.AppState == 1 {
			if err := e.StopApp(uuid); err != nil {
				c.JSON(200, OkWithData(err))
				return
			}
			c.JSON(200, OkWithData("app stopped:%s"+uuid))
			return
		}
	}
	c.JSON(200, Error400(fmt.Errorf("app not exists:%s", uuid)))
}

// 删除
func RemoveApp(c *gin.Context, hs *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	// 先把正在运行的给停了
	if app := e.GetApp(uuid); app != nil {
		app.Remove()
	}
	// 内存给清理了
	if err := e.RemoveApp(uuid); err != nil {
		c.JSON(200, Error400(err))
		return
	}
	// Sqlite 配置也给删了
	if err := hs.DeleteApp(uuid); err != nil {
		c.JSON(200, Error400(err))
		return
	}
	// lua的文件也删了
	if err := os.Remove("./apps/" + uuid + ".lua"); err != nil {
		c.JSON(200, Error400(err))
		return
	}
	c.JSON(200, OkWithData(fmt.Errorf("remove app successfully:%s", uuid)))
}
