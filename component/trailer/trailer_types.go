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
	UUID        string
	Name        string // 进程名
	GoodsType   string // LOCAL, EXTERNAL
	ExecuteType string // exe,elf,jar,py, nodejs....
	AutoStart   *bool  // 是否开启自启动，目前全部是自启动
	LocalPath   string // TCP or Unix Socket
	NetAddr     string // RPC addr
	Description string // Description text
	// Additional Args
	Args string // 使用空格分割 , such: la -al
}

//--------------------------------------------------------------------------------------------------
// GoodsProcess: 内存里的进程实例
//--------------------------------------------------------------------------------------------------

type GoodsProcess struct {
	Info          GoodsInfo
	psRunning     bool               // 本地进程是否启动了
	pid           int                // pid
	trailerClient TrailerClient      // Grpc客户端
	rpcStarted    bool               // RPC 网络服务是否开启
	ctx           context.Context    // Context
	cmd           *exec.Cmd          // Cmd
	cancel        context.CancelFunc // Cancel Func
	killedBy      string             // 进程被谁干死的, 一般用来处理要不要抢救进程
	mailBox       chan int           // 这里用来接收外部控制信号
}

func (goodsProcess *GoodsProcess) PsRunning() bool {
	return goodsProcess.psRunning

}
func (goodsProcess *GoodsProcess) Pid() int {
	return goodsProcess.pid
}

/*
*
* 尝试新建一个RPC客户端
*
 */
func (goodsProcess *GoodsProcess) ConnectToRpc() (TrailerClient, error) {
	if goodsProcess.trailerClient == nil {
		grpcConnection, err := grpc.Dial(goodsProcess.Info.NetAddr,
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
