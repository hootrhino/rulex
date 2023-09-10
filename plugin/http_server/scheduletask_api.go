package httpserver

import (
	"github.com/gin-gonic/gin"
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
	// TODO
	return nil, nil
}

func StartTask(c *gin.Context, hs *HttpApiServer) (any, error) {
	// TODO
	// 0. 更新数据库
	// 1. 调用cron的库进行调度
	return nil, nil
}

func StopTask(c *gin.Context, hs *HttpApiServer) (any, error) {
	// TODO
	// 0. 更新数据库
	// 1. 调用cron的库进行调度
	return nil, nil
}

func ListRunningTask(c *gin.Context, hs *HttpApiServer) (any, error) {
	// TODO
	// 1. 请求ProcessManager获取正在运行的列表
	return nil, nil
}

func TerminateRunningTask(c *gin.Context, hs *HttpApiServer) (any, error) {
	// TODO
	// 1. 请求ProcessManager进行停止正在运行的任务
	return nil, nil
}
