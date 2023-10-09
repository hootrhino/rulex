package service

import (
	"encoding/json"
	"errors"
	"github.com/hootrhino/rulex/component/cron_task"
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
		Script:   data.Script,
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

	// 停止已经在调度的任务
	manager := cron_task.GetCronManager()
	manager.DeleteTask(id)
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
	db := interdb.DB()
	cronTask := model.MCronTask{}
	d := db.Model(&model.MCronTask{})
	find := d.Find(&cronTask, data.ID)
	if find.Error != nil {
		return nil, find.Error
	}
	if find.RowsAffected == 0 {
		return nil, errors.New("定时任务不存在")
	}
	if cronTask.Enable == "1" {
		return nil, errors.New("请先暂停任务")
	}
	task := &model.MCronTask{
		Name:     data.Name,
		CronExpr: data.CronExpr,
		TaskType: data.TaskType,
		IsRoot:   data.IsRoot,
		Args:     data.Args,
	}

	t := db.Model(&task)
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
