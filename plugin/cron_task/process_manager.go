package cron_task

import (
	"encoding/json"
	"errors"
	"github.com/hootrhino/rulex/plugin/http_server/model"
	"io"
	"os/exec"
	"strings"
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

func (pm *ProcessManager) RunProcess(file io.Writer, task model.MScheduleTask) error {
	// 0. arguments
	// 1. working directory
	// 2. environment

	split := strings.Split(task.Args, " ")
	var command *exec.Cmd
	args := make([]string, 0)
	var name string
	if task.TaskType == 1 {
		name = "/bin/bash"
		args = append(args, task.Command)
		args = append(args, split...)
	} else {
		return errors.New("unknown taskType")
	}
	command = exec.Command(name, args...)
	command.SysProcAttr = GetSysProcAttr()
	command.Stderr = file
	command.Stdout = file
	command.Dir = task.WorkDir
	envSlice := make([]string, 0)
	err := json.Unmarshal([]byte(task.Env), &envSlice)
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

func (pm *ProcessManager) KillProcess(id int) error {
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

func (pm *ProcessManager) ListProcess() map[int32]*exec.Cmd {
	m := make(map[int32]*exec.Cmd)
	f := func(key, value any) bool {
		k := key.(int32)
		cmd := value.(*exec.Cmd)
		m[k] = cmd
		return true
	}
	pm.runningProcess.Range(f)
	return m
}
