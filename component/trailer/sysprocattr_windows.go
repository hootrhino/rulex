package trailer

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

func NewSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		HideWindow: true,
	}
}

/*
*
* 进程管理器拿到的进程详情
*
 */
// 定义一个结构体来表示进程的详细信息
type ProcessInfo struct {
	ImageName   string
	PID         string
	SessionName string
	SessionNum  string
}

// C:\Users\wangwenhai>tasklist /v /fi "PID eq 2000" /fo list

// 映像名称:     Code.exe
// PID:          2000
// 会话名      : Console
// 会话#   :     1
// 内存使用 :    247,624 K
// 状态  :       Unknown
// 用户名   :    DESKTOP-EMD3M3C\wangwenhai
// CPU 时间:     0:00:41
// 窗口标题    : 暂缺

func RunningProcessDetail(pid int) (ProcessInfo, error) {
	var processInfo ProcessInfo
	cmd := exec.Command("tasklist.exe", "/fi", fmt.Sprintf("PID eq %d", pid), "/fo", "csv", "/nh")
	fmt.Println(cmd.String())
	// 这里给我造成了困惑, 在windows下一直有个“Error: exec: already started”异常
	// 经过万分艰难才发现原来windows下Command就具备了执行文件的作用，而且貌似这个 tasklist 还会操作stdout
	// 反正这里有点另类，区别对待即可
	output, err := cmd.Output()
	if err != nil {
		return processInfo, err
	}
	outputLines := strings.Split(string(output), "\n")
	if len(outputLines) > 0 {
		fields := strings.Split(outputLines[0], ",")
		if len(fields) >= 6 {
			processInfo.ImageName = strings.Trim(fields[0], "\"")
			processInfo.PID = strings.Trim(fields[1], "\"")
			processInfo.SessionName = strings.Trim(fields[3], "\"")
		}
	}
	return processInfo, nil
}
