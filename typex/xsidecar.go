package typex

import (
	"context"
	"encoding/json"
	"os/exec"
)

/*
*
* 子进程的配置, 将 SocketAddr 传入 GRPC 客户端, Args 传入外挂的启动参数
*  $> /test_driver Args
 */
type Goods struct {
	UUID string
	// TCP or Unix Socket
	Addr string
	// Description text
	Description string
	// Additional Args
	Args []string
}

//--------------------------------------------------------------------------------------------------
// GoodsProcess
//--------------------------------------------------------------------------------------------------

type GoodsProcess struct {
	Running     bool
	Uuid        string
	Addr        string
	Description string
	Args        []string
	Ctx         context.Context
	Cmd         *exec.Cmd
	Cancel      context.CancelFunc
}

func (t GoodsProcess) String() string {
	r := map[string]interface{}{
		"running":     t.Running,
		"uuid":        t.Uuid,
		"addr":        t.Addr,
		"description": t.Description,
		"args":        t.Args,
	}
	b, _ := json.Marshal(r)
	return string(b)
}
func (scm *GoodsProcess) Stop() {
	if scm.Cmd != nil {
		if scm.Cmd.Process != nil {
			scm.Cmd.Process.Kill()
			scm.Cmd = nil
			scm.Cancel()
		}
	}
}

func NewGoodsProcess() *GoodsProcess {
	return &GoodsProcess{}
}
