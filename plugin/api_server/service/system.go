package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rulex/plugin/api_server/response"
	"github.com/hootrhino/rulex/typex"
)

type SystemService interface {
	GetPlugins(ctx *gin.Context)
}

type System struct {
	ruleEngine typex.RuleX
}

func (s System) GetPlugins(ctx *gin.Context) {
	var data []interface{}
	plugins := s.ruleEngine.AllPlugins()
	plugins.Range(func(key, value interface{}) bool {
		pi := value.(typex.XPlugin).PluginMetaInfo()
		data = append(data, pi)
		return true
	})
	response.OkWithData(data, ctx)
}

func NewSystemService(ruleEngine typex.RuleX) SystemService {
	return System{ruleEngine}
}
