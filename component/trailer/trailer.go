package trailer

//
// Trailer就是拖车,带着小车一起跑,比喻了Trailer实际上是个进程管理器
//                                ____
//      ______         ______     |   \
//     /|_||_\`.__    /|_||_\`.__ | |_ \---
//    (   _    _ _\  (   _    _ _\| | | ___\
//    =`-(_)--(_)-'  =`-(_)--(_)-'|_________\_
//  ______________________________|    |_o__ |
//  |            |[] ___ \_______|   / ___ \|
//  |_____________[]_/.-.\_\______|__/_/.-.\_[]
//                    (O)               (O)
// ---  ---   ---   ---   ---   ---   ------------
import (
	"context"
	"strings"
	"time"

	"os"
	"os/exec"
	"sync"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var __DefaultTrailerRuntime *TrailerRuntime

//--------------------------------------------------------------------------------------------------
// Trailer 接口
//--------------------------------------------------------------------------------------------------

type TrailerRuntime struct {
	ctx             context.Context
	re              typex.RuleX
	goodsProcessMap *sync.Map // Key: UUID, Value: GoodsProcess
}

func InitTrailerRuntime(re typex.RuleX) *TrailerRuntime {
	__DefaultTrailerRuntime = &TrailerRuntime{
		ctx:             typex.GCTX,
		re:              re,
		goodsProcessMap: &sync.Map{},
	}
	// 探针
	go func() {
		for {
			AllGoods().Range(func(key, value any) bool {
				goodsProcess := (value.(*GoodsProcess))
				grpcConnection, err := grpc.Dial(goodsProcess.NetAddr,
					grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					glogger.GLogger.Error(err)
				}
				client := NewTrailerClient(grpcConnection)
				probe(client, goodsProcess)
				grpcConnection.Close()
				return true
			})
			// 2秒停顿
			time.Sleep(2 * time.Second)
		}

	}()
	return __DefaultTrailerRuntime
}

/*
*
* 执行外
*
 */
func Fork(goods Goods) error {
	glogger.GLogger.Infof("fork goods process, (uuid = %v, addr = %v, args = %v)",
		goods.UUID, goods.LocalPath, goods.Args)
	Cmd := exec.Command(goods.LocalPath, goods.Args...)
	Cmd.SysProcAttr = NewSysProcAttr()
	Cmd.Stdin = os.Stdin
	Cmd.Stdout = os.Stdout
	Cmd.Stderr = os.Stderr
	ctx, Cancel := context.WithCancel(__DefaultTrailerRuntime.ctx)
	goodsProcess := &GoodsProcess{
		LocalPath:   goods.LocalPath,
		NetAddr:     goods.NetAddr,
		Uuid:        goods.UUID,
		Description: goods.Description,
		Args:        goods.Args,
		cmd:         Cmd,
		ctx:         ctx,
		cancel:      Cancel,
	}
	Save(goodsProcess)
	go run(goodsProcess) // 任务进程
	return nil
}

// 获取某个外挂
func Get(uuid string) *GoodsProcess {
	v, ok := __DefaultTrailerRuntime.goodsProcessMap.Load(uuid)
	if ok {
		return v.(*GoodsProcess)
	}
	return nil
}

// 保存进内存
func Save(goodsProcess *GoodsProcess) {
	__DefaultTrailerRuntime.goodsProcessMap.Store(goodsProcess.Uuid, goodsProcess)
}

// 从内存里删除, 删除后记得停止挂件, 通常外部配置表也要删除, 比如Sqlite
func Remove(uuid string) {
	v, ok := __DefaultTrailerRuntime.goodsProcessMap.Load(uuid)
	if ok {
		gp := (v.(*GoodsProcess))
		gp.Stop()
		__DefaultTrailerRuntime.goodsProcessMap.Delete(uuid)
	}
}

// 停止外挂运行时管理器, 这个要是停了基本上就是程序结束了
func Stop() {
	__DefaultTrailerRuntime.goodsProcessMap.Range(func(key, v interface{}) bool {
		gp := (v.(*GoodsProcess))
		gp.Stop()
		return true
	})
	__DefaultTrailerRuntime = nil
}

/*
*
* Cmd.Wait() 会阻塞, 但是当控制的子进程停止的时候会继续执行, 因此可以在defer里面释放资源
*
 */
func run(goodsProcess *GoodsProcess) error {
	defer func() {
		goodsProcess.Running = false
		goodsProcess.cancel() // 当监督器结束的时候探针probe进程也会被中断
		Remove(goodsProcess.Uuid)
	}()
	if err := goodsProcess.cmd.Start(); err != nil {
		glogger.GLogger.Error("exec command error:", err)
		return err
	}
	goodsProcess.Running = true
	glogger.GLogger.Infof("goods process(pid = %v, uuid = %v, addr = %v, args = %v) fork and started",
		goodsProcess.cmd.Process.Pid,
		goodsProcess.Uuid,
		goodsProcess.LocalPath,
		goodsProcess.Args)

	grpcConnection, err := grpc.Dial(goodsProcess.NetAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	defer grpcConnection.Close()
	client := NewTrailerClient(grpcConnection)
	if _, err := client.Init(goodsProcess.ctx, &Config{
		Kv: map[string]string{
			"uuid":        goodsProcess.Uuid,
			"description": goodsProcess.Description,
		},
	}); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	// Start
	if _, err := client.Start(goodsProcess.ctx, &Request{}); err != nil {
		glogger.GLogger.Error(err)
		return err
	}

	if err := goodsProcess.cmd.Wait(); err != nil {
		out, err1 := goodsProcess.cmd.Output()
		glogger.GLogger.Error("Cmd Wait error:", err, err1, string(out))
		return err
	}
	goodsProcess.Running = false
	return nil
}

// 探针
func probe(client TrailerClient, goodsProcess *GoodsProcess) {

	select {
	case <-goodsProcess.ctx.Done():
		{
			goodsProcess.Stop()
			glogger.GLogger.Infof("goods process(uuid = %v, addr = %v, args = %v) stopped",
				goodsProcess.Uuid,
				goodsProcess.NetAddr,
				goodsProcess.Args)
			return
		}
	default:
		{
			if goodsProcess.cmd != nil {
				if _, err := client.Status(goodsProcess.ctx, &Request{}); err != nil {
					glogger.GLogger.Error(err)
					goodsProcess.Running = false
				} else {
					goodsProcess.Running = true
					glogger.GLogger.Debug("goods Process is running:", goodsProcess.Uuid)
				}
			} else {
				goodsProcess.Running = false
			}
		}
	}

}

/*
*
* 返回外挂MAP
*
 */
func AllGoods() *sync.Map {
	return __DefaultTrailerRuntime.goodsProcessMap
}

/*
*
* 判断是否可执行(Linux Only)
*
 */
func IsExecutableFileUnix(filePath string) bool {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	if fileInfo.Mode()&0111 != 0 {
		return true
	}

	return false
}
func IsExecutableFileWin(filePath string) bool {
	filePath = strings.ToLower(filePath)
	return strings.HasSuffix(filePath, ".exe") ||
		strings.HasSuffix(filePath, ".jar") ||
		strings.HasSuffix(filePath, ".py") ||
		strings.HasSuffix(filePath, ".js") ||
		strings.HasSuffix(filePath, ".lua")

}
