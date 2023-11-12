package service

import (
	"encoding/json"
	"errors"
	"github.com/hootrhino/rulex/component/cron_task"
	"github.com/hootrhino/rulex/component/interdb"
	"github.com/hootrhino/rulex/plugin/http_server/dto"
	"github.com/hootrhino/rulex/plugin/http_server/model"
	"github.com/hootrhino/rulex/utils"
)

func CreateScheduleTask(data *dto.CronTaskCreateDTO) (*model.MCronTask, error) {
	db := interdb.DB()
	task := model.MCronTask{
		UUID:     utils.CronTaskUuid(),
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

func DeleteScheduleTask(uuid string) error {
	db := interdb.DB()
	task := model.MCronTask{}
	tx := db.Where("uuid = ?", uuid).Delete(&task)

	// 停止已经在调度的任务
	manager := cron_task.GetCronManager()
	manager.DeleteTask(uuid)
	return tx.Error
}

func ListScheduleTask(task model.MCronTask) (any, error) {
	db := interdb.DB()
	var records []model.MCronTask
	tx := db.Find(&records)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return records, nil
}

func UpdateScheduleTask(data *dto.CronTaskUpdateDTO) (*model.MCronTask, error) {
	db := interdb.DB()
	cronTask := model.MCronTask{}
	d := db.Model(&model.MCronTask{})
	tx := d.Where("uuid = ?", data.UUID).Find(&cronTask)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
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

	tx = db.Model(&task)
	if data.Env != nil {
		marshal, _ := json.Marshal(data.Env)
		task.Env = string(marshal)
	}

	tx = tx.Where("uuid = ?", data.UUID).Updates(task)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return task, nil
}
