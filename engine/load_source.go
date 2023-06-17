package engine

import (
	"context"
	"errors"
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
	return startSources(source, in, e)
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
		return startSources(config.Source, in, e)
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
func startSources(source typex.XSource, in *typex.InEnd, e *RuleEngine) error {
	//
	// 先注册, 如果出问题了直接删除就行
	//
	// 首先把资源ID给注册进去, 作为资源的全局索引，确保资源可以拿到配置
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
	// Set sources to inend
	in.Source = source
	// 然后启动资源
	if err := startSource(source, e); err != nil {
		glogger.GLogger.Error(err)
		e.RemoveInEnd(in.UUID)
		return err
	}
	go func(ctx context.Context) {
		ticker := time.NewTicker(time.Duration(time.Second * 5))
		// 5 seconds
	TICKER:
		select {
		case <-ctx.Done():
			{
				ticker.Stop()
				return
			}
		default:
			{
				<-ticker.C
				goto CHECK
			}
		}
	CHECK:
		{
			//
			// 通过HTTP删除资源的时候, 会把数据清了, 只要检测到资源没了, 这里也退出
			//
			if source.Details() == nil {
				return
			}
			//------------------------------------
			// 驱动挂了资源也挂了, 因此检查驱动状态在先
			//------------------------------------
			tryIfRestartSource(source, e)
			//------------------------------------
			goto TICKER
		}

	}(typex.GCTX)
	glogger.GLogger.Infof("InEnd [%v, %v] load successfully", in.Name, in.UUID)
	return nil
}

/*
*
* 检查是否需要重新拉起资源
* 这里也有优化点：不能手动控制内存回收可能会产生垃圾
*
 */
func checkSourceDriverState(source typex.XSource) {
	if source.Driver() == nil {
		return
	}

	// 只有资源启动状态才拉起驱动
	if source.Status() == typex.SOURCE_UP {
		// 必须资源启动, 驱动才有重启意义
		if source.Driver().State() == typex.DRIVER_DOWN {
			glogger.GLogger.Warn("Driver down:", source.Driver().DriverDetail().Name)
			// 只需要把资源给拉闸, 就会触发重启
			source.Stop()
		}

	}

}

// test SourceState
func tryIfRestartSource(source typex.XSource, e *RuleEngine) {
	checkSourceDriverState(source)
	if source.Status() == typex.SOURCE_STOP {
		return
	}
	if source.Status() == typex.SOURCE_DOWN {
		source.Details().SetState(typex.SOURCE_DOWN)
		//----------------------------------
		// 当资源挂了以后先给停止, 然后重启
		//----------------------------------
		glogger.GLogger.Warnf("Source %v %v down. try to restart it", source.Details().UUID, source.Details().Name)
		source.Stop()
		//----------------------------------
		// 主动垃圾回收一波
		//----------------------------------
		runtime.Gosched()
		runtime.GC() // GC 比较慢, 但是是良性卡顿, 问题不大
		startSource(source, e)
	} else {
		source.Details().SetState(typex.SOURCE_UP)
	}
}

func startSource(source typex.XSource, e *RuleEngine) error {
	//----------------------------------
	// 检查资源 如果是启动的，先给停了
	//----------------------------------
	ctx, cancelCTX := typex.NewCCTX()

	if err := source.Start(typex.CCTX{Ctx: ctx, CancelCTX: cancelCTX}); err != nil {
		glogger.GLogger.Error("Source start error:", err)
		if source.Status() == typex.SOURCE_UP {
			source.Stop()
		}
		if source.Driver() != nil {
			if source.Driver().State() == typex.DRIVER_UP {
				source.Driver().Stop()
			}
		}
		return err
	}
	//----------------------------------
	// 驱动也要停了, 然后重启
	//----------------------------------
	if source.Driver() != nil {
		if source.Driver().State() == typex.DRIVER_UP {
			source.Driver().Stop()
		}
		// Start driver
		if err := source.Driver().Init(map[string]string{}); err != nil {
			glogger.GLogger.Error("Driver initial error:", err)
			return errors.New("Driver initial error:" + err.Error())
		}
		glogger.GLogger.Infof("Try to start driver: [%v]", source.Driver().DriverDetail().Name)
		if err := source.Driver().Work(); err != nil {
			glogger.GLogger.Error("Driver work error:", err)
			return errors.New("Driver work error:" + err.Error())
		}
		glogger.GLogger.Infof("Driver start successfully: [%v]", source.Driver().DriverDetail().Name)
	}
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
