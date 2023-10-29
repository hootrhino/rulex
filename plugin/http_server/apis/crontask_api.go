package apis

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rulex/component/cron_task"
	"github.com/hootrhino/rulex/component/interdb"
	"github.com/hootrhino/rulex/plugin/http_server/dto"
	"github.com/hootrhino/rulex/plugin/http_server/model"
	"github.com/hootrhino/rulex/plugin/http_server/service"
	"github.com/hootrhino/rulex/typex"
)

// 参考 https://blog.csdn.net/newbieJ/article/details/127125140

// CreateCronTask godoc
// @BasePath /api/v1
// @Summary 创建定时任务
// @Tags crontask
// @param object body dto.CronTaskCreateDTO true "创建"
// @Accept json
// @Produce json
// @Success 200 {object} httpserver.R
// @Router /crontask/create [post]
func CreateCronTask(c *gin.Context, ruleEngine typex.RuleX) (any, error) {
	// 1. 从c中取出数据
	data := dto.CronTaskCreateDTO{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		return nil, err
	}
	// 2. 新增到数据库
	task, err := service.CreateScheduleTask(&data)
	if err != nil {
		return nil, err
	}

	return task.UUID, nil
}

// DeleteCronTask godoc
// @BasePath /api/v1
// @Summary 删除定时任务
// @Tags crontask
// @Param uuid query string true "uuid"
// @Accept json
// @Produce json
// @Success 200 {object} httpserver.R
// @Router /crontask/delete [delete]
func DeleteCronTask(c *gin.Context, ruleEngine typex.RuleX) (any, error) {
	uuid := c.Query("uuid")
	err := service.DeleteScheduleTask(uuid)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// ListCronTask godoc
// @BasePath /api/v1
// @Summary 获取所有定时任务
// @Tags crontask
// @Produce json
// @Success 200 {object} httpserver.R
// @Router /crontask/list [get]
func ListCronTask(c *gin.Context, ruleEngine typex.RuleX) (any, error) {
	condition := model.MCronTask{}
	scheduleTask, err := service.ListScheduleTask(condition)
	return scheduleTask, err
}

// UpdateCronTask godoc
// @BasePath /api/v1
// @Summary 更新定时任务
// @Tags crontask
// @param object body dto.CronTaskUpdateDTO true "更新"
// @Accept json
// @Produce json
// @Success 200 {object} httpserver.R
// @Router /crontask/update [put]
func UpdateCronTask(c *gin.Context, ruleEngine typex.RuleX) (any, error) {
	updateDTO := dto.CronTaskUpdateDTO{}
	err := c.ShouldBind(&updateDTO)
	if err != nil {
		return nil, err
	}

	task, err := service.UpdateScheduleTask(&updateDTO)
	if err != nil {
		return nil, err
	}
	return task, err
}

// StartTask godoc
// @BasePath /api/v1
// @Summary 启动定时任务
// @Tags crontask
// @Param uuid query string true "uuid"
// @Produce json
// @Success 200 {object} httpserver.R
// @Router /crontask/start [get]
func StartTask(c *gin.Context, ruleEngine typex.RuleX) (any, error) {
	uuid, ok := c.GetQuery("uuid")
	if !ok {
		return nil, errors.New("uuid must not be null")
	}
	// 0. 更新数据库
	db := interdb.DB()
	task := model.MCronTask{}
	task.Enable = "1"
	tx := db.Where("uuid = ?", uuid).Updates(&task)
	if tx.Error != nil {
		return nil, tx.Error
	}
	tx = db.Where("uuid = ?", uuid).Find(&task)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if task.UUID == "" {
		return nil, errors.New("cron task not exist")
	}

	// 1. 调用cron的库进行调度
	manager := cron_task.GetCronManager()
	err := manager.AddTask(task)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// StopTask godoc
// @BasePath /api/v1
// @Summary 停止定时任务
// @Tags crontask
// @Param uuid query string true "uuid"
// @Produce json
// @Success 200 {object} httpserver.R
// @Router /crontask/stop [get]
func StopTask(c *gin.Context, ruleEngine typex.RuleX) (any, error) {
	uuid, ok := c.GetQuery("uuid")
	if !ok {
		return nil, errors.New("uuid must not be null")
	}

	// 0. 更新数据库
	db := interdb.DB()
	task := model.MCronTask{}
	task.Enable = "0"
	tx := db.Where("uuid = ?", uuid).Updates(&task)
	if tx.Error != nil {
		return nil, tx.Error
	}

	// 1. 调用cron的库进行调度
	manager := cron_task.GetCronManager()
	manager.DeleteTask(uuid)
	return nil, nil
}

func ListRunningTask(c *gin.Context, ruleEngine typex.RuleX) (any, error) {
	// 1. 请求cronManager获取正在运行的列表
	manager := cron_task.GetCronManager()
	tasks := manager.ListRunningTask()
	return tasks, nil
}

func TerminateRunningTask(c *gin.Context, ruleEngine typex.RuleX) (any, error) {
	// 1. 请求cronManager进行停止正在运行的任务
	uuid, ok := c.GetQuery("uuid")
	if !ok {
		return nil, errors.New("id must not be null")
	}
	manager := cron_task.GetCronManager()
	err := manager.KillTask(uuid)
	if err != nil {
		return nil, err
	}
	return 0, nil
}
