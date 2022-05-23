package httpserver

import (
	"fmt"
	"rulex/core"
	"rulex/typex"
	"rulex/utils"

	"github.com/gin-gonic/gin"
)

//
// Get all rules
//
func Rules(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {

	uuid, _ := c.GetQuery("uuid")
	if uuid == "" {
		data := []interface{}{}
		allRules := e.AllRule()
		allRules.Range(func(key, value interface{}) bool {
			data = append(data, value)
			return true
		})
		c.JSON(200, Result{
			Code: 200,
			Msg:  SUCCESS,
			Data: data,
		})
	} else {
		c.JSON(200, Result{
			Code: 200,
			Msg:  SUCCESS,
			Data: e.GetRule(uuid),
		})
	}
}

//
// Create rule
//
func CreateRule(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	type Form struct {
		UUID        string   `json:"uuid"` // 如果空串就是新建，非空就是更新
		FromSource  []string `json:"fromSource" binding:"required"`
		FromDevice  []string `json:"fromDevice" binding:"required"`
		Name        string   `json:"name" binding:"required"`
		Description string   `json:"description"`
		Actions     string   `json:"actions"`
		Success     string   `json:"success"`
		Failed      string   `json:"failed"`
	}
	form := Form{}

	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(200, Error400(err))
		return
	}
	lenSources := len(form.FromSource)
	lenDevices := len(form.FromDevice)
	if lenSources > 0 {
		for _, id := range form.FromSource {
			in := e.GetInEnd(id)
			if in == nil {
				c.JSON(200, Error(`inend not exists: `+id))
				return
			}
		}
	}

	if lenDevices > 0 {
		for _, id := range form.FromDevice {
			in := e.GetDevice(id)
			if in == nil {
				c.JSON(200, Error(`device not exists: `+id))
				return
			}
		}
	}
	if (lenDevices < 1) && (lenSources < 1) {
		c.JSON(200, Error(`必须有一个数据输入项`))
		return
	}

	// tmpRule 是一个一次性的临时rule，用来验证规则，这么做主要是为了防止真实Lua Vm 被污染
	tmpRule := typex.NewRule(nil,
		"tmpRule",
		"tmpRule",
		"tmpRule",
		[]string{},
		[]string{},
		form.Success,
		form.Actions,
		form.Failed)

	if err := core.VerifyCallback(tmpRule); err != nil {
		c.JSON(200, Error400(err))
		return
	}
	// 如果是更新操作, 先删除规则
	if form.UUID != "" {
		if err1 := hh.DeleteMRule(form.UUID); err1 != nil {
			c.JSON(200, Error400(err1))
			return
		} else {
			e.RemoveRule(form.UUID)
		}
	}
	//
	mRule := &MRule{
		UUID:        utils.MakeUUID("RULE"),
		Name:        form.Name,
		Description: form.Description,
		FromSource:  form.FromSource,
		FromDevice:  form.FromDevice,
		Success:     form.Success,
		Failed:      form.Failed,
		Actions:     form.Actions,
	}
	if err := hh.InsertMRule(mRule); err != nil {
		c.JSON(200, Error400(err))
		return
	}
	rule := typex.NewRule(hh.ruleEngine,
		mRule.UUID,
		mRule.Name,
		mRule.Description,
		mRule.FromSource,
		mRule.FromDevice,
		mRule.Success,
		mRule.Actions,
		mRule.Failed)
	if err := e.LoadRule(rule); err != nil {
		c.JSON(200, Error400(err))
	} else {
		c.JSON(200, Ok())
	}
	return
}

//
// Delete rule by UUID
//
func DeleteRule(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	_, err0 := hh.GetMRule(uuid)
	if err0 != nil {
		c.JSON(200, Error400(err0))
		return
	}
	if err1 := hh.DeleteMRule(uuid); err1 != nil {
		c.JSON(200, Error400(err1))
	} else {
		e.RemoveRule(uuid)
		c.JSON(200, Ok())
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
		c.JSON(200, Error400(err))
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
	if err := core.VerifyCallback(tmpRule); err != nil {
		c.JSON(200, Error400(err))
	} else {
		c.JSON(200, Ok())
	}

}

/*
*
* 测试脚本执行效果
*
 */
func TestLuaCallback(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid") // InEnd
	data, _ := c.GetQuery("data") // Data
	_, err0 := hh.GetMRule(uuid)
	if err0 != nil {
		c.JSON(200, Error400(err0))
		return
	}
	value, ok := e.AllInEnd().Load(uuid)
	if !ok {
		c.JSON(200, Error(fmt.Sprintf("'InEnd' not exists: %v", uuid)))
		return
	}
	_, err1 := e.Work((value).(*typex.InEnd), data)
	if err1 != nil {
		c.JSON(200, Error400(err1))
	}
	c.JSON(200, Ok())
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
		c.JSON(200, Error400(err0))
		return
	}
	value, ok := e.AllOutEnd().Load(uuid)
	if !ok {
		c.JSON(200, Error((fmt.Sprintf("'OutEnd' not exists: %v", uuid))))
		return
	}
	e.PushOutQueue((value).(*typex.OutEnd), data)
	c.JSON(200, Ok())
}
