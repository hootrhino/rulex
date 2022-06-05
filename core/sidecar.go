package core

//
// Sidecar就是拖车,带着小车一起跑,比喻了SideCar实际上是个进程管理器
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
	"encoding/json"
	"syscall"

	"os"
	"os/exec"
	"sync"

	"github.com/ngaut/log"
)

//--------------------------------------------------------------------------------------------------
// SideCar 接口
//--------------------------------------------------------------------------------------------------

type SideCar interface {
	Fork(name string, uuid string, args []string) error
	Get(name string) *GoodsProcess
	Save(*GoodsProcess)
	Remove(*GoodsProcess)
	Stop()
}

//--------------------------------------------------------------------------------------------------
// GoodsProcess
//--------------------------------------------------------------------------------------------------

type GoodsProcess struct {
	running bool
	uuid    string
	name    string
	args    []string
	ctx     context.Context
	cmd     *exec.Cmd
	cancel  context.CancelFunc
}

func (t GoodsProcess) String() string {
	r := map[string]interface{}{
		"running": t.running,
		"uuid":    t.uuid,
		"name":    t.name,
		"args":    t.args,
	}
	b, _ := json.Marshal(r)
	return string(b)
}
func (scm *GoodsProcess) Stop() {
	scm.cmd.Process.Kill()
	scm.cancel()
	scm.cmd = nil
}

func NewGoodsProcess() *GoodsProcess {
	return &GoodsProcess{}
}

//--------------------------------------------------------------------------------------------------
// SideCar
//--------------------------------------------------------------------------------------------------

type SidecarManager struct {
	ctx             context.Context
	goodsProcessMap sync.Map // Key: UUID, Value: GoodsProcess
}

func NewSideCarManager(ctx context.Context) SideCar {
	return &SidecarManager{
		ctx:             ctx,
		goodsProcessMap: sync.Map{},
	}
}

/*
*
* 执行外挂
*
 */
func (scm *SidecarManager) Fork(name string, uuid string, args []string) error {
	log.Infof("fork goods process, (uuid = %v, name = %v, args = %v)", uuid, name, args)
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	ctx, cancel := context.WithCancel(scm.ctx)
	goodsProcess := &GoodsProcess{
		name:   name,
		uuid:   uuid,
		args:   args,
		cmd:    cmd,
		ctx:    ctx,
		cancel: cancel,
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
func (scm *SidecarManager) Get(uuid string) *GoodsProcess {
	v, ok := scm.goodsProcessMap.Load(uuid)
	if ok {
		return v.(*GoodsProcess)
	}
	return nil
}

// 保存进内存
func (scm *SidecarManager) Save(goodsProcess *GoodsProcess) {
	scm.goodsProcessMap.Store(goodsProcess.uuid, goodsProcess)
}

// 从内存里删除
func (scm *SidecarManager) Remove(goodsProcess *GoodsProcess) {
	_, ok := scm.goodsProcessMap.Load(goodsProcess.uuid)
	if ok {
		scm.goodsProcessMap.Delete(goodsProcess.uuid)
	}
}

// 停止外挂
func (scm *SidecarManager) Stop() {
	scm.goodsProcessMap.Range(func(key, value interface{}) bool {
		(value.(*GoodsProcess)).Stop()
		return true
	})

}

//
// cmd.Wait() 会阻塞, 但是当控制的子进程停止的时候会继续执行, 因此可以在defer里面释放资源
//
func (scm *SidecarManager) run(wg *sync.WaitGroup, goodsProcess *GoodsProcess) error {
	defer func() {
		goodsProcess.cancel()
	}()
	if err := goodsProcess.cmd.Start(); err != nil {
		log.Error("exec command error:", err)
		wg.Done()
		return err
	}
	wg.Done()
	goodsProcess.running = true
	log.Infof("goods process(pid = %v, uuid = %v, name = %v, args = %v) fork and started: ",
		goodsProcess.cmd.Process.Pid,
		goodsProcess.uuid,
		goodsProcess.name,
		goodsProcess.args)
	if err := goodsProcess.cmd.Wait(); err != nil {
		log.Error("cmd Wait error:", err)
		wg.Done()
		return err
	}
	goodsProcess.running = false
	return nil
}

// 探针
func (scm *SidecarManager) probe(wg *sync.WaitGroup, goodsProcess *GoodsProcess) {
	defer func() {
	}()
	wg.Done()
	for {
		select {
		case <-goodsProcess.ctx.Done():
			{
				process := goodsProcess.cmd.Process
				if process != nil {
					log.Infof("goods process(pid = %v,uuid = %v, name = %v, args = %v) stopped",
						goodsProcess.cmd.Process.Pid,
						goodsProcess.uuid,
						goodsProcess.name,
						goodsProcess.args)
					process.Kill()
					process.Signal(syscall.SIGKILL)
				} else {
					log.Infof("goods process(uuid = %v, name = %v, args = %v) stopped",
						goodsProcess.uuid,
						goodsProcess.name,
						goodsProcess.args)
				}
				scm.Remove(goodsProcess)
				return
			}
		default:
			{
				if goodsProcess.cmd.ProcessState != nil {
					goodsProcess.running = true
				} else {
					goodsProcess.running = false
				}
			}
		}
	}
}
