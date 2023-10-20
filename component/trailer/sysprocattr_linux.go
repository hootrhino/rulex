package trailer

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

func NewSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		Setsid: false,
	}
}

type ProcessInfo struct {
	CPU     string `json:"cpu"`
	Memory  string `json:"memory"`
	VSZ     string `json:"vsz"`
	RSS     string `json:"rss"`
	TTY     string `json:"tty"`
	Stat    string `json:"stat"`
	Start   string `json:"start"`
	Time    string `json:"time"`
	Command string `json:"command"`
}

// `ps` 命令的输出中各列的含义如下：
// 1. `%CPU`：进程的 CPU 使用率百分比。它表示进程占用 CPU 时间的百分比。
// 2. `%MEM`：进程的内存使用率百分比。它表示进程占用系统内存的百分比。
// 3. `VSZ`：进程的虚拟内存大小（以KB为单位）。它包括了进程占用的实际内存和虚拟内存，虚拟内存是物理内存和交换空间的组合。
// 4. `RSS`：进程的常驻内存大小（以KB为单位）。它表示实际分配给进程的物理内存。
// 5. `TT`：进程的终端类型，指示与进程关联的终端。
// 6. `STAT`：进程的状态，表示进程当前的状态。常见的状态包括 R（运行）、S（睡眠）、Z（僵尸）、T（暂停）等。
// 7. `STARTED`：进程的启动时间，表示进程启动的日期和时间。
// 8. `TIME`：进程的累积 CPU 时间，表示进程自启动以来已经占用的 CPU 时间。
// 9. `COMMAND`：进程的命令行，显示了启动进程的命令和参数。
// 这些列提供了有关正在运行的进程的重要信息，包括它们的 CPU 和内存使用情况、状态以及启动时间等。在分析和监视进程时，这些信息非常有用。
func RunningProcessDetail(pid int) (ProcessInfo, error) {
	// 创建一个 ProcessInfo 结构体实例
	var processInfo ProcessInfo

	// 构建命令
	cmd := exec.Command("ps", "-p", fmt.Sprint(pid), "-o", "%cpu,%mem,vsz,rss,tty,stat,start,time,command")

	// 执行命令并捕获输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		return processInfo, err
	}

	// 解析输出并填充结构体
	outputLines := strings.Split(string(output), "\n")
	if len(outputLines) > 1 {
		fields := strings.Fields(outputLines[1])
		if len(fields) >= 9 {
			processInfo.CPU = fields[0]
			processInfo.Memory = fields[1]
			processInfo.VSZ = fields[2]
			processInfo.RSS = fields[3]
			processInfo.TTY = fields[4]
			processInfo.Stat = fields[5]
			processInfo.Start = fields[6]
			processInfo.Time = fields[7]
			processInfo.Command = strings.Join(fields[8:], " ")
		}
	}

	return processInfo, nil
}
