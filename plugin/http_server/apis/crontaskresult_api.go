package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rulex/component/interdb"
	"github.com/hootrhino/rulex/plugin/http_server/model"
	"github.com/hootrhino/rulex/plugin/http_server/service"
	"github.com/hootrhino/rulex/typex"
)

// PageCronTaskResult godoc
// @BasePath /api/v1
// @Summary 分页获取定时任务执行结果
// @Tags crontask
// @Param current query string false "current"
// @Param size query string false "size"
// @Param uuid query string false "uuid"
// @Accept json
// @Produce json
// @Success 200 {object} httpserver.R
// @Router /crontask/results/page [get]
func PageCronTaskResult(c *gin.Context, ruleEngine typex.RuleX) (any, error) {
	page, err := service.ReadPageRequest(c)
	if err != nil {
		return nil, err
	}

	cronResult := model.MCronResult{}
	uuid := c.Query("uuid")
	if uuid != "" {
		cronResult.TaskUuid = uuid
	}

	db := interdb.DB()
	var count int64
	t := db.Model(&model.MCronResult{}).Where(&model.MCronResult{}, &cronResult).Count(&count)
	if t.Error != nil {
		return nil, t.Error
	}

	tx := db.Scopes(service.Paginate(*page))
	var records []model.MCronResult
	result := tx.Order("created_at DESC").Find(&records, &cronResult)
	if result.Error != nil {
		return nil, result.Error
	}
	return service.WrapPageResult(*page, records, count), nil
}
