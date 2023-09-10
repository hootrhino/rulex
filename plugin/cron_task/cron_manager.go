package cron_task

import (
	"github.com/hootrhino/rulex/plugin/http_server/model"
	"github.com/robfig/cron/v3"
	"sync"
)

type CronManager struct {
	cronEngine     *cron.Cron
	m              map[uint]cron.EntryID
	runningTask    sync.Map
	processManager *LinuxProcessManager
}

func NewCronManager() *CronManager {
	engine := cron.New(cron.WithSeconds())
	manager := CronManager{
		cronEngine:     engine,
		m:              make(map[uint]cron.EntryID),
		processManager: NewProcessManager(),
	}
	return &manager
}

func (m *CronManager) AddTask(task model.MScheduleTask) error {
	cronExpr := task.CronExpr
	id := task.ID
	entryId, err := m.cronEngine.AddFunc(cronExpr, func() {
		m.runningTask.Store(id, nil)
		defer m.runningTask.Delete(id)
		// 1. 调用process manager启动进程
		m.processManager.RunProcess(task)
		// 2. 阻塞等待其运行完成

	})
	if err != nil {
		return err
	}
	m.m[id] = entryId
	return nil
}

func (m *CronManager) DeleteTask(id uint) {
	entryID, ok := m.m[id]
	if !ok {
		return
	}
	m.cronEngine.Remove(entryID)
	delete(m.m, id)
}
