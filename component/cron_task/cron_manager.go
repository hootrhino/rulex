package cron_task

import (
	"os"
	"os/exec"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/hootrhino/rulex/component/interdb"
	"github.com/hootrhino/rulex/component/rulex_api_server/model"
	"github.com/hootrhino/rulex/glogger"
	"github.com/robfig/cron/v3"
)

var cronManager *CronManager

type CronManager struct {
	cronEngine     *cron.Cron
	mtx            sync.Mutex
	crontab        map[string]cron.EntryID
	runningTask    sync.Map
	processManager *ProcessManager
}

func GetCronManager() *CronManager {
	if cronManager == nil {
		cronManager = NewCronManager()
	}
	return cronManager
}

func NewCronManager() *CronManager {
	engine := cron.New(
		cron.WithChain(
			cron.SkipIfStillRunning(cron.DefaultLogger),
			cron.Recover(cron.DefaultLogger),
		),
		cron.WithSeconds(),
	)
	manager := CronManager{
		cronEngine:     engine,
		crontab:        make(map[string]cron.EntryID),
		processManager: NewProcessManager(),
		mtx:            sync.Mutex{},
	}
	engine.Start()
	return &manager
}

func (m *CronManager) AddTask(task model.MCronTask) error {
	cronExpr := task.CronExpr
	id := task.UUID

	m.mtx.Lock()
	defer m.mtx.Unlock()
	if _, ok := m.crontab[id]; ok {
		return nil
	}
	dir, _ := os.Getwd()
	task.WorkDir = path.Join(".")
	err := os.MkdirAll(dir, PERM_0777)
	if err != nil {
		return err
	}
	entryId, err := m.cronEngine.AddFunc(cronExpr, func() {
		// create a task execute log record
		result := model.MCronResult{
			TaskUuid:  task.UUID,
			Status:    CRON_RESULT_STATUS_RUNNING,
			StartTime: time.Now(),
		}
		saveResults(&result)

		taskLogger := glogger.GLogger.WithField("task_uuid", task.UUID)
		taskLogger.Info("---------------Start task---------------")

		m.runningTask.Store(id, task)
		defer m.runningTask.Delete(id)
		// 调用process manager启动任务并等待其完成
		err = m.processManager.RunProcess(taskLogger.Logger.Out, task)
		exitCode := 0
		if err != nil {
			taskLogger.Error("Task Return Error, err=", err)
			if exitError, ok := err.(*exec.ExitError); ok {
				// 进程退出时返回非零状态码
				exitCode = exitError.ExitCode()
			} else {
				exitCode = -1
			}
		}
		taskLogger.Info("---------------End   task---------------")

		result.EndTime = time.Now()
		result.Status = CRON_RESULT_STATUS_END
		result.ExitCode = strconv.Itoa(exitCode)
		saveResults(&result)
	})
	if err != nil {
		return err
	}
	m.crontab[id] = entryId
	return nil
}

func saveResults(m *model.MCronResult) {
	db := interdb.DB()
	if m.ID == 0 {
		db.Create(m)
	} else {
		db.Updates(m)
	}
}

func (m *CronManager) DeleteTask(id string) {
	entryID, ok := m.crontab[id]
	if !ok {
		return
	}
	err := m.processManager.KillProcess(id)
	if err != nil {
		glogger.GLogger.Error("kill process failed, err=%+v", err)
		return
	}
	m.cronEngine.Remove(entryID)
	delete(m.crontab, id)
}

func (m *CronManager) KillTask(uuid string) error {
	return m.processManager.KillProcess(uuid)
}

func (m *CronManager) ListRunningTask() []model.MCronTask {
	tasks := make([]model.MCronTask, 0)
	m.runningTask.Range(func(key, value any) bool {
		task := value.(model.MCronTask)
		tasks = append(tasks, task)
		return true
	})
	return tasks
}
