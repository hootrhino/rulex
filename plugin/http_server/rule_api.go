package httpserver

import (
	"fmt"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/interqueue"
	common "github.com/hootrhino/rulex/plugin/http_server/common"
	"github.com/hootrhino/rulex/plugin/http_server/model"
	"github.com/sirupsen/logrus"

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

func RuleDetail(c *gin.Context, hh *HttpApiServer) {
	uuid, _ := c.GetQuery("uuid")
	rule, err := hh.GetMRuleWithUUID(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400EmptyObj(err))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(ruleVo{
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
	}))
}

// Get all rules
func Rules(c *gin.Context, hh *HttpApiServer) {
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
		c.JSON(common.HTTP_OK, common.OkWithData(DataList))
	} else {
		rule, err := hh.GetMRuleWithUUID(uuid)
		if err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		c.JSON(common.HTTP_OK, common.OkWithData(ruleVo{
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
		}))
	}
}

// Create rule
func CreateRule(c *gin.Context, hh *HttpApiServer) {
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
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if !utils.SContains([]string{"lua"}, form.Type) {
		c.JSON(common.HTTP_OK, common.Error(`rule type must be 'lua' but now is:`+form.Type))
		return
	}
	lenSources := len(form.FromSource)
	lenDevices := len(form.FromDevice)
	if lenSources > 0 {
		for _, id := range form.FromSource {
			in, _ := hh.GetMInEndWithUUID(id)
			if in == nil {
				c.JSON(common.HTTP_OK, common.Error(`inend not exists: `+id))
				return
			}
		}
	}

	if lenDevices > 0 {
		for _, id := range form.FromDevice {
			in, _ := hh.GetMDeviceWithUUID(id)
			if in == nil {
				c.JSON(common.HTTP_OK, common.Error(`device not exists: `+id))
				return
			}
		}
	}
	// tmpRule 是一个一次性的临时rule，用来验证规则，这么做主要是为了防止真实Lua Vm 被污染
	tmpRule := typex.NewRule(nil, "_", "_", "_", []string{}, []string{},
		form.Success, form.Actions, form.Failed)

	if err := core.VerifyLuaSyntax(tmpRule); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}

	//
	mRule := &model.MRule{
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

	if form.Type != "lua" {
		// 未来可能支持别的脚本
	}
	rule := typex.NewLuaRule(hh.ruleEngine,
		mRule.UUID,
		mRule.Name,
		mRule.Description,
		mRule.FromSource,
		mRule.FromDevice,
		mRule.Success,
		mRule.Actions,
		mRule.Failed)
	hh.ruleEngine.RemoveRule(rule.UUID)
	if err := hh.ruleEngine.LoadRule(rule); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 更新FromSource RULE到Device表中
	for _, inId := range form.FromSource {
		InEnd, _ := hh.GetMInEndWithUUID(inId)
		if InEnd == nil {
			c.JSON(common.HTTP_OK, common.Error(`inend not exists: `+inId))
			return
		}
		// 去重旧的
		ruleMap := map[string]string{}
		for _, rule := range InEnd.BindRules {
			ruleMap[rule] = rule
		}
		// 追加新的ID
		ruleMap[inId] = mRule.UUID
		// 最后ID列表
		BindRules := []string{}
		for _, iid := range ruleMap {
			BindRules = append(BindRules, iid)
		}
		InEnd.BindRules = BindRules
		if err := hh.UpdateMInEnd(InEnd.UUID, &model.MInEnd{
			BindRules: BindRules,
		}); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		// SaveDB
		if err := hh.InsertMRule(mRule); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		// LoadNewest!!!
		if err := hh.LoadNewestInEnd(inId); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}

	}

	// FromDevice
	for _, devId := range form.FromDevice {
		Device, _ := hh.GetMDeviceWithUUID(devId)
		if Device == nil {
			c.JSON(common.HTTP_OK, common.Error(`device not exists: `+devId))
			return
		}
		// 去重旧的
		ruleMap := map[string]string{}
		for _, rule := range Device.BindRules {
			ruleMap[rule] = rule
		}
		// 追加新的ID
		ruleMap[devId] = mRule.UUID
		// 最后ID列表
		BindRules := []string{}
		for _, iid := range ruleMap {
			BindRules = append(BindRules, iid)
		}
		Device.BindRules = BindRules
		if err := hh.UpdateDevice(Device.UUID, &model.MDevice{
			BindRules: BindRules,
		}); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		// SaveDB
		if err := hh.InsertMRule(mRule); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		// LoadNewest!!!
		if err := hh.LoadNewestDevice(devId); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}

	}
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
* Update
*
 */
func UpdateRule(c *gin.Context, hh *HttpApiServer) {
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
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if !utils.SContains([]string{"lua"}, form.Type) {
		c.JSON(common.HTTP_OK, common.Error(`rule type must be 'lua' but now is:`+form.Type))
		return
	}
	lenSources := len(form.FromSource)
	lenDevices := len(form.FromDevice)
	if lenSources > 0 {
		for _, id := range form.FromSource {
			in := hh.ruleEngine.GetInEnd(id)
			if in == nil {
				c.JSON(common.HTTP_OK, common.Error(`inend not exists: `+id))
				return
			}
		}
	}

	if lenDevices > 0 {
		for _, id := range form.FromDevice {
			in := hh.ruleEngine.GetDevice(id)
			if in == nil {
				c.JSON(common.HTTP_OK, common.Error(`device not exists: `+id))
				return
			}
		}
	}
	// tmpRule 是一个一次性的临时rule，用来验证规则，这么做主要是为了防止真实Lua Vm 被污染
	tmpRule := typex.NewRule(nil, "_", "_", "_", []string{}, []string{},
		form.Success, form.Actions, form.Failed)

	if err := core.VerifyLuaSyntax(tmpRule); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}

	if form.Type == "lua" {
		mRule, err := hh.GetMRuleWithUUID(form.UUID)
		if err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		rule := typex.NewLuaRule(hh.ruleEngine,
			mRule.UUID,
			mRule.Name,
			mRule.Description,
			mRule.FromSource,
			mRule.FromDevice,
			mRule.Success,
			mRule.Actions,
			mRule.Failed)
		hh.ruleEngine.RemoveRule(rule.UUID)
		if err := hh.ruleEngine.LoadRule(rule); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		// SaveDB
		//
		if err := hh.UpdateMRule(mRule.UUID, &model.MRule{
			Name:        form.Name,
			Type:        form.Type,
			Expression:  form.Expression,
			Description: form.Description,
			FromSource:  form.FromSource,
			FromDevice:  form.FromDevice,
			Success:     form.Success,
			Failed:      form.Failed,
			Actions:     form.Actions,
		}); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		// 更新FromSource RULE到Device表中
		for _, inId := range form.FromSource {
			if inId != "" {
				InEnd, _ := hh.GetMInEndWithUUID(inId)
				if InEnd == nil {
					c.JSON(common.HTTP_OK, common.Error(`inend not exists: `+inId))
					return
				}
				// 去重旧的
				ruleMap := map[string]string{}
				for _, rule := range InEnd.BindRules {
					ruleMap[rule] = rule
				}
				// 追加新的ID
				ruleMap[inId] = mRule.UUID
				// 最后ID列表
				BindRules := []string{}
				for _, iid := range ruleMap {
					BindRules = append(BindRules, iid)
				}
				InEnd.BindRules = BindRules
				if err := hh.UpdateMInEnd(InEnd.UUID, &model.MInEnd{
					BindRules: BindRules,
				}); err != nil {
					c.JSON(common.HTTP_OK, common.Error400(err))
					return
				}
				// LoadNewest!!!
				if err := hh.LoadNewestInEnd(inId); err != nil {
					c.JSON(common.HTTP_OK, common.Error400(err))
					return
				}
			}

		}
		// FromDevice
		for _, devId := range form.FromDevice {
			if devId != "" {
				Device, _ := hh.GetMDeviceWithUUID(devId)
				if Device == nil {
					c.JSON(common.HTTP_OK, common.Error(`device not exists: `+devId))
					return
				}
				// 去重旧的
				ruleMap := map[string]string{}
				for _, rule := range Device.BindRules {
					ruleMap[rule] = rule
				}
				// 追加新的ID
				ruleMap[devId] = mRule.UUID
				// 最后ID列表
				BindRules := []string{}
				for _, iid := range ruleMap {
					BindRules = append(BindRules, iid)
				}
				Device.BindRules = BindRules
				if err := hh.UpdateDevice(Device.UUID, &model.MDevice{
					BindRules: BindRules,
				}); err != nil {
					c.JSON(common.HTTP_OK, common.Error400(err))
					return
				}
				// LoadNewest!!!
				if err := hh.LoadNewestDevice(devId); err != nil {
					c.JSON(common.HTTP_OK, common.Error400(err))
					return
				}
			}
		}

		c.JSON(common.HTTP_OK, common.Ok())
		return
	}
	c.JSON(common.HTTP_OK, common.Error("rule type invalid:"+form.Type))

}

// Delete rule by UUID
func DeleteRule(c *gin.Context, hh *HttpApiServer) {
	uuid, _ := c.GetQuery("uuid")
	mRule, err0 := hh.GetMRule(uuid)
	if err0 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err0))
		return
	}
	// 更新FromSource RULE到Device表中
	for _, id := range mRule.FromSource {
		InEnd, _ := hh.GetMInEndWithUUID(id)
		if InEnd == nil {
			c.JSON(common.HTTP_OK, common.Error(`inend not exists: `+id))
			return
		}
		// 去重旧的
		BindRules := []string{}
		for _, iid := range InEnd.BindRules {
			if iid != mRule.UUID {
				BindRules = append(BindRules, iid)
			}
		}
		if err := hh.UpdateMInEnd(InEnd.UUID, &model.MInEnd{
			BindRules: BindRules,
		}); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
	}
	// FromDevice
	for _, devId := range mRule.FromDevice {
		Device, _ := hh.GetMDeviceWithUUID(devId)
		if Device == nil {
			c.JSON(common.HTTP_OK, common.Error(`device not exists: `+devId))
			return
		}
		// 去重旧的
		BindRules := []string{}
		for _, iid := range Device.BindRules {
			if iid != mRule.UUID {
				BindRules = append(BindRules, iid)
			}
		}
		if err := hh.UpdateDevice(Device.UUID, &model.MDevice{
			BindRules: BindRules,
		}); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}

	}
	if err := hh.DeleteMRule(uuid); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	hh.ruleEngine.RemoveRule(uuid)
	//
	// 内存里的数据刷新完了以后更新数据库，最后重启
	//
	for _, devId := range mRule.FromDevice {
		if err := hh.LoadNewestDevice(devId); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
	}
	for _, devId := range mRule.FromSource {
		if err := hh.LoadNewestInEnd(devId); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
	}
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
* 验证lua语法
*
 */
func ValidateLuaSyntax(c *gin.Context, hh *HttpApiServer) {
	type Form struct {
		FromSource []string `json:"fromSource" binding:"required"`
		FromDevice []string `json:"fromDevice" binding:"required"`
		Actions    string   `json:"actions" binding:"required"`
		Success    string   `json:"success" binding:"required"`
		Failed     string   `json:"failed" binding:"required"`
	}
	form := Form{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
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
		c.JSON(common.HTTP_OK, common.Error400(err))
	} else {
		c.JSON(common.HTTP_OK, common.Ok())
	}

}

/*
*
* 测试脚本执行效果
*
 */
func TestSourceCallback(c *gin.Context, hh *HttpApiServer) {
	type Form struct {
		UUID     string `json:"uuid"`
		TestData string `json:"testData"`
	}
	form := Form{}
	if err := c.BindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}

	inend := hh.ruleEngine.GetInEnd(form.UUID)
	if inend == nil {
		c.JSON(common.HTTP_OK, common.Error(fmt.Sprintf("'InEnd' not exists: %v", form.UUID)))
		return
	}
	glogger.GLogger.WithFields(logrus.Fields{
		"topic": "rule/test/" + form.UUID,
	}).Debug(form.TestData)
	err1 := interqueue.DefaultDataCacheQueue.PushInQueue(inend, form.TestData)
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
* 测试 OutEnd 的结果
*
 */
func TestOutEndCallback(c *gin.Context, hh *HttpApiServer) {
	type Form struct {
		UUID     string `json:"uuid"`
		TestData string `json:"testData"`
	}
	form := Form{}
	if err := c.BindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}

	outend := hh.ruleEngine.GetOutEnd(form.UUID)
	if outend == nil {
		c.JSON(common.HTTP_OK, common.Error(fmt.Sprintf("'OutEnd' not exists: %v", form.UUID)))
		return
	}
	glogger.GLogger.WithFields(logrus.Fields{
		"topic": "rule/test/" + form.UUID,
	}).Debug(form.TestData)
	err1 := interqueue.DefaultDataCacheQueue.PushOutQueue(outend, form.TestData)
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
* Device
*
 */
func TestDeviceCallback(c *gin.Context, hh *HttpApiServer) {
	type Form struct {
		UUID     string `json:"uuid"`
		TestData string `json:"testData"`
	}
	form := Form{}
	if err := c.BindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	glogger.GLogger.WithFields(logrus.Fields{
		"topic": "rule/test/" + form.UUID,
	}).Debug(form.TestData)
	device := hh.ruleEngine.GetDevice(form.UUID)
	if device == nil {
		c.JSON(common.HTTP_OK, common.Error(fmt.Sprintf("'Device' not exists: %v", form.UUID)))
		return
	}
	err1 := interqueue.DefaultDataCacheQueue.PushDeviceQueue(device, form.TestData)
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}
