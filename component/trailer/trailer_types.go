package trailer

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"syscall"

	"github.com/hootrhino/rulex/glogger"
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
	Args string // such la -al
}

//--------------------------------------------------------------------------------------------------
// GoodsProcess
//--------------------------------------------------------------------------------------------------

type GoodsProcess struct {
	PsRunning bool   `json:"running,omitempty"`
	Name      string `json:"name,omitempty"`
	Uuid      string `json:"uuid,omitempty"`
	Pid       int    `json:"pid,omitempty"`
	// 首先启动本地文件，然后用网络路径去发送RPC
	LocalPath   string `json:"local_path,omitempty"` // 文件路径
	NetAddr     string `json:"net_addr,omitempty"`   // RPC网络请求路径
	Description string `json:"description,omitempty"`
	Args        string `json:"args,omitempty"`
	rpcStarted  bool
	ctx         context.Context
	cmd         *exec.Cmd
	cancel      context.CancelFunc
}

func (t GoodsProcess) String() string {
	b, _ := json.Marshal(t)
	return string(b)
}
func (goodsPs *GoodsProcess) Stop() {
	if goodsPs.cmd != nil {
		if goodsPs.cmd.Process != nil {
			goodsPs.cmd.Process.Kill()
			goodsPs.cmd.Process.Signal(syscall.SIGTERM)
		}
		if goodsPs.cmd.Process != nil {
			pr, err := os.FindProcess(goodsPs.Pid)
			if err != nil {
				glogger.GLogger.Error(err)
				return
			}
			pr.Kill()
		}
	}
	goodsPs.cancel()
}

func NewGoodsProcess() *GoodsProcess {
	return &GoodsProcess{}
}
