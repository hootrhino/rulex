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
	"fmt"
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
			select {
			case <-typex.GCTX.Done():
				{
					return
				}
			default:
				{
				}
			}
			// glogger.GLogger.Debug("Start prob process.")
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
* 直接启动
*
 */
func StartProcess(goods Goods) error {
	return fork(goods)
}

/*
*
* 直接关闭
*
 */
func StopProcess(goods Goods) error {
	v, ok := __DefaultTrailerRuntime.goodsProcessMap.Load(goods.UUID)
	if ok {
		gp := (v.(*GoodsProcess))
		gp.Stop()
		return nil
	}
	return fmt.Errorf("Goods not exists:%s", goods.UUID)
}

/*
*
* Fork 一个进程来执行
*
 */

type goodsStdInOut struct {
	ps *GoodsProcess
}

func NewWSStdInOut(ps *GoodsProcess) goodsStdInOut {
	return goodsStdInOut{ps: ps}
}

func (hk goodsStdInOut) Write(p []byte) (n int, err error) {
	glogger.Logrus.WithField("topic",
		fmt.Sprintf("goods/console/%s", hk.ps.Uuid)).Debug(string(p))
	return 0, nil
}

/*
*
* 分离进程
*
 */
func fork(goods Goods) error {
	glogger.GLogger.Infof("fork goods process, (uuid = %v, addr = %v, args = %v)",
		goods.UUID, goods.LocalPath, goods.Args)
	ctx, Cancel := context.WithCancel(__DefaultTrailerRuntime.ctx)
	Cmd := exec.CommandContext(ctx, goods.LocalPath, goods.Args)
	Cmd.SysProcAttr = NewSysProcAttr()
	goodsProcess := &GoodsProcess{
		Pid:         0,
		LocalPath:   goods.LocalPath,
		NetAddr:     goods.NetAddr,
		Uuid:        goods.UUID,
		Description: goods.Description,
		Args:        goods.Args,
		cmd:         Cmd,
		ctx:         ctx,
		cancel:      Cancel,
		mailBox:     make(chan int, 1),
	}
	// out := NewWSStdInOut(goodsProcess)
	Cmd.Stdin = nil
	Cmd.Stdout = os.Stdout
	Cmd.Stderr = os.Stdout
	saveProcessMetaToMap(goodsProcess)
	go runLocalProcess(goodsProcess) // 任务进程
	return nil
}

/*
*
* Cmd.Wait() 会阻塞, 但是当控制的子进程停止的时候会继续执行, 因此可以在defer里面释放资源
*  先保证本地Process进程启动，然后再回调RPC
 */
func runLocalProcess(goodsProcess *GoodsProcess) error {
	defer func() {
		// Remove 已经包含了Stop, cancel
		Remove(goodsProcess.Uuid)
	}()
	// 到这里就挂了的,说明参数错了,不值得救活
	// 最好是直接删了或者更新配置
	if err := goodsProcess.cmd.Start(); err != nil {
		glogger.GLogger.Error("exec command error:", err)
		return err
	}

	var client TrailerClient
	// 迁移到prob,改造成监督进程
	// Load OS process
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		for {
			select {
			case <-goodsProcess.ctx.Done():
				{
					glogger.GLogger.Debug("goodsProcess.ctx.Done():", goodsProcess.NetAddr)
					goodsProcess.cancel()
					return
				}
			case <-ctx.Done():
				{
					glogger.GLogger.Debug("goods Process Start timeout:", goodsProcess.NetAddr)
					goodsProcess.cancel()
					return
				}
			default:
				{
				}
			}
			time.Sleep(2 * time.Second)
			// 尝试启动RPC
			glogger.GLogger.Debug("Wait Grpc Start:", goodsProcess.NetAddr)
			if loadRpc(goodsProcess) {
				glogger.GLogger.Debug("Grpc Started:", goodsProcess.NetAddr)
				return
			}
		}
	}()

	glogger.GLogger.Infof("goods process(pid = %v, uuid = %v, addr = %v, args = %v) fork and started",
		goodsProcess.cmd.Process.Pid,
		goodsProcess.Uuid,
		goodsProcess.LocalPath,
		goodsProcess.Args)
	// Start 以后即可拿到Pid
	goodsProcess.Pid = goodsProcess.cmd.Process.Pid
	if err := goodsProcess.cmd.Wait(); err != nil {
		State := goodsProcess.cmd.ProcessState
		if !State.Success() {
			out, err1 := goodsProcess.cmd.Output()
			glogger.GLogger.Error("Cmd Exit With State:", err, err1, string(out), ",State:", State)
		}
		// 非正常结束, 极有可能是被kill的，所以要尝试抢救
		// 问题：怎么知道是被kill或者自动结束?
		if State.ExitCode() != 0 {
			// 如果是被 RULEX 干死的就不抢救了，说明触发了 Stop 和 Remove；
			//    killedBy 如果是别的原因就有抢救机会
			if goodsProcess.killedBy != "RULEX" {
				glogger.GLogger.Warn("Goods process Exit, May be a accident, try to rescue it:", goodsProcess.Uuid)
				// 说明是用户操作停止
				time.Sleep(2 * time.Second)
				go fork(Goods{
					UUID:        goodsProcess.Uuid,
					LocalPath:   goodsProcess.LocalPath,
					NetAddr:     goodsProcess.NetAddr,
					Description: goodsProcess.Uuid,
					Args:        goodsProcess.Args,
				})
			} else {
				glogger.GLogger.Debug("Goods process killed by Rulex, No need to rescue it:", goodsProcess.Uuid)
			}
		}

	}
	if client != nil {
		client.Stop(goodsProcess.ctx, &Request{})
	}
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
func saveProcessMetaToMap(goodsProcess *GoodsProcess) {
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

// 从内存里删除, 删除后记得停止挂件, 通常外部配置表也要删除, 比如Sqlite
func RemoveBy(uuid, by string) {
	v, ok := __DefaultTrailerRuntime.goodsProcessMap.Load(uuid)
	if ok {
		gp := (v.(*GoodsProcess))
		gp.StopBy(by)
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
}

/*
*
* 这个探针的主要作用就是监控进程挂了没，如果挂了要不要救活等
*
 */
func probe(client TrailerClient, goodsProcess *GoodsProcess) {
	select {
	case <-goodsProcess.ctx.Done():
		{
			glogger.GLogger.Infof("goods process(uuid = %v, addr = %v, args = %v) stopped",
				goodsProcess.Uuid,
				goodsProcess.NetAddr,
				goodsProcess.Args)
			return
		}
	default:
		{
			if goodsProcess.cmd != nil {
				goodsProcess.PsRunning = true
				if goodsProcess.rpcStarted {
					response, err := client.Status(goodsProcess.ctx, &Request{})
					Status := response.GetStatus()
					if Status == StatusResponse_RUNNING && err == nil {
						goodsProcess.rpcStarted = true
					} else {
						glogger.GLogger.Error(err)
						goodsProcess.rpcStarted = false
					}
				}
				return
			} else {
				// 进程没起来，RPC也不会起来
				// 进程的 cmd==nil 时，说明已经挂了，尝试将其救活, 默认最多抢救5次
				goodsProcess.PsRunning = false
				goodsProcess.rpcStarted = false
				glogger.GLogger.Debug("Goods Process is down:",
					goodsProcess.Uuid, " try to restart")
				go fork(Goods{
					UUID:        goodsProcess.Uuid,
					LocalPath:   goodsProcess.LocalPath,
					NetAddr:     goodsProcess.NetAddr,
					Description: goodsProcess.Uuid,
					Args:        goodsProcess.Args,
				})
				return
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

/*
*
* 回调RPC接口, 让远程接口响应
*
 */
func loadRpc(goodsProcess *GoodsProcess) bool {
	grpcConnection, err := grpc.Dial(goodsProcess.NetAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		glogger.GLogger.Error(err)
	}
	defer grpcConnection.Close()
	client := NewTrailerClient(grpcConnection)
	// glogger.GLogger.Debug("Try to start:", goodsProcess.NetAddr)
	// 等进程起来以后RPC调用
	if goodsProcess.cmd != nil {
		if _, err := client.Init(goodsProcess.ctx, &Config{
			Kv: []byte(goodsProcess.Args),
		}); err != nil {
			glogger.GLogger.Error("Init error:", goodsProcess.NetAddr, ", error:", err)
			return false
		}
		// Start
		if _, err := client.Start(goodsProcess.ctx, &Request{}); err != nil {
			glogger.GLogger.Error("Start error:", goodsProcess.NetAddr, ", error:", err)
			return false
		} else {
			goodsProcess.rpcStarted = true
			return true
		}
	}
	return false
}
