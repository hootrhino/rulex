package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rulex/component/interdb"
	common "github.com/hootrhino/rulex/component/rulex_api_server/common"
	"github.com/hootrhino/rulex/component/rulex_api_server/dto"
	"github.com/hootrhino/rulex/component/rulex_api_server/model"
	"github.com/hootrhino/rulex/component/rulex_api_server/service"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
	"gorm.io/gorm"
)

type UserLuaTemplateVo struct {
	Gid       string                     `json:"gid,omitempty"`  // 分组ID
	UUID      string                     `json:"uuid,omitempty"` // 名称
	Label     string                     `json:"label"`          // 快捷代码名称
	Apply     string                     `json:"apply"`          // 快捷代码
	Type      string                     `json:"type"`           // 类型 固定为function类型detail
	Detail    string                     `json:"detail"`         // 细节
	Variables []dto.LuaTemplateVariables `json:"variables"`      // 变量
}

/*
*
* 新建用户模板
*
 */

func CreateUserLuaTemplate(c *gin.Context, ruleEngine typex.RuleX) {
	form := UserLuaTemplateVo{Type: "function"}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	_, err0 := service.GetGenericGroupWithUUID(form.Gid)
	if err0 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err0))
		return
	}
	MUserLuaTemplate := model.MUserLuaTemplate{
		UUID:   utils.UserLuaUuid(),
		Label:  form.Label,
		Type:   "function",
		Apply:  form.Apply,
		Detail: form.Detail,
		Gid:    form.Gid,
	}
	if err := service.InsertUserLuaTemplate(MUserLuaTemplate); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 新建用户模板的时候必须给一个分组
	if err := service.BindResource(form.Gid, MUserLuaTemplate.UUID); err != nil {
		c.JSON(common.HTTP_OK, common.Error("Group not found"))
		return
	}
	// 返回新建的用户模板字段 用来跳转编辑器
	c.JSON(common.HTTP_OK, common.OkWithData(map[string]string{
		"uuid": MUserLuaTemplate.UUID,
	}))

}

/*
*
* 更新用户模板
*
 */
func UpdateUserLuaTemplate(c *gin.Context, ruleEngine typex.RuleX) {
	form := UserLuaTemplateVo{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	MUserLuaTemplate := model.MUserLuaTemplate{
		UUID:   form.UUID,
		Label:  form.Label,
		Type:   form.Type,
		Apply:  form.Apply,
		Detail: form.Detail,
		Gid:    form.Gid,
	}
	Variables, errVariables := MUserLuaTemplate.GenVariables(form.Variables)
	if errVariables != nil {
		c.JSON(common.HTTP_OK, common.Error400(errVariables))
		return
	}
	MUserLuaTemplate.Variables = Variables
	// 事务
	txErr := service.ReBindResource(func(tx *gorm.DB) error {
		return tx.Model(MUserLuaTemplate).
			Where("uuid=?", MUserLuaTemplate.UUID).
			Updates(&MUserLuaTemplate).Error
	}, form.UUID, form.Gid)
	if txErr != nil {
		c.JSON(common.HTTP_OK, common.Error400(txErr))
		return
	}
	// 返回新建的用户模板字段 用来跳转编辑器
	c.JSON(common.HTTP_OK, common.OkWithData(map[string]string{
		"uuid": MUserLuaTemplate.UUID,
	}))
}

/*
*
* 删除用户模板
*
 */
func DeleteUserLuaTemplate(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	txErr := interdb.DB().Transaction(func(tx *gorm.DB) error {
		Group := service.GetVisualGroup(uuid)
		if err1 := service.DeleteUserLuaTemplate(uuid); err1 != nil {
			c.JSON(common.HTTP_OK, common.Error400(err1))
			return err1
		}
		// 解除关联
		err2 := interdb.DB().Where("gid=? and rid =?", Group.UUID, uuid).
			Delete(&model.MGenericGroupRelation{}).Error
		if err2 != nil {
			return err2
		}
		return nil
	})
	if txErr != nil {
		c.JSON(common.HTTP_OK, common.Error400(txErr))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())

}
func ListUserLuaTemplateGroup(c *gin.Context, ruleEngine typex.RuleX) {
	visuals := []MGenericGroupVo{}
	for _, vv := range service.ListByGroupType("USER_LUA_TEMPLATE") {
		visuals = append(visuals, MGenericGroupVo{
			UUID:   vv.UUID,
			Name:   vv.Name,
			Type:   vv.Type,
			Parent: vv.Parent,
		})
	}
	c.JSON(common.HTTP_OK, common.OkWithData(visuals))
}

/*
*
* 模糊查询
*
 */
func SearchUserLuaTemplateGroup(c *gin.Context, ruleEngine typex.RuleX) {
	visuals := []UserLuaTemplateVo{}
	keyword, _ := c.GetQuery("keyword")
	for _, vv := range service.SearchUserLuaTemplate(keyword, keyword) {
		visuals = append(visuals, UserLuaTemplateVo{
			UUID:   vv.UUID,
			Label:  vv.Label,
			Type:   vv.Type,
			Apply:  vv.Apply,
			Detail: vv.Detail,
			Gid:    vv.Gid,
		})
	}
	c.JSON(common.HTTP_OK, common.OkWithData(visuals))
}

/*
*
* 用户模板列表
*
 */
func ListUserLuaTemplate(c *gin.Context, ruleEngine typex.RuleX) {
	UserLuaTemplates := []UserLuaTemplateVo{}
	for _, vv := range service.AllUserLuaTemplate() {
		Vo := UserLuaTemplateVo{
			UUID:      vv.UUID,
			Label:     vv.Label,
			Type:      vv.Type,
			Apply:     vv.Apply,
			Detail:    vv.Detail,
			Gid:       vv.Gid,
			Variables: vv.GetVariables(),
		}
		Group := service.GetUserLuaTemplateGroup(vv.UUID)
		if Group.UUID != "" {
			Vo.Gid = Group.UUID
		} else {
			Vo.Gid = ""
		}
		UserLuaTemplates = append(UserLuaTemplates, Vo)
	}
	c.JSON(common.HTTP_OK, common.OkWithData(UserLuaTemplates))

}

/*
*
* 用户模板分组查看
*
 */
func ListUserLuaTemplateByGroup(c *gin.Context, ruleEngine typex.RuleX) {
	Gid, _ := c.GetQuery("uuid")
	UserLuaTemplates := []UserLuaTemplateVo{}
	MUserLuaTemplates := service.FindUserTemplateByGroup(Gid)
	for _, vv := range MUserLuaTemplates {
		Vo := UserLuaTemplateVo{
			UUID:      vv.UUID,
			Label:     vv.Label,
			Type:      vv.Type,
			Apply:     vv.Apply,
			Detail:    vv.Detail,
			Gid:       vv.Gid,
			Variables: vv.GetVariables(),
		}
		Group := service.GetUserLuaTemplateGroup(vv.UUID)
		Vo.Gid = Group.UUID
		UserLuaTemplates = append(UserLuaTemplates, Vo)
	}
	c.JSON(common.HTTP_OK, common.OkWithData(UserLuaTemplates))
}

/*
*
* 用户模板详情
*
 */
func UserLuaTemplateDetail(c *gin.Context, ruleEngine typex.RuleX) {
	uuid, _ := c.GetQuery("uuid")
	mUserLuaTemplate, err := service.GetUserLuaTemplateWithUUID(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	Vo := UserLuaTemplateVo{
		UUID:      mUserLuaTemplate.UUID,
		Label:     mUserLuaTemplate.Label,
		Type:      mUserLuaTemplate.Type,
		Apply:     mUserLuaTemplate.Apply,
		Detail:    mUserLuaTemplate.Detail,
		Variables: mUserLuaTemplate.GetVariables(),
	}
	Group := service.GetUserLuaTemplateGroup(mUserLuaTemplate.UUID)
	if Group.UUID != "" {
		Vo.Gid = Group.UUID
	} else {
		Vo.Gid = ""
	}
	c.JSON(common.HTTP_OK, common.OkWithData(Vo))
}
