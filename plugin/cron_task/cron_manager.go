package cron_task

import (
	"fmt"
	"github.com/hootrhino/rulex/glogger"
	sqlitedao "github.com/hootrhino/rulex/plugin/http_server/dao/sqlite"
	"github.com/hootrhino/rulex/plugin/http_server/model"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

var cronManager *CronManager

type CronManager struct {
	cronEngine     *cron.Cron
	crontab        map[uint]cron.EntryID
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
	engine := cron.New(cron.WithSeconds())
	manager := CronManager{
		cronEngine:     engine,
		crontab:        make(map[uint]cron.EntryID),
		processManager: NewProcessManager(),
	}
	// 每天0点10分清理日志
	engine.AddFunc("0 10 0 * * *", func() {
		thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
		// 指定文件夹路径
		folderPath := "cron_logs"
		// 遍历文件夹下的所有文件
		err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// 检查是否为文件夹
			if info.IsDir() {
				// 解析文件夹名称为日期
				dirName := filepath.Base(path)
				date, err := time.Parse("2006-01-02", dirName)
				if err != nil {
					return err
				}

				// 检查是否为 30 天前的文件夹
				if date.Before(thirtyDaysAgo) {
					// 删除文件夹及其内容
					err := os.RemoveAll(path)
					if err != nil {
						return err
					}
					glogger.GLogger.Info("Deleted folder: %s", path)
				}
			}
			return nil
		})
		glogger.GLogger.Error("clean cron logs failed, err=", err)
	})
	engine.Start()
	return &manager
}

func (m *CronManager) AddTask(task model.MCronTask) error {
	cronExpr := task.CronExpr
	id := task.ID
	dir, _ := os.Getwd()
	task.WorkDir = path.Join(dir, task.WorkDir)
	task.Command = path.Join(dir, task.Command)
	entryId, err := m.cronEngine.AddFunc(cronExpr, func() {
		// 打开一个新的logger
		now := time.Now()
		logPath := fmt.Sprintf("cron_logs/%s/%v", now.Format("2006-01-02"), id)
		err := os.MkdirAll(logPath, 0777)
		if err != nil {
			glogger.GLogger.Error(err)
		}
		logTask := logrus.New()
		filePath := path.Join(logPath, now.Format("15-04-05")+".log")
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			glogger.GLogger.Error(err)
			logTask.Out = os.Stdout
		} else {
			logTask.Out = file
		}
		defer func() {
			if file != nil {
				file.Close()
			}
		}()

		//save
		result := model.MCronResult{
			TaskId:    id,
			Status:    "1",
			LogPath:   filePath,
			StartTime: now,
		}
		saveResults(&result)

		logTask.Info("---------------Start task---------------")

		m.runningTask.Store(id, task)
		defer m.runningTask.Delete(id)
		// 调用process manager启动任务并等待其完成
		err = m.processManager.RunProcess(logTask.Out, task)
		exitCode := 0
		if err != nil {
			logTask.Error("Task Return Error, err=", err)
			if exitError, ok := err.(*exec.ExitError); ok {
				// 进程退出时返回非零状态码
				exitCode = exitError.ExitCode()
			} else {
				fmt.Println("执行命令出错:", err)
				exitCode = -1
			}
		}
		logTask.Info("---------------End   task---------------")

		result.EndTime = time.Now()
		result.Status = "2"
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
	db := sqlitedao.Sqlite.DB()
	if m.ID == 0 {
		db.Create(m)
	} else {
		db.Updates(m)
	}
}

func (m *CronManager) DeleteTask(id uint) {
	entryID, ok := m.crontab[id]
	if !ok {
		return
	}
	err := m.processManager.KillProcess(int(id))
	if err != nil {
		glogger.GLogger.Error("kill process failed, err=%+v", err)
	}
	m.cronEngine.Remove(entryID)
	delete(m.crontab, id)
}

func (m *CronManager) KillTask(id int) error {
	return m.processManager.KillProcess(id)
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
