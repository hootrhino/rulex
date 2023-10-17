package trailer

import (
	"context"
	"encoding/json"
	"os/exec"
	"syscall"
)

/*
*
* 子进程的配置, 将 SocketAddr 传入 GRPC 客户端, Args 传入外挂的启动参数
*  $> /test_driver Args
 */
type Goods struct {
	UUID string
	// TCP or Unix Socket
	LocalPath string
	NetAddr   string
	// Description text
	Description string
	// Additional Args
	Args []string
}

//--------------------------------------------------------------------------------------------------
// GoodsProcess
//--------------------------------------------------------------------------------------------------

type GoodsProcess struct {
	Running bool   `json:"running,omitempty"`
	Uuid    string `json:"uuid,omitempty"`
	// 首先启动本地文件，然后用网络路径去发送RPC
	LocalPath   string   `json:"local_path,omitempty"` // 文件路径
	NetAddr     string   `json:"net_addr,omitempty"`   // RPC网络请求路径
	Description string   `json:"description,omitempty"`
	Args        []string `json:"args,omitempty"`
	ctx         context.Context
	cmd         *exec.Cmd
	cancel      context.CancelFunc
}

func (t GoodsProcess) String() string {
	r := map[string]interface{}{
		"running":     t.Running,
		"uuid":        t.Uuid,
		"LocalPath":   t.LocalPath,
		"NetAddr":     t.NetAddr,
		"description": t.Description,
		"args":        t.Args,
	}
	b, _ := json.Marshal(r)
	return string(b)
}
func (scm *GoodsProcess) Stop() {
	if scm.cmd != nil {
		if scm.cmd.Process != nil {
			scm.cancel()
			scm.cmd.Process.Kill()
			scm.cmd.Process.Signal(syscall.SIGKILL)
		}
	}
}

func NewGoodsProcess() *GoodsProcess {
	return &GoodsProcess{}
}
