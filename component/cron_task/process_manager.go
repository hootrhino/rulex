package cron_task

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/plugin/http_server/model"
	"io"
	"os/exec"
	"sync"
)

// ProcessManager
type ProcessManager struct {
	runningProcess sync.Map
}

func NewProcessManager() *ProcessManager {
	manager := ProcessManager{
		runningProcess: sync.Map{},
	}
	return &manager
}

func (pm *ProcessManager) RunProcess(file io.Writer, task model.MCronTask) error {
	// 0. arguments
	// 1. working directory
	// 2. environment

	var command *exec.Cmd
	args := make([]string, 0)
	var name string
	script, err := base64.StdEncoding.DecodeString(task.Script)
	if err != nil {
		glogger.GLogger.Error("parse script failed", err)
		return err
	}
	if task.TaskType == CRON_TASK_TYPE_LINUX_SHELL {
		name = "/bin/bash"
		args = append(args, "-c")
		args = append(args, string(script))
		args = append(args, "bash") // 占据$0位置
		if task.Args != nil {
			args = append(args, *task.Args) // 占据$1位置，此时$@与$1相同
		}

	} else {
		return errors.New("unknown taskType")
	}
	command = exec.Command(name, args...)
	command.SysProcAttr = GetSysProcAttr()
	command.Stderr = file
	command.Stdout = file
	command.Dir = task.WorkDir
	envSlice := make([]string, 0)
	err = json.Unmarshal([]byte(task.Env), &envSlice)
	if err != nil {
		return err
	}
	command.Env = envSlice

	pm.runningProcess.Store(task.ID, command)
	defer pm.runningProcess.Delete(task.ID)

	err = command.Run()
	if err != nil {
		return err
	}

	return nil
}

func (pm *ProcessManager) KillProcess(id string) error {
	value, ok := pm.runningProcess.Load(id)
	if !ok {
		// not exist, return success
		return nil
	}
	cmd := value.(*exec.Cmd)
	err := KillProcess(cmd)
	if err != nil {
		return err
	}
	return nil
}

func (pm *ProcessManager) ListProcess() map[string]*exec.Cmd {
	m := make(map[string]*exec.Cmd)
	f := func(key, value any) bool {
		k := key.(string)
		cmd := value.(*exec.Cmd)
		m[k] = cmd
		return true
	}
	pm.runningProcess.Range(f)
	return m
}
