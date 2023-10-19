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
	UUID      string
	AutoStart *bool
	// TCP or Unix Socket
	LocalPath string
	NetAddr   string
	// Description text
	Description string
	// Additional Args
	Args string // 使用空格分割 , such: la -al
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
	killedBy    string   // 被谁干死的
	mailBox     chan int // 这里用来接收外部控制信号
}

func (t GoodsProcess) String() string {
	b, _ := json.Marshal(t)
	return string(b)
}
func (goodsPs *GoodsProcess) StopBy(r string) {
	goodsPs.killedBy = r
	if goodsPs.cmd != nil {
		if goodsPs.cmd.Process != nil {
			goodsPs.cmd.Process.Kill()
			goodsPs.cmd.Process.Signal(syscall.SIGTERM)
		}
	}
	goodsPs.cancel()
}
func (goodsPs *GoodsProcess) Stop() {
	if goodsPs.cmd != nil {
		if goodsPs.cmd.Process != nil {
			goodsPs.cmd.Process.Kill()
			goodsPs.cmd.Process.Signal(syscall.SIGTERM)
		}
	}
	goodsPs.cancel()
}

func NewGoodsProcess() *GoodsProcess {
	return &GoodsProcess{}
}
