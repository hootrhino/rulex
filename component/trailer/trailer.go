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
	"net"
	"os"
	"strings"
	"time"

	"os/exec"
	"sync"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var __DefaultTrailerRuntime *TrailerRuntime

// --------------------------------------------------------------------------------------------------
// Trailer 接口
// --------------------------------------------------------------------------------------------------
type __TrailerRpcServer struct {
	UnimplementedTrailerServer
}

/*
*
* 只实现 OnStream 别的暂时先不管
*
 */
func (__TrailerRpcServer) OnStream(s Trailer_OnStreamServer) error {
	s.Send(&StreamResponse{Code: 1, Data: []byte("OK")})
	return nil
}

type TrailerRuntime struct {
	ctx             context.Context
	re              typex.RuleX
	goodsProcessMap *sync.Map // Key: UUID, Value: GoodsProcess
	rpcServer       *grpc.Server
	pid             int
	running         bool // true: running; false: stop
}

/*
*
* RULEX RPC Server 默认运行在 2588
*
 */
func InitTrailerRuntime(re typex.RuleX) *TrailerRuntime {
	__DefaultTrailerRuntime = &TrailerRuntime{
		ctx:             typex.GCTX,
		re:              re,
		goodsProcessMap: &sync.Map{},
		pid:             os.Getppid(),
		running:         true,
	}
	Listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 2588))
	if err != nil {
		glogger.GLogger.Error(err)
		return nil
	}
	__DefaultTrailerRuntime.rpcServer = grpc.NewServer(grpc.EmptyServerOption{})
	RegisterTrailerServer(__DefaultTrailerRuntime.rpcServer, __TrailerRpcServer{})
	// Stream Server
	go __DefaultTrailerRuntime.rpcServer.Serve(Listener)
	glogger.GLogger.Info("Trailer Runtime Proto Server Listening at [::]:2588")

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
				client, err := goodsProcess.ConnectToRpc()
				if err != nil {
					glogger.GLogger.Debug("ConnectToRpc error:", err)
					goodsProcess.ConnectToRpc() // 尝试重连
				} else {
					probe(client, goodsProcess) // 尝试重连
				}
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
func StartProcess(goods GoodsInfo) error {
	Remove(goods.UUID)
	return fork(goods)
}

/*
*
* 直接关闭
*
 */
func StopProcess(goods GoodsInfo) error {
	// Defer 删了以后这里不一定能拿到UUID
	v, ok := __DefaultTrailerRuntime.goodsProcessMap.Load(goods.UUID)
	if ok {
		gp := (v.(*GoodsProcess))
		gp.Stop()
		return nil
	}
	return fmt.Errorf("goods not exists:%s", goods.UUID)
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
		fmt.Sprintf("goods/console/%s", hk.ps.Info.UUID)).Debug(string(p))
	return len(p), nil
}
func (hk goodsStdInOut) Read(p []byte) (n int, err error) {
	return len(p), nil
}

/*
*
* 分离进程
*
 */
func fork(info GoodsInfo) error {
	glogger.GLogger.Infof("fork goods process, (uuid = %v, addr = %v, args = %v)",
		info.UUID, info.LocalPath, info.Args)
	ctx, Cancel := context.WithCancel(__DefaultTrailerRuntime.ctx)

	var Cmd *exec.Cmd
	args := strings.Split(info.Args, " ")
	tArgs := []string{info.LocalPath}
	// python main.py args...
	if info.ExecuteType == "PYTHON" {
		tArgs = append(tArgs, args...)
		Cmd = exec.CommandContext(ctx, "python", tArgs...)
	}
	// node main.js  args...
	if info.ExecuteType == "NODEJS" {
		tArgs = append(tArgs, args...)
		Cmd = exec.CommandContext(ctx, "node", tArgs...)
	}
	// lua main.lua args...
	if info.ExecuteType == "LUA" {
		tArgs = append(tArgs, args...)
		Cmd = exec.CommandContext(ctx, "lua", tArgs...)
	}
	//$ java -jar JarExample.jar args...
	if info.ExecuteType == "JAVA" {
		jarArgs := []string{"-jar"}
		tArgs = append(tArgs, args...)
		jarArgs = append(jarArgs, tArgs...)
		Cmd = exec.CommandContext(ctx, "java", jarArgs...)
	}
	if info.ExecuteType == "ELF" {
		Cmd = exec.CommandContext(ctx, info.LocalPath, args...)
	}
	if info.ExecuteType == "EXE" {
		Cmd = exec.CommandContext(ctx, info.LocalPath, args...)
	}
	if Cmd == nil {
		Cancel()
		return fmt.Errorf("unsupported executable file:%s", info.LocalPath)
	}
	glogger.GLogger.Debug("Execute system process:", Cmd.String())
	Cmd.SysProcAttr = NewSysProcAttr()
	goodsProcess := &GoodsProcess{
		Info:    info,
		cmd:     Cmd,
		ctx:     ctx,
		cancel:  Cancel,
		mailBox: make(chan int, 1),
	}
	inOut := NewWSStdInOut(goodsProcess)
	goodsProcess.cmd.Stdin = &inOut
	goodsProcess.cmd.Stdout = &inOut
	goodsProcess.cmd.Stderr = &inOut
	saveProcessMetaToMap(goodsProcess)
	go runLocalProcess(goodsProcess) // 任务进程
	return nil
}

/*
*
* 抢救进程
*
 */
func rescueRunLocalProcess(goodsProcess *GoodsProcess) error {
	return runLocalProcess(goodsProcess)
}

/*
*
* Cmd.Wait() 会阻塞, 但是当控制的子进程停止的时候会继续执行, 因此可以在defer里面释放资源
*  先保证本地Process进程启动，然后再回调RPC
 */
func runLocalProcess(goodsProcess *GoodsProcess) error {
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
					glogger.GLogger.Debug("goodsProcess.ctx.Done():",
						goodsProcess.Info.NetAddr)
					goodsProcess.cancel()
					return
				}
			case <-ctx.Done():
				{
					glogger.GLogger.Debug("goods Process Start timeout:",
						goodsProcess.Info.NetAddr)
					goodsProcess.cancel()
					return
				}
			default:
				{
				}
			}
			time.Sleep(2 * time.Second)
			// 尝试启动RPC
			glogger.GLogger.Debug("Wait Grpc Start:", goodsProcess.Info.NetAddr)
			if loadRpc(goodsProcess) {
				glogger.GLogger.Debug("Grpc Started:", goodsProcess.Info.NetAddr)
				return
			} else {
				glogger.GLogger.Debug("Grpc Started Failed:", goodsProcess.Info.NetAddr)
				goodsProcess.psRunning = false
			}
		}
	}()

	glogger.GLogger.Info("goods started:", goodsProcess.String())
	// Start 以后即可拿到Pid
	goodsProcess.pid = goodsProcess.cmd.Process.Pid
	if err := goodsProcess.cmd.Wait(); err != nil {
		State := goodsProcess.cmd.ProcessState
		if !State.Success() {
			glogger.GLogger.Error("Cmd Exit With State:", State)
			// 非正常结束, 极有可能是被kill的，所以要尝试抢救
			// 如果是被 RULEX 干死的就不抢救了，说明触发了 Stop 和 Remove；
			//    killedBy 如果是别的原因就有抢救机会
			// 还需要判断是否是主进程结束
			if !__DefaultTrailerRuntime.running {
				glogger.GLogger.Error("Trailer Runtime exited:", State)
				return nil
			}
			if goodsProcess.Info.KilledBy != "RULEX" {
				glogger.GLogger.Warn("Goods process Exit, May be a accident, try to rescue it:",
					goodsProcess.String())
				// 说明是用户操作停止
				time.Sleep(4 * time.Second)
				goodsProcess.Stop()
				PPId := os.Getppid()
				if PPId == __DefaultTrailerRuntime.pid {
					go rescueRunLocalProcess(goodsProcess)
				}
			} else {
				glogger.GLogger.Debug("Goods process killed by Rulex, No need to rescue it:",
					goodsProcess.String())
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
	__DefaultTrailerRuntime.goodsProcessMap.Store(goodsProcess.Info.UUID, goodsProcess)
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
	__DefaultTrailerRuntime.running = false
	__DefaultTrailerRuntime.goodsProcessMap.Range(func(key, v interface{}) bool {
		gp := (v.(*GoodsProcess))
		gp.StopBy("RULEX")
		return true
	})
	if __DefaultTrailerRuntime.rpcServer != nil {
		__DefaultTrailerRuntime.rpcServer.Stop()
	}
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
			glogger.GLogger.Info("goods process stopped:", goodsProcess.String())
			return
		}
	default:
		{
			if goodsProcess.cmd != nil {
				goodsProcess.psRunning = true
				if goodsProcess.rpcStarted {
					response, errStatus := client.Status(goodsProcess.ctx, &Request{})
					Status := response.GetStatus()
					if Status == StatusResponse_RUNNING && errStatus == nil {
						goodsProcess.rpcStarted = true
					} else {
						glogger.GLogger.Error(errStatus)
						goodsProcess.rpcStarted = false
					}
				}
				return
			} else {
				// 进程没起来，RPC也不会起来
				// 进程的 cmd==nil 时，说明已经挂了，尝试将其救活, 默认最多抢救5次
				goodsProcess.psRunning = false
				goodsProcess.rpcStarted = false
				glogger.GLogger.Debug("Goods Process is down try to restart:",
					goodsProcess.String())
				// 未来会根据 AutoStart判断是否重启进程
				// 字段已经在0.6.4加入
				goodsProcess.Stop()
				// 防止Windows下出现僵尸进程
				PPId := os.Getppid()
				if PPId == __DefaultTrailerRuntime.pid {
					go StartProcess(goodsProcess.Info)
				}
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
* 回调RPC接口, 让远程接口响应
*
 */
func loadRpc(goodsProcess *GoodsProcess) bool {
	grpcConnection, err := grpc.Dial(goodsProcess.Info.NetAddr,
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
			Kv: []byte(goodsProcess.Info.Args),
		}); err != nil {
			glogger.GLogger.Error("Init error:", goodsProcess.Info.NetAddr, ", error:", err)
			return false
		}
		// Start
		if _, err := client.Start(goodsProcess.ctx, &Request{}); err != nil {
			glogger.GLogger.Error("Start error:", goodsProcess.Info.NetAddr, ", error:", err)
			return false
		} else {
			goodsProcess.rpcStarted = true
			return true
		}
	}
	return false
}
