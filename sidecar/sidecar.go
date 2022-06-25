package sidecar

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

	"github.com/i4de/rulex/glogger"
)

//--------------------------------------------------------------------------------------------------
// SideCar 接口
//--------------------------------------------------------------------------------------------------

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

type SideCar interface {
	Fork(Goods) error
	Get(addr string) *GoodsProcess
	Save(*GoodsProcess)
	Remove(uuid string)
	AllGoods() *sync.Map
	Stop()
}

//--------------------------------------------------------------------------------------------------
// GoodsProcess
//--------------------------------------------------------------------------------------------------

type GoodsProcess struct {
	running     bool
	uuid        string
	addr        string
	description string
	args        []string
	ctx         context.Context
	cmd         *exec.Cmd
	cancel      context.CancelFunc
}

func (t *GoodsProcess) Running() bool {
	return t.running
}
func (t *GoodsProcess) Description() string {
	return t.description
}
func (t *GoodsProcess) UUID() string {
	return t.uuid
}
func (t *GoodsProcess) Addr() string {
	return t.addr
}
func (t *GoodsProcess) Args() []string {
	return t.args
}
func (t GoodsProcess) String() string {
	r := map[string]interface{}{
		"running":     t.running,
		"uuid":        t.uuid,
		"addr":        t.addr,
		"description": t.description,
		"args":        t.args,
	}
	b, _ := json.Marshal(r)
	return string(b)
}
func (scm *GoodsProcess) Stop() {
	if scm.cmd != nil {
		if scm.cmd.Process != nil {
			scm.cmd.Process.Kill()
			scm.cmd = nil
			scm.cancel()
		}
	}
}

func NewGoodsProcess() *GoodsProcess {
	return &GoodsProcess{}
}

//--------------------------------------------------------------------------------------------------
// SideCar
//--------------------------------------------------------------------------------------------------

type SidecarManager struct {
	ctx             context.Context
	goodsProcessMap *sync.Map // Key: UUID, Value: GoodsProcess
}

func NewSideCarManager(ctx context.Context) SideCar {
	return &SidecarManager{
		ctx:             ctx,
		goodsProcessMap: &sync.Map{},
	}
}

/*
*
* 执行外
*
 */
func (scm *SidecarManager) Fork(goods Goods) error {
	glogger.GLogger.Infof("fork goods process, (uuid = %v, addr = %v, args = %v)", goods.UUID, goods.Addr, goods.Args)
	cmd := exec.Command(goods.Addr, goods.Args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	ctx, cancel := context.WithCancel(scm.ctx)
	goodsProcess := &GoodsProcess{
		addr:        goods.Addr,
		uuid:        goods.UUID,
		description: goods.Description,
		args:        goods.Args,
		cmd:         cmd,
		ctx:         ctx,
		cancel:      cancel,
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

//
// 从内存里删除, 删除后记得停止挂件, 通常外部配置表也要删除, 比如Sqlite
//
func (scm *SidecarManager) Remove(uuid string) {
	v, ok := scm.goodsProcessMap.Load(uuid)
	if ok {
		(v.(*GoodsProcess)).Stop()
		scm.goodsProcessMap.Delete(uuid)
	}
}

//
// 停止外挂运行时管理器, 这个要是停了基本上就是程序结束了
//
func (scm *SidecarManager) Stop() {
	scm.goodsProcessMap.Range(func(key, value interface{}) bool {
		(value.(*GoodsProcess)).Stop()
		return true
	})
	scm = nil
}

//
// cmd.Wait() 会阻塞, 但是当控制的子进程停止的时候会继续执行, 因此可以在defer里面释放资源
//
func (scm *SidecarManager) run(wg *sync.WaitGroup, goodsProcess *GoodsProcess) error {
	defer func() {
		goodsProcess.cancel()
	}()
	if err := goodsProcess.cmd.Start(); err != nil {
		glogger.GLogger.Error("exec command error:", err)
		wg.Done()
		return err
	}
	wg.Done()
	goodsProcess.running = true
	glogger.GLogger.Infof("goods process(pid = %v, uuid = %v, addr = %v, args = %v) fork and started",
		goodsProcess.cmd.Process.Pid,
		goodsProcess.uuid,
		goodsProcess.addr,
		goodsProcess.args)
	if err := goodsProcess.cmd.Wait(); err != nil {
		glogger.GLogger.Error("cmd Wait error:", err)
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
				if goodsProcess.cmd != nil {
					process := goodsProcess.cmd.Process
					if process != nil {
						glogger.GLogger.Infof("goods process(pid = %v,uuid = %v, addr = %v, args = %v) stopped",
							goodsProcess.cmd.Process.Pid,
							goodsProcess.uuid,
							goodsProcess.addr,
							goodsProcess.args)
						process.Kill()
						process.Signal(syscall.SIGKILL)
					} else {
						glogger.GLogger.Infof("goods process(uuid = %v, addr = %v, args = %v) stopped",
							goodsProcess.uuid,
							goodsProcess.addr,
							goodsProcess.args)
					}
				}
				scm.Remove(goodsProcess.uuid)
				return
			}
		default:
			{
				if goodsProcess.cmd != nil {
					if goodsProcess.cmd.ProcessState != nil {
						goodsProcess.running = true
					} else {
						goodsProcess.running = false
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
func (scm *SidecarManager) AllGoods() *sync.Map {
	return scm.goodsProcessMap
}
