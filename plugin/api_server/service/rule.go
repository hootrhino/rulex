package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rulex/global"
	"github.com/hootrhino/rulex/plugin/api_server/model"
	"github.com/hootrhino/rulex/plugin/api_server/response"
	"github.com/hootrhino/rulex/typex"
)

type RuleService interface {
	GetRuleList(ctx *gin.Context)
}

type Rule struct {
	ruleEngine typex.RuleX
}

func (r Rule) GetRuleList(ctx *gin.Context) {
	//TODO 业务逻辑
	var rule model.MRule
	global.RULEX_DB.First(&rule, 1)
	response.OkWithData(rule, ctx)
}

func NewRuleService(ruleEngine typex.RuleX) RuleService {
	return Rule{ruleEngine}
}
