package scheduletask_service

import (
	"encoding/json"
	sqlitedao "github.com/hootrhino/rulex/plugin/http_server/dao/sqlite"
	"github.com/hootrhino/rulex/plugin/http_server/dto"
	"github.com/hootrhino/rulex/plugin/http_server/model"
	"gorm.io/gorm"
)

func CreateScheduleTask(data *dto.CronTaskCreateDTO) (*model.MCronTask, error) {
	db := sqlitedao.Sqlite.DB()
	task := model.MCronTask{
		Name:     data.Name,
		CronExpr: data.CronExpr,
		Enable:   "0",
		TaskType: data.TaskType,
		Args:     data.Args,
		IsRoot:   data.IsRoot,
	}
	if data.Env != nil {
		marshal, _ := json.Marshal(data.Env)
		task.Env = string(marshal)
	} else {
		task.Env = "[]"
	}

	tx := db.Create(&task)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &task, nil
}

func DeleteScheduleTask(id uint) error {
	db := sqlitedao.Sqlite.DB()
	task := model.MCronTask{}
	task.ID = id
	tx := db.Delete(&task)
	return tx.Error
}

func PageScheduleTask(page model.PageRequest, task model.MCronTask) (any, error) {
	db := sqlitedao.Sqlite.DB()
	var records []model.MCronTask
	var count int64
	t := db.Model(&model.MCronTask{}).Where(&model.MCronTask{}, &task).Count(&count)
	if t.Error != nil {
		return nil, t.Error
	}
	tx := db.Scopes(Paginate(page)).Find(&records, &task)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return WrapPageResult(page, records, count), nil
}

func Paginate(page model.PageRequest) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		current := page.Current
		size := page.Size

		offset := (current - 1) * size
		return db.Offset(offset).Limit(size)
	}
}

func WrapPageResult(page model.PageRequest, records any, count int64) model.PageResult {
	return model.PageResult{
		Current: page.Current,
		Size:    page.Size,
		Total:   int(count),
		Records: records,
	}
}
