package engine

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"time"

	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/target"
	"github.com/i4de/rulex/typex"
)

/*
*
* 加载用户自定义输出资源
*
 */
func (e *RuleEngine) LoadUserOutEnd(target typex.XTarget, out *typex.OutEnd) error {
	return startTarget(target, out, e)
}

/*
*
* 加载内建输出资源
*
 */

func (e *RuleEngine) LoadBuiltinOutEnd(out *typex.OutEnd) error {
	return e.LoadOutEnd(out)
}
func (e *RuleEngine) LoadOutEnd(out *typex.OutEnd) error {
	if out.Type == typex.MONGO_SINGLE {
		return startTarget(target.NewMongoTarget(e), out, e)
	}
	if out.Type == typex.MQTT_TARGET {
		return startTarget(target.NewMqttTarget(e), out, e)
	}
	if out.Type == typex.NATS_TARGET {
		return startTarget(target.NewNatsTarget(e), out, e)
	}
	if out.Type == typex.HTTP_TARGET {
		return startTarget(target.NewHTTPTarget(e), out, e)
	}
	if out.Type == typex.TDENGINE_TARGET {
		return startTarget(target.NewTdEngineTarget(e), out, e)
	}
	if out.Type == typex.GRPC_CODEC_TARGET {
		return startTarget(target.NewCodecTarget(e), out, e)
	}
	return errors.New("unsupported target type:" + out.Type.String())
}

//
// Start output target
//
// Target life cycle:
//     Register -> Start -> running/restart cycle
//
func startTarget(target typex.XTarget, out *typex.OutEnd, e typex.RuleX) error {
	//
	// 先注册, 如果出问题了直接删除就行
	//
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
	// 然后启动资源
	ctx, cancelCTX := typex.NewCCTX()
	if err := target.Start(typex.CCTX{Ctx: ctx, CancelCTX: cancelCTX}); err != nil {
		glogger.GLogger.Error(err)
		e.RemoveOutEnd(out.UUID)
		return err
	}
	// Set sources to inend
	out.Target = target
	//
	ticker := time.NewTicker(time.Duration(time.Second * 5))
	tryIfRestartTarget(target, e, out.UUID)
	go func(ctx context.Context) {

		// 5 seconds
		//
	TICKER:
		<-ticker.C
		select {
		case <-ctx.Done():
			{
				ticker.Stop()
				return
			}
		default:
			{
				goto CHECK
			}
		}
	CHECK:
		{
			if target.Details() == nil {
				return
			}
			tryIfRestartTarget(target, e, out.UUID)
			goto TICKER
		}

	}(typex.GCTX)
	glogger.GLogger.Infof("Target [%v, %v] load successfully", out.Name, out.UUID)
	return nil
}

//
// 监测状态, 如果挂了重启
//
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
		ctx, cancelCTX := typex.NewCCTX()
		target.Start(typex.CCTX{Ctx: ctx, CancelCTX: cancelCTX})
	} else {
		target.Details().State = typex.SOURCE_UP
	}
}
