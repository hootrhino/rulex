package engine

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
)

/*
*
* 加载用户资源
*
 */
func (e *RuleEngine) LoadUserInEnd(source typex.XSource, in *typex.InEnd) error {
	return e.loadSource(source, in)
}

/*
*
* 内建资源
*
 */
func (e *RuleEngine) LoadBuiltInEnd(in *typex.InEnd) error {
	return e.LoadInEnd(in)
}

/*
*
* 加载输入资源
*
 */
func (e *RuleEngine) LoadInEnd(in *typex.InEnd) error {
	if config := e.SourceTypeManager.Find(in.Type); config != nil {
		return e.loadSource(config.NewSource(e), in)
	}
	return fmt.Errorf("unsupported InEnd type:%s", in.Type)
}

//
// start Sources
//
/*
* Life cycle
+------------------+       +------------------+   +---------------+        +---------------+
|     Init         |------>|   Start          |-->|     Test      |--+ --->|  Stop         |
+------------------+  ^    +------------------+   +---------------+  |     +---------------+
                      |                                              |
                      |                                              |
                      +-------------------Error ---------------------+
*/
func (e *RuleEngine) loadSource(source typex.XSource, in *typex.InEnd) error {
	in.Source = source
	e.SaveInEnd(in)
	// Load config
	config := e.GetInEnd(in.UUID).Config
	if config == nil {
		e.RemoveInEnd(in.UUID)
		err := fmt.Errorf("source [%v] config is nil", in.Name)
		return err
	}
	if err := source.Init(in.UUID, config); err != nil {
		glogger.GLogger.Error(err)
		e.RemoveInEnd(in.UUID)
		return err
	}
	// 然后启动资源
	ctx, cancelCTX := typex.NewCCTX()
	go func(ctx1 context.Context, source1 typex.XSource) {
		acc++
		defer func() {
			println("、、、、、、、、、、、、Defer", acc)
		}()
		for {
			ticker := time.NewTicker(time.Duration(time.Second * 5))
			select {
			case <-ctx1.Done():
				{
					println("、、、、、、、、、、、、收到消息了 把这个进程停了，重启下一个", acc)
					ticker.Stop()
					return
				}
			default:
				{
				}
			}
			if source1.Status() == typex.SOURCE_STOP {
				return
			}
			println("、、、、、、、、、、、、检查资源状态", source1.Status(), acc)
			Status := source.Status()
			if Status == typex.SOURCE_STOP {
				return
			}
			if Status == typex.SOURCE_DOWN {
				source.Details().State = typex.SOURCE_DOWN
				glogger.GLogger.Warnf("Device [%v, %v] down. try to restart.",
					source.Details().UUID, source.Details().Name)
				source.Stop()
				runtime.Gosched()
				runtime.GC()
				time.Sleep(5 * time.Second)
				e.loadSource(source, source.Details())
				return
			}
			<-ticker.C
		}
	}(ctx, source)
	if err := e.startSource(source, ctx, cancelCTX); err != nil {
		return err
	}
	glogger.GLogger.Infof("InEnd [%v, %v] load successfully", in.Name, in.UUID)
	return nil
}

var acc int = 100

func (e *RuleEngine) startSource(source typex.XSource,
	ctx context.Context, cancelCTX context.CancelFunc) error {
	//----------------------------------
	// 检查资源 如果是启动的，先给停了
	//----------------------------------
	if err := source.Start(typex.CCTX{Ctx: ctx, CancelCTX: cancelCTX}); err != nil {
		glogger.GLogger.Error("Source start error:", err)
		source.Stop()
		return err
	}
	// LoadNewestSource
	// 2023-06-14新增： 重启成功后数据会丢失,还得加载最新的Rule到设备中
	Source := source.Details()
	if Source != nil {
		for _, rule := range Source.BindRules {
			RuleInstance := typex.NewLuaRule(e,
				rule.UUID,
				rule.Name,
				rule.Description,
				rule.FromSource,
				rule.FromDevice,
				rule.Success,
				rule.Actions,
				rule.Failed)
			if err1 := e.LoadRule(RuleInstance); err1 != nil {
				return err1
			}
		}
	}
	return nil
}

/*
*
* 检查是否需要重新拉起资源
* 这里也有优化点：不能手动控制内存回收可能会产生垃圾
*
 */
// test SourceState
func (e *RuleEngine) tryIfRestartSource(source typex.XSource) {
	Status := source.Status()
	if Status == typex.SOURCE_STOP {
		return
	}
	if Status == typex.SOURCE_DOWN {
		source.Details().State = typex.SOURCE_DOWN
		glogger.GLogger.Warnf("Device [%v, %v] down. try to restart.",
			source.Details().UUID, source.Details().Name)
		source.Stop()
		runtime.Gosched()
		runtime.GC()
		e.loadSource(source, source.Details())
	}

}
