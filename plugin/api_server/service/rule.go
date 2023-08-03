package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rulex/global"
	"github.com/hootrhino/rulex/plugin/api_server/model"
	"github.com/hootrhino/rulex/plugin/api_server/response"
)

type RuleService interface {
	GetRuleList(ctx *gin.Context)
}

type Rule struct{}

func (d Rule) GetRuleList(ctx *gin.Context) {
	//TODO 业务逻辑
	var rule model.MRule
	global.RULEX_DB.First(&rule, 1)
	response.OkWithData(rule, ctx)
}

func NewDemoService() RuleService {
	return Rule{}
}
