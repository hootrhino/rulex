package apis

import (
	"github.com/gin-gonic/gin"
	sqlitedao "github.com/hootrhino/rulex/plugin/http_server/dao/sqlite"
	"github.com/hootrhino/rulex/plugin/http_server/model"
	"github.com/hootrhino/rulex/plugin/http_server/service"
	"github.com/hootrhino/rulex/typex"
	"strconv"
)

func PageCronTaskResult(c *gin.Context, ruleEngine typex.RuleX) (any, error) {
	page, err := service.ReadPageRequest(c)
	if err != nil {
		return nil, err
	}

	cronResult := model.MCronResult{}
	taskId := c.Query("taskId")
	if taskId != "" {
		atoi, err := strconv.Atoi(taskId)
		if err != nil {
			return nil, err
		}
		cronResult.TaskId = uint(atoi)
	}

	db := sqlitedao.Sqlite.DB()
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
