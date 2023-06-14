package httpserver

import (
	"fmt"

	"github.com/hootrhino/rulex/core"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"

	"github.com/gin-gonic/gin"
)

type ruleVo struct {
	UUID        string   `json:"uuid"`
	FromSource  []string `json:"fromSource"`
	FromDevice  []string `json:"fromDevice"`
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Status      int      `json:"status"`
	Expression  string   `json:"expression"`
	Description string   `json:"description"`
	Actions     string   `json:"actions"`
	Success     string   `json:"success"`
	Failed      string   `json:"failed"`
}

// Get all rules
func Rules(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	if uuid == "" {
		DataList := []ruleVo{}
		allRules, _ := hh.GetAllMRule()
		for _, rule := range allRules {
			DataList = append(DataList, ruleVo{
				UUID:        rule.UUID,
				Name:        rule.Name,
				Type:        rule.Type,
				Status:      1,
				Expression:  rule.Expression,
				Description: rule.Description,
				FromSource:  rule.FromSource,
				FromDevice:  rule.FromDevice,
				Success:     rule.Success,
				Failed:      rule.Failed,
				Actions:     rule.Actions,
			})
		}
		c.JSON(HTTP_OK, OkWithData(DataList))
	} else {
		rule, _ := hh.GetMRuleWithUUID(uuid)
		c.JSON(HTTP_OK, OkWithData(ruleVo{
			UUID:        rule.UUID,
			Name:        rule.Name,
			Type:        rule.Type,
			Expression:  rule.Expression,
			Description: rule.Description,
			FromSource:  rule.FromSource,
			FromDevice:  rule.FromDevice,
			Success:     rule.Success,
			Failed:      rule.Failed,
			Actions:     rule.Actions,
		}))
	}
}

// Create rule
func CreateRule(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	type Form struct {
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
			in, _ := hh.GetMInEndWithUUID(id)
			if in == nil {
				c.JSON(HTTP_OK, Error(`inend not exists: `+id))
				return
			}
		}
	}

	if lenDevices > 0 {
		for _, id := range form.FromDevice {
			in, _ := hh.GetDeviceWithUUID(id)
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
		UUID:        utils.RuleUuid(),
		Expression:  form.Expression,
		Description: form.Description,
		FromSource:  form.FromSource,
		FromDevice:  form.FromDevice,
		Success:     form.Success,
		Failed:      form.Failed,
		Actions:     form.Actions,
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
		}
		// 更新FromSource RULE到Device表中
		for _, id := range form.FromSource {
			InEnd, _ := hh.GetDeviceWithUUID(id)
			if InEnd == nil {
				c.JSON(HTTP_OK, Error(`inend not exists: `+id))
				return
			} else {
				// 去重旧的
				ruleMap := map[string]string{}
				for _, rule := range InEnd.BindRules {
					ruleMap[rule] = rule
				}
				// 追加新的ID
				ruleMap[id] = mRule.UUID
				// 最后ID列表
				BindRules := []string{}
				for _, iid := range ruleMap {
					BindRules = append(BindRules, iid)
				}
				InEnd.BindRules = BindRules
				if err := hh.UpdateMInEnd(InEnd.UUID, &MInEnd{
					BindRules: BindRules,
				}); err != nil {
					c.JSON(HTTP_OK, Error400(err))
					return
				}
			}
		}
		// FromDevice
		for _, id := range form.FromDevice {
			Device, _ := hh.GetDeviceWithUUID(id)
			if Device == nil {
				c.JSON(HTTP_OK, Error(`device not exists: `+id))
				return
			} else {
				// 去重旧的
				ruleMap := map[string]string{}
				for _, rule := range Device.BindRules {
					ruleMap[rule] = rule
				}
				// 追加新的ID
				ruleMap[id] = mRule.UUID
				// 最后ID列表
				BindRules := []string{}
				for _, iid := range ruleMap {
					BindRules = append(BindRules, iid)
				}
				Device.BindRules = BindRules
				if err := hh.UpdateDevice(Device.UUID, &MDevice{
					BindRules: BindRules,
				}); err != nil {
					c.JSON(HTTP_OK, Error400(err))
					return
				}
			}
		}
		// SaveDB
		if err := hh.InsertMRule(mRule); err != nil {
			c.JSON(HTTP_OK, Error400(err))
			return
		}
		c.JSON(HTTP_OK, Ok())
		return
	}
	c.JSON(HTTP_OK, Error("unsupported type:"+form.Type))

}

/*
*
* Update
*
 */
func UpdateRule(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	type Form struct {
		UUID        string   `json:"uuid" binding:"required"` // 如果空串就是新建，非空就是更新
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
		}
		// 更新FromSource RULE到Device表中
		for _, id := range form.FromSource {
			InEnd, _ := hh.GetDeviceWithUUID(id)
			if InEnd == nil {
				c.JSON(HTTP_OK, Error(`inend not exists: `+id))
				return
			}
			// 去重旧的
			ruleMap := map[string]string{}
			for _, rule := range InEnd.BindRules {
				ruleMap[rule] = rule
			}
			// 追加新的ID
			ruleMap[id] = mRule.UUID
			// 最后ID列表
			BindRules := []string{}
			for _, iid := range ruleMap {
				BindRules = append(BindRules, iid)
			}
			InEnd.BindRules = BindRules
			if err := hh.UpdateMInEnd(InEnd.UUID, &MInEnd{
				BindRules: BindRules,
			}); err != nil {
				c.JSON(HTTP_OK, Error400(err))
				return
			}

		}
		// FromDevice
		for _, id := range form.FromDevice {
			Device, _ := hh.GetDeviceWithUUID(id)
			if Device == nil {
				c.JSON(HTTP_OK, Error(`device not exists: `+id))
				return
			}
			// 去重旧的
			ruleMap := map[string]string{}
			for _, rule := range Device.BindRules {
				ruleMap[rule] = rule
			}
			// 追加新的ID
			ruleMap[id] = mRule.UUID
			// 最后ID列表
			BindRules := []string{}
			for _, iid := range ruleMap {
				BindRules = append(BindRules, iid)
			}
			Device.BindRules = BindRules
			if err := hh.UpdateDevice(Device.UUID, &MDevice{
				BindRules: BindRules,
			}); err != nil {
				c.JSON(HTTP_OK, Error400(err))
				return
			}

		}
		if err := hh.UpdateMRule(form.UUID, mRule); err != nil {
			c.JSON(HTTP_OK, Error400(err))
			return
		}
		c.JSON(HTTP_OK, Ok())
		return
	}
	c.JSON(HTTP_OK, Error("rule not exists:"+form.UUID))

}

// Delete rule by UUID
func DeleteRule(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	mRule, err0 := hh.GetMRule(uuid)
	if err0 != nil {
		c.JSON(HTTP_OK, Error400(err0))
		return
	}
	// 更新FromSource RULE到Device表中
	for _, id := range mRule.FromSource {
		InEnd, _ := hh.GetDeviceWithUUID(id)
		if InEnd == nil {
			c.JSON(HTTP_OK, Error(`inend not exists: `+id))
			return
		}
		// 去重旧的
		ruleMap := map[string]string{}
		for _, rule := range InEnd.BindRules {
			ruleMap[rule] = rule
		}
		// 删除ID
		delete(ruleMap, mRule.UUID)
		// 最后ID列表
		BindRules := []string{}
		for _, iid := range ruleMap {
			BindRules = append(BindRules, iid)
		}
		InEnd.BindRules = BindRules
		if err := hh.UpdateMInEnd(InEnd.UUID, &MInEnd{
			BindRules: BindRules,
		}); err != nil {
			c.JSON(HTTP_OK, Error400(err))
			return
		}
	}
	// FromDevice
	for _, id := range mRule.FromDevice {
		Device, _ := hh.GetDeviceWithUUID(id)
		if Device == nil {
			c.JSON(HTTP_OK, Error(`device not exists: `+id))
			return
		}
		// 去重旧的
		ruleMap := map[string]string{}
		for _, rule := range Device.BindRules {
			ruleMap[rule] = rule
		}
		// 删除ID
		delete(ruleMap, mRule.UUID)
		// 最后ID列表
		BindRules := []string{}
		for _, iid := range ruleMap {
			BindRules = append(BindRules, iid)
		}
		Device.BindRules = BindRules
		if err := hh.UpdateDevice(Device.UUID, &MDevice{
			BindRules: BindRules,
		}); err != nil {
			c.JSON(HTTP_OK, Error400(err))
			return
		}
	}
	if err := hh.DeleteMRule(uuid); err != nil {
		c.JSON(HTTP_OK, Error400(err))
		return
	}

	e.RemoveRule(uuid)
	c.JSON(HTTP_OK, Ok())
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
