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
* 加载用户自定义输出资源
*
 */
func (e *RuleEngine) LoadUserOutEnd(target typex.XTarget, out *typex.OutEnd) error {
	return loadTarget(target, out, e)
}

/*
*
* 加载内建输出资源
 */
func (e *RuleEngine) LoadBuiltinOutEnd(out *typex.OutEnd) error {
	return e.LoadOutEnd(out)
}
func (e *RuleEngine) LoadOutEnd(out *typex.OutEnd) error {
	if config := e.TargetTypeManager.Find(out.Type); config != nil {
		return loadTarget(config.Target, out, e)
	}
	return errors.New("unsupported target type:" + out.Type.String())
}

// Start output target
//
// Target life cycle:
//
//	Register -> Start -> running/restart cycle
func loadTarget(target typex.XTarget, out *typex.OutEnd, e typex.RuleX) error {
	// Set sources to inend
	out.Target = target
	e.SaveOutEnd(out)
	// Load config
	config := e.GetOutEnd(out.UUID).Config
	if config == nil {
		e.RemoveOutEnd(out.UUID)
		err := fmt.Errorf("target [%v] config is nil", out.Name)
		return err
	}
	if err := target.Init(out.UUID, config); err != nil {
		glogger.GLogger.Error(err)
		e.RemoveInEnd(out.UUID)
		return err
	}
	// start
	// if err := startTarget(target, e); err != nil {
	// 	glogger.GLogger.Error(err)
	// 	e.RemoveOutEnd(out.UUID)
	// 	return err
	// }
	//
	startTarget(target, e)
	go func(ctx context.Context) {
		ticker := time.NewTicker(time.Duration(time.Second * 5))
		for {
			select {
			case <-ctx.Done():
				{
					ticker.Stop()
					return
				}
			default:
				{
				}
			}
			<-ticker.C
			if target.Details() == nil {
				return
			}
			tryIfRestartTarget(target, e, out.UUID)
		}
	}(typex.GCTX)
	glogger.GLogger.Infof("Target [%v, %v] load successfully", out.Name, out.UUID)
	return nil
}

// 监测状态, 如果挂了重启
func tryIfRestartTarget(target typex.XTarget, e typex.RuleX, id string) {
	if target.Status() == typex.SOURCE_STOP {
		return
	}
	if target.Status() == typex.SOURCE_DOWN {
		target.Details().State = typex.SOURCE_DOWN
		glogger.GLogger.Warnf("Target [%v, %v] down. try to restart it", target.Details().Name, target.Details().UUID)
		target.Stop()
		runtime.Gosched()
		runtime.GC()
		startTarget(target, e)
	} else {
		target.Details().State = typex.SOURCE_UP
	}
}

func startTarget(target typex.XTarget, e typex.RuleX) error {
	ctx, cancelCTX := typex.NewCCTX()
	if err := target.Start(typex.CCTX{Ctx: ctx, CancelCTX: cancelCTX}); err != nil {
		glogger.GLogger.Error("abstractDevice start error:", err)
		return err
	}
	return nil
}
