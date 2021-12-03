package httpserver

import (
	"errors"
	"net/http"
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
		c.JSON(http.StatusOK, Result{
			Code: http.StatusOK,
			Msg:  "Success",
			Data: data,
		})
	} else {
		c.JSON(http.StatusOK, Result{
			Code: http.StatusOK,
			Msg:  "Success",
			Data: e.GetRule(uuid),
		})
	}
}

//
// Create rule
//
func CreateRule(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	type Form struct {
		From        []string `json:"from" binding:"required"`
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

	if len(form.From) > 0 {
		for _, id := range form.From {
			if e.GetInEnd(id).UUID == "" {
				c.JSON(200, errors.New(`"inend not exists:" `+id))
				return
			}
		}

		tmpRule := typex.NewRule(nil,
			"",
			form.Name,
			form.Description,
			nil,
			form.Success,
			form.Actions,
			form.Failed)

		if err := core.VerifyCallback(tmpRule); err != nil {
			c.JSON(200, Error400(err))
			return
		} else {
			mRule := &MRule{
				UUID:        utils.MakeUUID("RULE"),
				Name:        form.Name,
				Description: form.Description,
				From:        form.From,
				Success:     form.Success,
				Failed:      form.Failed,
				Actions:     form.Actions,
			}
			if err := hh.InsertMRule(mRule); err != nil {
				c.JSON(200, gin.H{"msg": err.Error()})
				return
			}
			rule := typex.NewRule(hh.ruleEngine,
				mRule.UUID,
				mRule.Name,
				mRule.Description,
				mRule.From,
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
	} else {
		c.JSON(200, Error400(errors.New("from can't empty")))
		return
	}

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
		From    []string `json:"from" binding:"required"`
		Actions string   `json:"actions" binding:"required"`
		Success string   `json:"success" binding:"required"`
		Failed  string   `json:"failed" binding:"required"`
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
		form.From,
		form.Success,
		form.Actions,
		form.Failed)
	if err := core.VerifyCallback(tmpRule); err != nil {
		c.JSON(200, Error400(err))
	} else {
		c.JSON(200, Ok())
	}

}
