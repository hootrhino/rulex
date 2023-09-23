package apis

import (
	"fmt"
	"regexp"

	common "github.com/hootrhino/rulex/plugin/http_server/common"
	"github.com/hootrhino/rulex/plugin/http_server/model"
	"github.com/hootrhino/rulex/plugin/http_server/service"

	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rulex/component/appstack"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
)

/*
*
* 其实这个结构体扮演的角色VO层
*
 */
type appStackDto struct {
	UUID        string `json:"uuid,omitempty"`      // 名称
	Name        string `json:"name,omitempty"`      // 名称
	Version     string `json:"version,omitempty"`   // 版本号
	AutoStart   *bool  `json:"autoStart,omitempty"` // 自动启动
	AppState    int    `json:"appState,omitempty"`  // 状态: 1 运行中, 0 停止
	Type        string `json:"type,omitempty"`      // 默认就是lua, 留个扩展以后可能支持别的
	LuaSource   string `json:"luaSource,omitempty"`
	Description string `json:"description,omitempty"`
}

/*
*
* APP 详情
*
 */
func AppDetail(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	// uuid
	appInfo, err1 := service.GetMAppWithUUID(uuid)
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400EmptyObj(err1))
		return
	}
	web_data := appStackDto{
		UUID:      appInfo.UUID,
		Name:      appInfo.Name,
		Version:   appInfo.Version,
		AutoStart: appInfo.AutoStart,
		Type:      "lua",
		AppState: func() int {
			if a := appstack.GetApp(appInfo.UUID); a != nil {
				return int(a.AppState)
			}
			return 0
		}(),
		Description: appInfo.Description,
		LuaSource:   appInfo.LuaSource,
	}
	c.JSON(common.HTTP_OK, common.OkWithData(web_data))
}

// 列表
func Apps(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	// 从配置拿App
	if uuid == "" {
		result := []appStackDto{}
		for _, app := range appstack.AllApp() {
			web_data := appStackDto{
				UUID:      app.UUID,
				Name:      app.Name,
				Version:   app.Version,
				AutoStart: &app.AutoStart,
				Type:      "lua",
				AppState: func() int {
					if a := appstack.GetApp(app.UUID); a != nil {
						return int(a.AppState)
					}
					return 0
				}(),
				Description: "",
			}
			result = append(result, web_data)
		}
		c.JSON(common.HTTP_OK, common.OkWithData(result))
		return
	}
	// uuid
	appInfo, err1 := service.GetMAppWithUUID(uuid)
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	web_data := appStackDto{
		UUID:      appInfo.UUID,
		Name:      appInfo.Name,
		Version:   appInfo.Version,
		AutoStart: appInfo.AutoStart,
		Type:      "lua",
		AppState: func() int {
			if a := appstack.GetApp(appInfo.UUID); a != nil {
				return int(a.AppState)
			}
			return 0
		}(),
		Description: appInfo.Description,
		LuaSource:   appInfo.LuaSource,
	}
	c.JSON(common.HTTP_OK, common.OkWithData(web_data))

}

/*
*
* 直接新建一个文件，文件名为 UUID.lua
*
 */
const semVerRegexExpr = `^(0|[1-9]+[0-9]*)\.(0|[1-9]+[0-9]*)\.(0|[1-9]+[0-9]*)(-(0|[1-9A-Za-z-][0-9A-Za-z-]*)(\.[0-9A-Za-z-]+)*)?(\+[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?$`
const luaTemplate = `
--
-- App use lua syntax, goto https://hootrhino.github.io for more document
-- APPID: %s
--
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
	applib:Debug("Hello World:" .. applib:Time())
	return 0
end
`

func CreateApp(c *gin.Context, ruleEngine typex.RuleX) {
	form := appStackDto{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	match, _ := regexp.Match(semVerRegexExpr, []byte(form.Version))
	if !match {
		c.JSON(common.HTTP_OK, common.Error400(fmt.Errorf("version not match server style:%s", form.Version)))
		return
	}
	newUUID := utils.AppUuid()
	mAPP := &model.MApp{
		UUID:    newUUID,
		Name:    form.Name,
		Version: form.Version,
		LuaSource: fmt.Sprintf(luaTemplate,
			newUUID, form.Name, form.Version, form.Description, defaultLuaMain),
		AutoStart:   form.AutoStart,
		Description: form.Description,
	}
	if err := service.InsertApp(mAPP); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 立即加载
	if err := appstack.LoadApp(
		typex.NewApplication(newUUID, form.Name, form.Version), mAPP.LuaSource); err != nil {
		glogger.GLogger.Error("app Load failed:", err)
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 自启动立即运行
	if *form.AutoStart {
		glogger.GLogger.Debugf("App autoStart allowed:%s-%s-%s", newUUID, form.Version, form.Name)
		if err2 := appstack.StartApp(newUUID); err2 != nil {
			glogger.GLogger.Error("App autoStart failed:", err2)
		}
	}
	c.JSON(common.HTTP_OK, common.OkWithData("app create successfully"))
}

/*
*
* Update app
*
 */
func UpdateApp(c *gin.Context, ruleEngine typex.RuleX) {
	form := appStackDto{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 校验版本号
	match, _ := regexp.Match(semVerRegexExpr, []byte(form.Version))
	if !match {
		c.JSON(common.HTTP_OK, common.Error400(fmt.Errorf("version not match server style:%s", form.Version)))
		return
	}
	// 校验语法
	if err1 := appstack.ValidateLuaSyntax([]byte(form.LuaSource)); err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}

	if err := service.UpdateApp(&model.MApp{
		UUID:        form.UUID,
		Name:        form.Name,
		Version:     form.Version,
		AutoStart:   form.AutoStart,
		LuaSource:   form.LuaSource,
		Description: form.Description,
	}); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 如果内存里面有, 先把内存里的清理了
	if app := appstack.GetApp(form.UUID); app != nil {
		glogger.GLogger.Debug("Already loaded, will try to stop:", form.UUID)
		// 已经启动了就不能再启动
		if app.AppState == 1 {
			appstack.StopApp(form.UUID)
		}
		appstack.RemoveApp(app.UUID)
	}
	//
	if *form.AutoStart {
		glogger.GLogger.Debugf("App autoStart allowed:%s-%s-%s", form.UUID, form.Version, form.Name)
		// 必须先load后start
		if err := appstack.LoadApp(typex.NewApplication(
			form.UUID, form.Name, form.Version), form.LuaSource); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		if err2 := appstack.StartApp(form.UUID); err2 != nil {
			glogger.GLogger.Error("App autoStart failed:", err2)
			c.JSON(common.HTTP_OK, common.Error400(err2))
			return
		}
	}
	c.JSON(common.HTTP_OK, common.OkWithData("app update successfully:"+form.UUID))
}

/*
*
* 启动应用: 用来从数据库里面启动, 有2种情况：
* 1 停止了的, 就需要重启一下
* 2 还未被加载进来的（刚新建），先load后start
 */
func StartApp(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	// 检查数据库
	mApp, err := service.GetMAppWithUUID(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 如果内存里面有, 判断状态
	if app := appstack.GetApp(uuid); app != nil {
		glogger.GLogger.Debug("Already loaded, will try to start:", uuid)
		// 已经启动了就不能再启动
		if app.AppState == 1 {
			c.JSON(common.HTTP_OK, common.Error400(fmt.Errorf("app is running now:%s", uuid)))
		}
		if app.AppState == 0 {
			if err := appstack.StartApp(uuid); err != nil {
				c.JSON(common.HTTP_OK, common.Error400(err))
			} else {
				c.JSON(common.HTTP_OK, common.OkWithData("app start successfully:"+uuid))
			}
		}
		return
	}
	// 如果内存里面没有，尝试从配置加载
	glogger.GLogger.Debug("No loaded, will try to load:", uuid)
	if err := appstack.LoadApp(typex.NewApplication(
		mApp.UUID, mApp.Name, mApp.Version), mApp.LuaSource); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	glogger.GLogger.Debug("app loaded, will try to start:", uuid)
	if err := appstack.StartApp(uuid); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData("app start successfully:"+uuid))
}

// 停止, 但是不删除，仅仅是把虚拟机进程给杀死
func StopApp(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	if app := appstack.GetApp(uuid); app != nil {
		if app.AppState == 0 {
			c.JSON(common.HTTP_OK, common.Error400(fmt.Errorf("app is stopping now:%s", uuid)))
			return
		}
		if app.AppState == 1 {
			if err := appstack.StopApp(uuid); err != nil {
				c.JSON(common.HTTP_OK, common.Error400(err))
				return
			}
			c.JSON(common.HTTP_OK, common.OkWithData("app stopped:%s"+uuid))
			return
		}
	}
	c.JSON(common.HTTP_OK, common.Error400(fmt.Errorf("app not exists:%s", uuid)))
}

// 删除
func RemoveApp(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	// 先把正在运行的给停了
	if app := appstack.GetApp(uuid); app != nil {
		app.Remove()
	}
	// 内存给清理了
	if err := appstack.RemoveApp(uuid); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// Sqlite 配置也给删了
	if err := service.DeleteApp(uuid); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(fmt.Sprintf("remove app successfully:%s", uuid)))
}
