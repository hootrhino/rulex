package httpserver

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rulex/plugin/cron_task"
	sqlitedao "github.com/hootrhino/rulex/plugin/http_server/dao/sqlite"
	"github.com/hootrhino/rulex/plugin/http_server/model"
	"github.com/hootrhino/rulex/plugin/http_server/service/scheduletask_service"
	"strconv"
)

// CreateScheduleTask
// 创建定时任务
func CreateScheduleTask(c *gin.Context, hs *HttpApiServer) (any, error) {
	// 1. 从c中取出数据
	dto := model.MScheduleTask{}
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		return nil, err
	}
	// TODO 同时处理文件
	// 2. 新增到数据库
	err = scheduletask_service.CreateScheduleTask(&dto)
	if err != nil {
		return nil, err
	}
	return dto.ID, nil
}

func DeleteScheduleTask(c *gin.Context, hs *HttpApiServer) (any, error) {
	dto := model.MScheduleTask{}
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		return nil, err
	}
	err = scheduletask_service.DeleteScheduleTask(dto.ID)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func PageScheduleTask(c *gin.Context, hs *HttpApiServer) (any, error) {
	page := model.PageRequest{}
	var err error
	page.Current, err = strconv.Atoi(c.DefaultQuery("current", "1"))
	if err != nil {
		return nil, err
	}
	page.Size, err = strconv.Atoi(c.DefaultQuery("size", "25"))
	if err != nil {
		return nil, err
	}

	condition := model.MScheduleTask{}
	scheduleTask, err := scheduletask_service.PageScheduleTask(page, condition)
	return scheduleTask, err
}

func UpdateScheduleTask(c *gin.Context, hs *HttpApiServer) (any, error) {
	// TODO 更新其他信息
	return nil, nil
}

func StartTask(c *gin.Context, hs *HttpApiServer) (any, error) {
	id, ok := c.GetQuery("id")
	if !ok {
		return nil, errors.New("id must not be null")
	}
	idNum, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	// 0. 更新数据库
	db := sqlitedao.Sqlite.DB()
	task := model.MScheduleTask{}
	task.ID = uint(idNum)
	tx := db.Model(&model.MScheduleTask{}).Save(&task)
	if tx.Error != nil {
		return nil, tx.Error
	}
	condition := &model.MScheduleTask{}
	condition.ID = task.ID
	find := db.Where(condition).Find(&task)
	if find.Error != nil {
		return nil, find.Error
	}

	// 1. 调用cron的库进行调度
	manager := cron_task.GetCronManager()
	err = manager.AddTask(task)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func StopTask(c *gin.Context, hs *HttpApiServer) (any, error) {
	id, ok := c.GetQuery("id")
	if !ok {
		return nil, errors.New("id must not be null")
	}
	idNum, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	// 0. 更新数据库
	db := sqlitedao.Sqlite.DB()
	tx := db.Model(&model.MScheduleTask{
		RulexModel: model.RulexModel{
			ID: uint(idNum),
		},
	}).Update("enable", 0)
	if tx.Error != nil {
		return nil, tx.Error
	}

	// 1. 调用cron的库进行调度
	manager := cron_task.GetCronManager()
	manager.DeleteTask(uint(idNum))
	return nil, nil
}

func ListRunningTask(c *gin.Context, hs *HttpApiServer) (any, error) {
	// 1. 请求cronManager获取正在运行的列表
	manager := cron_task.GetCronManager()
	tasks := manager.ListRunningTask()
	return tasks, nil
}

func TerminateRunningTask(c *gin.Context, hs *HttpApiServer) (any, error) {
	// 1. 请求cronManager进行停止正在运行的任务
	id, ok := c.GetQuery("id")
	if !ok {
		return nil, errors.New("id must not be null")
	}
	idNum, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	manager := cron_task.GetCronManager()
	err = manager.KillTask(idNum)
	return nil, nil
}
