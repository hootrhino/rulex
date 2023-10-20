package trailer

import (
	"context"
	"encoding/json"
	"os/exec"
	"syscall"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/*
*
* 子进程的配置, 将 SocketAddr 传入 GRPC 客户端, Args 传入外挂的启动参数
*  $> /test_driver Args
 */
type GoodsInfo struct {
	UUID      string
	Name      string // 进程名
	GoodsType string // LOCAL, EXTERNAL
	AutoStart *bool
	// TCP or Unix Socket
	LocalPath string
	// RPC addr
	NetAddr string
	// Description text
	Description string
	// Additional Args
	Args string // 使用空格分割 , such: la -al
}

//--------------------------------------------------------------------------------------------------
// GoodsProcess: 内存里的进程实例
//--------------------------------------------------------------------------------------------------

type GoodsProcess struct {
	Uuid          string        `json:"uuid"`       // UUID
	Pid           int           `json:"pid"`        // pid
	PsRunning     bool          `json:"running"`    // 本地进程是否启动了
	Name          string        `json:"name"`       // 进程名
	GoodsType     string        `json:"goodsType"`  // LOCAL, EXTERNAL
	LocalPath     string        `json:"local_path"` // 文件保存路径
	NetAddr       string        `json:"net_addr"`   // RPC网络请求路径
	Args          string        `json:"args"`       // 进程参数
	Description   string        `json:"description"`
	trailerClient TrailerClient // Grpc客户端
	rpcStarted    bool          // RPC 网络服务是否开启
	ctx           context.Context
	cmd           *exec.Cmd
	cancel        context.CancelFunc
	killedBy      string   // 被谁干死的, 一般用来处理要不要抢救进程
	mailBox       chan int // 这里用来接收外部控制信号
}

/*
*
* 尝试新建一个RPC客户端
*
 */
func (goodsProcess *GoodsProcess) ConnectToRpc() (TrailerClient, error) {
	if goodsProcess.trailerClient == nil {
		grpcConnection, err := grpc.Dial(goodsProcess.NetAddr,
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, err
		}
		goodsProcess.trailerClient = NewTrailerClient(grpcConnection)
	}
	return goodsProcess.trailerClient, nil
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
