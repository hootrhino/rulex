package cron_task

import (
	"fmt"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/plugin/http_server/model"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"sync"
	"time"
)

type CronManager struct {
	cronEngine     *cron.Cron
	m              map[uint]cron.EntryID
	runningTask    sync.Map
	processManager *ProcessManager
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
		// 打开一个新的logger
		now := time.Now()
		now.Format("15:04:05")
		path := fmt.Sprintf("./cron_task/%s/%v/%s.log", now.Format("2006-01-02"), task.ID, now.Format("15-04-05"))
		logTask := logrus.New()
		file, err2 := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err2 != nil {
			glogger.GLogger.Error(err2)
			logTask.Out = os.Stdout
		} else {
			logTask.Out = file
		}
		defer file.Close()
		logTask.Info("---------------Start task---------------")

		m.runningTask.Store(id, nil)
		defer m.runningTask.Delete(id)
		// 调用process manager启动任务并等待其完成
		_, err := m.processManager.RunProcess(logTask.Out, task)
		if err != nil {
			logTask.Info("--")
			return
		}
		logTask.Info("---------------End   task---------------")
		logTask = nil
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
	err := m.processManager.KillProcess(int(id))
	if err != nil {
		log.Default().Printf("kill process failed; %+v", err)
	}
	m.cronEngine.Remove(entryID)
	delete(m.m, id)
}
