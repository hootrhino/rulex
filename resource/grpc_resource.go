package resource

import (
	"context"
	"fmt"
	"net"
	"rulex/core"
	"rulex/rulexrpc"
	"rulex/typex"
	"rulex/utils"

	"github.com/ngaut/log"
	"google.golang.org/grpc"
)

const (
	DEFAULT_TRANSPORT = "tcp"
)

//
type grpcConfig struct {
	Port       uint16             `json:"port" validate:"required" title:"端口" info:""`
	DataModels []typex.XDataModel `json:"dataModels" title:"数据模型" info:""`
}

type RulexRpcServer struct {
	grpcInEndResource *grpcInEndResource
	rulexrpc.UnimplementedRulexRpcServer
}

//
// Resource interface
//
type grpcInEndResource struct {
	typex.XStatus
	rulexServer *RulexRpcServer
	rpcServer   *grpc.Server
}

//
func NewGrpcInEndResource(inEndId string, e typex.RuleX) typex.XResource {
	g := grpcInEndResource{}
	g.PointId = inEndId
	g.RuleEngine = e
	return &g
}

//
func (g *grpcInEndResource) Start() error {
	inEnd := g.RuleEngine.GetInEnd(g.PointId)
	config := inEnd.Config
	var mainConfig grpcConfig
	if err := utils.BindResourceConfig(config, &mainConfig); err != nil {
		return err
	}

	listener, err := net.Listen(DEFAULT_TRANSPORT, fmt.Sprintf(":%d", mainConfig.Port))
	if err != nil {
		return err
	}
	// Important !!!
	g.rpcServer = grpc.NewServer()
	g.rulexServer = new(RulexRpcServer)
	g.rulexServer.grpcInEndResource = g
	//
	rulexrpc.RegisterRulexRpcServer(g.rpcServer, g.rulexServer)
	go func(c context.Context) {
		log.Info("RulexRpc resource started on", listener.Addr())
		g.rpcServer.Serve(listener)
	}(context.Background())

	return nil
}

//
func (g *grpcInEndResource) DataModels() []typex.XDataModel {
	return []typex.XDataModel{}
}

//
func (g *grpcInEndResource) Stop() {
	if g.rpcServer != nil {
		g.rpcServer.Stop()
	}

}
func (g *grpcInEndResource) Reload() {

}
func (g *grpcInEndResource) Pause() {

}
func (g *grpcInEndResource) Status() typex.ResourceState {
	return typex.UP
}

func (g *grpcInEndResource) Register(inEndId string) error {
	g.PointId = inEndId
	return nil
}

func (g *grpcInEndResource) Test(inEndId string) bool {
	return true
}

func (g *grpcInEndResource) Enabled() bool {
	return true
}

func (g *grpcInEndResource) Details() *typex.InEnd {
	return g.RuleEngine.GetInEnd(g.PointId)
}
func (m *grpcInEndResource) OnStreamApproached(data string) error {
	return nil
}
func (*grpcInEndResource) Driver() typex.XExternalDriver {
	return nil
}
func (*grpcInEndResource) Configs() typex.XConfig {
	config, err := core.RenderConfig("GRPC", "", grpcConfig{})
	if err != nil {
		log.Error(err)
		return typex.XConfig{}
	} else {
		return config
	}
}

//
func (r *RulexRpcServer) Work(ctx context.Context, in *rulexrpc.Data) (*rulexrpc.Response, error) {
	ok, err := r.grpcInEndResource.RuleEngine.Work(
		r.grpcInEndResource.RuleEngine.GetInEnd(r.grpcInEndResource.PointId),
		in.Value,
	)
	if ok {
		return &rulexrpc.Response{
			Code:    0,
			Message: "OK",
		}, nil
	} else {
		return &rulexrpc.Response{
			Code:    1,
			Message: err.Error(),
		}, err
	}

}

//
// 拓扑
//
func (*grpcInEndResource) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}
