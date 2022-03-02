package target

import (
	"fmt"
	"rulex/core"
	"rulex/rulexrpc"
	"rulex/typex"
	"rulex/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var _state typex.SourceState

type _codecTargetConfig struct {
	Host string `json:"host" validate:"required"`
	Port int    `json:"port" validate:"required"`
	Type string `json:"type" validate:"required"`
}

type codecTarget struct {
	typex.XStatus
	host          string
	port          int
	_type         string
	client        rulexrpc.CodecClient
	rpcConnection *grpc.ClientConn
}

func NewCodecTarget(rx typex.RuleX) typex.XTarget {
	_state = typex.DOWN
	ct := &codecTarget{}
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
func (ct *codecTarget) Init(outEndId string, config map[string]interface{}) error {
	var mainConfig _codecTargetConfig
	if err := utils.BindSourceConfig(config, &mainConfig); err != nil {
		return err
	}

	ct.PointId = outEndId
	ct.host = mainConfig.Host
	ct.port = mainConfig.Port
	ct._type = mainConfig.Type

	return nil
}

//
// 启动资源
//
func (ct *codecTarget) Start(cctx typex.CCTX) error {
	ct.Ctx = cctx.Ctx
	ct.CancelCTX = cctx.CancelCTX
	rpcConnection, err := grpc.Dial(fmt.Sprintf("%s:%d", ct.host, ct.port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	ct.rpcConnection = rpcConnection
	ct.client = rulexrpc.NewCodecClient(rpcConnection)
	_state = typex.UP
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
	return _state

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
	return core.GenOutConfig(typex.GRPC_CODEC_TARGET, "GRPC_CODEC_TARGET", httpConfig{})

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
	if ct._type == "DECODE" {
		response, err = ct.client.Decode(ct.Ctx, dataRequest)
	} else if ct._type == "ENCODE" {
		response, err = ct.client.Encode(ct.Ctx, dataRequest)
	} else {
		_state = typex.DOWN
		return nil, fmt.Errorf("unknown operate type:%s", ct._type)
	}
	if err != nil {
		_state = typex.DOWN
		return nil, err

	}
	_state = typex.UP
	return response.GetData(), nil
}

//
// 不经过规则引擎处理的直达数据
//
func (ct *codecTarget) OnStreamApproached(data string) error {
	return nil
}

//
// 停止资源, 用来释放资源
//
func (ct *codecTarget) Stop() {
	ct.rpcConnection.Close()
}
