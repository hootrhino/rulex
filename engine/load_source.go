package engine

import (
	"context"
	"errors"
	"fmt"
	"rulex/source"
	"rulex/typex"
	"runtime"
	"time"

	"github.com/ngaut/log"
)

/*
*
* TODO: 0.3.0重构此处，换成 SourceRegistry 形式
*
 */
func (e *RuleEngine) LoadInEnd(in *typex.InEnd) error {
	if in.Type == typex.MQTT {
		return startSources(source.NewMqttInEndSource(in.UUID, e), in, e)
	}
	if in.Type == typex.HTTP {
		return startSources(source.NewHttpInEndSource(in.UUID, e), in, e)
	}
	if in.Type == typex.COAP {
		return startSources(source.NewCoAPInEndSource(in.UUID, e), in, e)
	}
	if in.Type == typex.GRPC {
		return startSources(source.NewGrpcInEndSource(in.UUID, e), in, e)
	}
	if in.Type == typex.UART_MODULE {
		return startSources(source.NewUartModuleSource(in.UUID, e), in, e)
	}
	if in.Type == typex.MODBUS_MASTER {
		return startSources(source.NewModbusMasterSource(in.UUID, e), in, e)
	}
	if in.Type == typex.SNMP_SERVER {
		return startSources(source.NewSNMPInEndSource(in.UUID, e), in, e)
	}
	if in.Type == typex.NATS_SERVER {
		return startSources(source.NewNatsSource(e), in, e)
	}
	if in.Type == typex.SIEMENS_S7 {
		return startSources(source.NewSiemensS7Source(e), in, e)
	}
	if in.Type == typex.RULEX_UDP {
		return startSources(source.NewUdpInEndSource(e), in, e)
	}
	if in.Type == typex.RTU485_THER {
		return startSources(source.NewRtu485THerSource(e), in, e)
	}
	return fmt.Errorf("unsupported InEnd type:%s", in.Type)
}

//
// start Sources
//
/*
* Life cycle
+------------------+       +------------------+   +---------------+        +---------------+
|     Register     |------>|   Start          |-->|     Test      |--+ --->|  Stop         |
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
		log.Error(err)
		e.RemoveInEnd(in.UUID)
		return err
	}
	// Set sources to inend
	in.Source = source
	// 然后启动资源
	if err := startSource(source, e); err != nil {
		log.Error(err)
		e.RemoveInEnd(in.UUID)
		return err
	}
	ticker := time.NewTicker(time.Duration(time.Second * 5))
	go func(ctx context.Context) {
		// 5 seconds
	TICKER:
		<-ticker.C
		select {
		case <-ctx.Done():
			{
				return
			}
		default:
			{
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
	log.Infof("InEnd [%v, %v] load successfully", in.Name, in.UUID)
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
		if source.Driver().State() == typex.DRIVER_STOP {
			log.Warn("Driver stopped:", source.Driver().DriverDetail().Name)
			// 只需要把资源给拉闸, 就会触发重启
			source.Stop()
		}

	}

}

//
// test SourceState
//
func tryIfRestartSource(source typex.XSource, e *RuleEngine) {
	checkSourceDriverState(source)
	if source.Status() == typex.SOURCE_DOWN {
		source.Details().SetState(typex.SOURCE_DOWN)
		//----------------------------------
		// 当资源挂了以后先给停止, 然后重启
		//----------------------------------
		log.Warnf("Source %v %v down. try to restart it", source.Details().UUID, source.Details().Name)
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

//
//
//
func startSource(source typex.XSource, e *RuleEngine) error {
	//----------------------------------
	// 检查资源 如果是启动的，先给停了
	//----------------------------------
	ctx, cancelCTX := typex.NewCCTX()

	if err := source.Start(typex.CCTX{Ctx: ctx, CancelCTX: cancelCTX}); err != nil {
		log.Error("Source start error:", err)
		if source.Status() == typex.SOURCE_UP {
			source.Stop()
		}
		if source.Driver() != nil {
			if source.Driver().State() == typex.DRIVER_RUNNING {
				source.Driver().Stop()
			}
		}
		return err
	}
	//----------------------------------
	// 驱动也要停了, 然后重启
	//----------------------------------
	if source.Driver() != nil {
		if source.Driver().State() == typex.DRIVER_RUNNING {
			source.Driver().Stop()
		}
		// Start driver
		if err := source.Driver().Init(map[string]string{}); err != nil {
			log.Error("Driver initial error:", err)
			return errors.New("Driver initial error:" + err.Error())
		}
		log.Infof("Try to start driver: [%v]", source.Driver().DriverDetail().Name)
		if err := source.Driver().Work(); err != nil {
			log.Error("Driver work error:", err)
			return errors.New("Driver work error:" + err.Error())
		}
		log.Infof("Driver start successfully: [%v]", source.Driver().DriverDetail().Name)
	}
	return nil
}
