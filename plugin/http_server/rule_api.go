package httpserver

import (
	"fmt"

	"github.com/hootrhino/rulex/core"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"

	"github.com/gin-gonic/gin"
)

// Get all rules
func Rules(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {

	uuid, _ := c.GetQuery("uuid")
	if uuid == "" {
		data := []interface{}{}
		allRules := e.AllRule()
		allRules.Range(func(key, value interface{}) bool {
			data = append(data, value)
			return true
		})
		c.JSON(HTTP_OK, OkWithData(data))
	} else {
		c.JSON(HTTP_OK, OkWithData(e.GetRule(uuid)))
	}
}

// Create rule
func CreateRule(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	type Form struct {
		UUID        string   `json:"uuid"` // 如果空串就是新建，非空就是更新
		FromSource  []string `json:"fromSource" binding:"required"`
		FromDevice  []string `json:"fromDevice" binding:"required"`
		Name        string   `json:"name" binding:"required"`
		Type        string   `json:"type"`
		Expression  string   `json:"expression"`
		Description string   `json:"description"`
		Actions     string   `json:"actions"`
		Success     string   `json:"success"`
		Failed      string   `json:"failed"`
	}
	form := Form{Type: "lua"}

	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(HTTP_OK, Error400(err))
		return
	}
	if !utils.SContains([]string{"lua", "expr"}, form.Type) {
		c.JSON(HTTP_OK, Error(`rule type must one of 'lua' or 'expr':`+form.Type))
		return
	}
	lenSources := len(form.FromSource)
	lenDevices := len(form.FromDevice)
	if lenSources > 0 {
		for _, id := range form.FromSource {
			in := e.GetInEnd(id)
			if in == nil {
				c.JSON(HTTP_OK, Error(`inend not exists: `+id))
				return
			}
		}
	}

	if lenDevices > 0 {
		for _, id := range form.FromDevice {
			in := e.GetDevice(id)
			if in == nil {
				c.JSON(HTTP_OK, Error(`device not exists: `+id))
				return
			}
		}
	}
	// tmpRule 是一个一次性的临时rule，用来验证规则，这么做主要是为了防止真实Lua Vm 被污染
	tmpRule := typex.NewRule(nil, "_", "_", "_", []string{}, []string{},
		form.Success, form.Actions, form.Failed)
	if form.Type == "lua" {
		if err := core.VerifyLuaSyntax(tmpRule); err != nil {
			c.JSON(HTTP_OK, Error400(err))
			return
		}
	}
	if form.Type == "expr" {
		if err := core.VerifyExprSyntax(tmpRule); err != nil {
			c.JSON(HTTP_OK, Error400(err))
			return
		}
	}
	//
	mRule := &MRule{
		Name:        form.Name,
		Type:        form.Type,
		Expression:  form.Expression,
		Description: form.Description,
		FromSource:  form.FromSource,
		FromDevice:  form.FromDevice,
		Success:     form.Success,
		Failed:      form.Failed,
		Actions:     form.Actions,
	}
	// 更新操作
	if form.UUID != "" {
		mRule.UUID = form.UUID
		if err := hh.UpdateMRule(form.UUID, mRule); err != nil {
			c.JSON(HTTP_OK, Error400(err))
			return
		}
	}
	// 新建操作
	if form.UUID == "" {
		mRule.UUID = utils.RuleUuid()
		if err := hh.InsertMRule(mRule); err != nil {
			c.JSON(HTTP_OK, Error400(err))
			return
		}
	}

	if form.Type == "lua" {
		rule := typex.NewLuaRule(hh.ruleEngine,
			mRule.UUID,
			mRule.Name,
			mRule.Description,
			mRule.FromSource,
			mRule.FromDevice,
			mRule.Success,
			mRule.Actions,
			mRule.Failed)
		e.RemoveRule(rule.UUID)
		if err := e.LoadRule(rule); err != nil {
			c.JSON(HTTP_OK, Error400(err))
			return
		} else {
			c.JSON(HTTP_OK, Ok())
			return
		}
	}
	if form.Type == "expr" {
		rule := typex.NewExprRule(hh.ruleEngine,
			mRule.UUID,
			mRule.Name,
			mRule.Type,
			mRule.Expression,
			mRule.Description,
			mRule.FromSource,
			mRule.FromDevice,
			mRule.Success,
			mRule.Actions,
			mRule.Failed)
		e.RemoveRule(rule.UUID)
		if err := e.LoadRule(rule); err != nil {
			c.JSON(HTTP_OK, Error400(err))
			return
		} else {
			c.JSON(HTTP_OK, Ok())
			return
		}
	}
	c.JSON(HTTP_OK, Error("unsupported type:"+form.Type))

}

// Delete rule by UUID
func DeleteRule(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	_, err0 := hh.GetMRule(uuid)
	if err0 != nil {
		c.JSON(HTTP_OK, Error400(err0))
		return
	}
	if err1 := hh.DeleteMRule(uuid); err1 != nil {
		c.JSON(HTTP_OK, Error400(err1))
		return

	} else {
		e.RemoveRule(uuid)
		c.JSON(HTTP_OK, Ok())
		return
	}

}

/*
*
* 验证lua语法
*
 */
func ValidateLuaSyntax(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	type Form struct {
		FromSource []string `json:"fromSource" binding:"required"`
		FromDevice []string `json:"fromDevice" binding:"required"`
		Actions    string   `json:"actions" binding:"required"`
		Success    string   `json:"success" binding:"required"`
		Failed     string   `json:"failed" binding:"required"`
	}
	form := Form{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(HTTP_OK, Error400(err))
		return
	}
	tmpRule := typex.NewRule(
		nil, // 不需要该字段
		"",  // 不需要该字段
		"",  // 不需要该字段
		"",  // 不需要该字段
		form.FromSource,
		form.FromDevice,
		form.Success,
		form.Actions,
		form.Failed)
	if err := core.VerifyLuaSyntax(tmpRule); err != nil {
		c.JSON(HTTP_OK, Error400(err))
	} else {
		c.JSON(HTTP_OK, Ok())
	}

}

/*
*
* 测试脚本执行效果
*
 */
func TestSourceCallback(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid") // InEnd
	data, _ := c.GetQuery("data") // Data
	_, err0 := hh.GetMRule(uuid)
	if err0 != nil {
		c.JSON(HTTP_OK, Error400(err0))
		return
	}
	value, ok := e.AllInEnd().Load(uuid)
	if !ok {
		c.JSON(HTTP_OK, Error(fmt.Sprintf("'InEnd' not exists: %v", uuid)))
		return
	}
	_, err1 := e.WorkInEnd((value).(*typex.InEnd), data)
	if err1 != nil {
		c.JSON(HTTP_OK, Error400(err1))
		return
	}
	c.JSON(HTTP_OK, Ok())
}

/*
*
* 测试 OutEnd 的结果
*
 */
func TestOutEndCallback(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid") // OutEnd
	data, _ := c.GetQuery("data") // Data
	_, err0 := hh.GetMRule(uuid)
	if err0 != nil {
		c.JSON(HTTP_OK, Error400(err0))
		return
	}
	value, ok := e.AllOutEnd().Load(uuid)
	if !ok {
		c.JSON(HTTP_OK, Error((fmt.Sprintf("'OutEnd' not exists: %v", uuid))))
		return
	}
	c.JSON(HTTP_OK, OkWithData(e.PushOutQueue((value).(*typex.OutEnd), data)))
}

/*
*
* 测试 Device 的结果
*
 */
func TestDeviceCallback(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid") // Device
	data, _ := c.GetQuery("data") // Data, Read or write
	_, err0 := hh.GetMRule(uuid)
	if err0 != nil {
		c.JSON(HTTP_OK, Error400(err0))
		return
	}
	value, ok := e.AllDevices().Load(uuid)
	if !ok {
		c.JSON(HTTP_OK, Error((fmt.Sprintf("'Device' not exists: %v", uuid))))
		return
	}
	c.JSON(HTTP_OK, OkWithData(e.PushDeviceQueue((value).(*typex.Device), data)))
}
