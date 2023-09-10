package cron_task

import (
	"errors"
	"github.com/hootrhino/rulex/plugin/http_server/model"
	"log"
	"os/exec"
	"strings"
	"sync"
)

// LinuxProcessManager
// Linux进程管理器
type LinuxProcessManager struct {
	runningProcess sync.Map
}

func NewProcessManager() *LinuxProcessManager {
	manager := LinuxProcessManager{
		runningProcess: sync.Map{},
	}
	return &manager
}

func (pm *LinuxProcessManager) RunProcess(task model.MScheduleTask) (int32, error) {
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
		return 0, errors.New("unknown taskType")
	}
	command = exec.Command(name, args...)

	// if linux
	//command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	// TODO 设置指定的logger文件
	command.Stderr = log.Default().Writer()
	command.Stdout = log.Default().Writer()
	command.Dir = task.WorkDir
	command.Env = strings.Split(task.Env, ";")

	err := command.Start()
	if err != nil {
		return 0, err
	}
	pm.runningProcess.Store(task.ID, command)

	err = command.Wait()
	if err != nil {
		log.Println("process return error, error=", err.Error())
		// TODO update DB
	} else {
		// TODO update DB
	}
	return 0, err
}

func (pm *LinuxProcessManager) KillProcess(id int) error {
	value, ok := pm.runningProcess.Load(id)
	if !ok {
		// no exist, return success
		return nil
	}
	cmd := value.(*exec.Cmd)
	err := cmd.Process.Kill()
	if err != nil {
		return err
	}
	pm.runningProcess.Delete(id)
	return nil
}

func (pm *LinuxProcessManager) ListProcess() map[int32]*exec.Cmd {
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
