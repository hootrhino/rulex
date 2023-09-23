package service

import (
	"encoding/json"

	"github.com/hootrhino/rulex/component/interdb"
	"github.com/hootrhino/rulex/plugin/http_server/dto"
	"github.com/hootrhino/rulex/plugin/http_server/model"
)

func CreateScheduleTask(data *dto.CronTaskCreateDTO) (*model.MCronTask, error) {
	db := interdb.DB()
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
	db := interdb.DB()
	task := model.MCronTask{}
	task.ID = id
	tx := db.Delete(&task)
	return tx.Error
}

func PageScheduleTask(page model.PageRequest, task model.MCronTask) (any, error) {
	db := interdb.DB()
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

func UpdateScheduleTask(data *dto.CronTaskUpdateDTO) (*model.MCronTask, error) {
	task := &model.MCronTask{
		Name:     data.Name,
		CronExpr: data.CronExpr,
		TaskType: data.TaskType,
		IsRoot:   data.IsRoot,
	}

	db := interdb.DB()
	t := db.Model(&task)
	if data.Args != nil {
		task.Args = *data.Args
		// 为了能更新为空
		if task.Args == "" {
			task.Args = " "
		}
	}
	task.ID = uint(data.ID)
	if data.Env != nil {
		marshal, _ := json.Marshal(data.Env)
		task.Env = string(marshal)
	}

	tx := t.Updates(task)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return task, nil
}
