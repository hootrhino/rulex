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
	"syscall"

	"os"
	"os/exec"
	"sync"

	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
)

//--------------------------------------------------------------------------------------------------
// Trailer 接口
//--------------------------------------------------------------------------------------------------

type TrailerManager struct {
	ctx             context.Context
	re              typex.RuleX
	goodsProcessMap *sync.Map // Key: UUID, Value: GoodsProcess
}

func NewTrailerManager(re typex.RuleX) typex.XTrailer {
	return &TrailerManager{
		ctx:             typex.GCTX,
		re:              re,
		goodsProcessMap: &sync.Map{},
	}
}

/*
*
* 执行外
*
 */
func (scm *TrailerManager) Fork(goods typex.Goods) error {
	glogger.GLogger.Infof("fork goods process, (uuid = %v, addr = %v, args = %v)", goods.UUID, goods.Addr, goods.Args)
	Cmd := exec.Command(goods.Addr, goods.Args...)
	Cmd.SysProcAttr = NewSysProcAttr()
	Cmd.Stdin = os.Stdin
	Cmd.Stdout = os.Stdout
	Cmd.Stderr = os.Stderr
	ctx, Cancel := context.WithCancel(scm.ctx)
	goodsProcess := &typex.GoodsProcess{
		Addr:        goods.Addr,
		Uuid:        goods.UUID,
		Description: goods.Description,
		Args:        goods.Args,
		Cmd:         Cmd,
		Ctx:         ctx,
		Cancel:      Cancel,
	}
	wg := sync.WaitGroup{}
	wg.Add(2)
	scm.Save(goodsProcess)
	go scm.run(&wg, goodsProcess)
	go scm.probe(&wg, goodsProcess)
	wg.Wait()
	return nil
}

// 获取某个外挂
func (scm *TrailerManager) Get(uuid string) *typex.GoodsProcess {
	v, ok := scm.goodsProcessMap.Load(uuid)
	if ok {
		return v.(*typex.GoodsProcess)
	}
	return nil
}

// 保存进内存
func (scm *TrailerManager) Save(goodsProcess *typex.GoodsProcess) {
	scm.goodsProcessMap.Store(goodsProcess.Uuid, goodsProcess)
}

// 从内存里删除, 删除后记得停止挂件, 通常外部配置表也要删除, 比如Sqlite
func (scm *TrailerManager) Remove(uuid string) {
	v, ok := scm.goodsProcessMap.Load(uuid)
	if ok {
		(v.(*typex.GoodsProcess)).Stop()
		scm.goodsProcessMap.Delete(uuid)
	}
}

// 停止外挂运行时管理器, 这个要是停了基本上就是程序结束了
func (scm *TrailerManager) Stop() {
	scm.goodsProcessMap.Range(func(key, value interface{}) bool {
		(value.(*typex.GoodsProcess)).Stop()
		return true
	})
	scm = nil
}

// Cmd.Wait() 会阻塞, 但是当控制的子进程停止的时候会继续执行, 因此可以在defer里面释放资源
func (scm *TrailerManager) run(wg *sync.WaitGroup, goodsProcess *typex.GoodsProcess) error {
	defer func() {
		goodsProcess.Cancel()
	}()
	if err := goodsProcess.Cmd.Start(); err != nil {
		glogger.GLogger.Error("exec command error:", err)
		wg.Done()
		return err
	}
	wg.Done()
	goodsProcess.Running = true
	glogger.GLogger.Infof("goods process(pid = %v, uuid = %v, addr = %v, args = %v) fork and started",
		goodsProcess.Cmd.Process.Pid,
		goodsProcess.Uuid,
		goodsProcess.Addr,
		goodsProcess.Args)
	if err := goodsProcess.Cmd.Wait(); err != nil {
		glogger.GLogger.Error("Cmd Wait error:", err)
		wg.Done()
		return err
	}
	goodsProcess.Running = false
	return nil
}

// 探针
func (scm *TrailerManager) probe(wg *sync.WaitGroup, goodsProcess *typex.GoodsProcess) {
	defer func() {
	}()
	wg.Done()
	for {
		select {
		case <-goodsProcess.Ctx.Done():
			{
				if goodsProcess.Cmd != nil {
					process := goodsProcess.Cmd.Process
					if process != nil {
						glogger.GLogger.Infof("goods process(pid = %v,uuid = %v, addr = %v, args = %v) stopped",
							goodsProcess.Cmd.Process.Pid,
							goodsProcess.Uuid,
							goodsProcess.Addr,
							goodsProcess.Args)
						process.Kill()
						process.Signal(syscall.SIGKILL)
					} else {
						glogger.GLogger.Infof("goods process(uuid = %v, addr = %v, args = %v) stopped",
							goodsProcess.Uuid,
							goodsProcess.Addr,
							goodsProcess.Args)
					}
				}
				scm.Remove(goodsProcess.Uuid)
				return
			}
		default:
			{
				if goodsProcess.Cmd != nil {
					if goodsProcess.Cmd.ProcessState != nil {
						goodsProcess.Running = true
					} else {
						goodsProcess.Running = false
					}
				}

			}
		}
	}
}

/*
*
* 返回外挂MAP
*
 */
func (scm *TrailerManager) AllGoods() *sync.Map {
	return scm.goodsProcessMap
}
