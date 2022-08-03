package target

import (
	"fmt"

	"github.com/i4de/rulex/common"
	"github.com/i4de/rulex/core"
	"github.com/i4de/rulex/rulexrpc"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type codecTarget struct {
	typex.XStatus
	client        rulexrpc.CodecClient
	rpcConnection *grpc.ClientConn
	mainConfig    common.GrpcConfig
	status        typex.SourceState
}

func NewCodecTarget(rx typex.RuleX) typex.XTarget {
	ct := &codecTarget{}
	ct.mainConfig = common.GrpcConfig{}
	ct.RuleEngine = rx
	return ct
}

//
// 测试资源是否可用
//
func (ct *codecTarget) Test(outEndId string) bool {
	return true
}

//
// 用来初始化传递资源配置
//
func (ct *codecTarget) Init(outEndId string, configMap map[string]interface{}) error {
	ct.PointId = outEndId
	//
	if err := utils.BindSourceConfig(configMap, &ct.mainConfig); err != nil {
		return err
	}

	return nil
}

//
// 启动资源
//
func (ct *codecTarget) Start(cctx typex.CCTX) error {
	ct.Ctx = cctx.Ctx
	ct.CancelCTX = cctx.CancelCTX
	//
	rpcConnection, err := grpc.Dial(fmt.Sprintf("%s:%d", ct.mainConfig.Host, ct.mainConfig.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	ct.rpcConnection = rpcConnection
	ct.client = rulexrpc.NewCodecClient(rpcConnection)
	ct.status = typex.SOURCE_UP
	return nil

}

//
// 资源是否被启用
//
func (ct *codecTarget) Enabled() bool {
	return true
}

//
// 重载: 比如可以在重启的时候把某些数据保存起来
//
func (ct *codecTarget) Reload() {

}

//
// 挂起资源, 用来做暂停资源使用
//
func (ct *codecTarget) Pause() {

}

//
// 获取资源状态
//
func (ct *codecTarget) Status() typex.SourceState {
	return ct.status

}

//
// 获取资源绑定的的详情
//
func (ct *codecTarget) Details() *typex.OutEnd {
	out := ct.RuleEngine.GetOutEnd(ct.PointId)
	return out

}

//
//
//
func (ct *codecTarget) Configs() *typex.XConfig {
	return core.GenOutConfig(typex.GRPC_CODEC_TARGET, "GRPC_CODEC_TARGET", common.GrpcConfig{})

}

//
// 数据出口
//
func (ct *codecTarget) To(data interface{}) (interface{}, error) {
	dataRequest := &rulexrpc.CodecRequest{
		Value: []byte(data.(string)),
	}
	var response *rulexrpc.CodecResponse
	var err error
	if ct.mainConfig.Type == "DECODE" {
		response, err = ct.client.Decode(ct.Ctx, dataRequest)
	}
	if ct.mainConfig.Type == "ENCODE" {
		response, err = ct.client.Encode(ct.Ctx, dataRequest)
	}
	if err != nil {
		return nil, err
	}
	return response.GetData(), nil
}

//
// 停止资源, 用来释放资源
//
func (ct *codecTarget) Stop() {
	ct.rpcConnection.Close()
	ct.status = typex.SOURCE_STOP
}
