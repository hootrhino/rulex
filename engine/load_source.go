package engine

import (
	"context"
	"fmt"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
)

/*
*
* 加载用户资源
*
 */
func (e *RuleEngine) LoadUserInEnd(source typex.XSource, in *typex.InEnd) error {
	return nil
}

/*
*
* 内建资源
*
 */
func (e *RuleEngine) LoadBuiltInEnd(in *typex.InEnd) error {
	return nil
}

/*
*
* 加载输入资源
*
 */
func (e *RuleEngine) LoadInEndWithCtx(in *typex.InEnd,
	ctx context.Context, cancelCTX context.CancelFunc) error {
	if config := e.SourceTypeManager.Find(in.Type); config != nil {
		return e.loadSource(config.NewSource(e), in, ctx, cancelCTX)
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
func (e *RuleEngine) loadSource(source typex.XSource, in *typex.InEnd,
	ctx context.Context, cancelCTX context.CancelFunc) error {
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

	e.startSource(source, ctx, cancelCTX)
	glogger.GLogger.Infof("InEnd [%v, %v] load successfully", in.Name, in.UUID)
	return nil
}

func (e *RuleEngine) startSource(source typex.XSource,
	ctx context.Context, cancelCTX context.CancelFunc) error {

	if err := source.Start(typex.CCTX{Ctx: ctx, CancelCTX: cancelCTX}); err != nil {
		glogger.GLogger.Error("Source start error:", err)
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
