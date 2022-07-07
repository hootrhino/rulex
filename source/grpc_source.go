package source

import (
	"context"
	"fmt"
	"net"

	"github.com/i4de/rulex/core"
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/rulexrpc"
	"github.com/i4de/rulex/typex"
	"github.com/i4de/rulex/utils"

	"google.golang.org/grpc"
)

const (
	DefaultTransport = "tcp"
)

//
type grpcConfig struct {
	Port       uint16             `json:"port" validate:"required" title:"端口" info:""`
	DataModels []typex.XDataModel `json:"dataModels" title:"数据模型" info:""`
}

type RulexRpcServer struct {
	grpcInEndSource *grpcInEndSource
	rulexrpc.UnimplementedRulexRpcServer
}

//
// Source interface
//
type grpcInEndSource struct {
	typex.XStatus
	rulexServer *RulexRpcServer
	rpcServer   *grpc.Server
}

//
func NewGrpcInEndSource(inEndId string, e typex.RuleX) typex.XSource {
	g := grpcInEndSource{}
	g.PointId = inEndId
	g.RuleEngine = e
	return &g
}

//
func (g *grpcInEndSource) Start(cctx typex.CCTX) error {
	g.Ctx = cctx.Ctx
	g.CancelCTX = cctx.CancelCTX
	config := g.RuleEngine.GetInEnd(g.PointId).Config
	var mainConfig grpcConfig
	if err := utils.BindSourceConfig(config, &mainConfig); err != nil {
		return err
	}

	listener, err := net.Listen(DefaultTransport, fmt.Sprintf(":%d", mainConfig.Port))
	if err != nil {
		return err
	}
	// Important !!!
	g.rpcServer = grpc.NewServer()
	g.rulexServer = new(RulexRpcServer)
	g.rulexServer.grpcInEndSource = g
	//
	rulexrpc.RegisterRulexRpcServer(g.rpcServer, g.rulexServer)

	go func(c context.Context) {
		glogger.GLogger.Info("RulexRpc source started on", listener.Addr())
		g.rpcServer.Serve(listener)
	}(g.Ctx)

	return nil
}

//
func (g *grpcInEndSource) DataModels() []typex.XDataModel {
	return g.XDataModels
}

//
func (g *grpcInEndSource) Stop() {
	if g.rpcServer != nil {
		g.rpcServer.Stop()
	}
	g.CancelCTX()

}
func (g *grpcInEndSource) Reload() {

}
func (g *grpcInEndSource) Pause() {

}
func (g *grpcInEndSource) Status() typex.SourceState {
	return typex.SOURCE_UP
}

func (g *grpcInEndSource) Init(inEndId string, cfg map[string]interface{}) error {

	g.PointId = inEndId
	return nil
}
func (g *grpcInEndSource) Test(inEndId string) bool {
	return true
}

func (g *grpcInEndSource) Enabled() bool {
	return true
}

func (g *grpcInEndSource) Details() *typex.InEnd {
	return g.RuleEngine.GetInEnd(g.PointId)
}

func (*grpcInEndSource) Driver() typex.XExternalDriver {
	return nil
}
func (*grpcInEndSource) Configs() *typex.XConfig {
	return core.GenInConfig(typex.GRPC, "GRPC", grpcConfig{})
}

//
func (r *RulexRpcServer) Work(ctx context.Context, in *rulexrpc.Data) (*rulexrpc.Response, error) {
	ok, err := r.grpcInEndSource.RuleEngine.WorkInEnd(
		r.grpcInEndSource.RuleEngine.GetInEnd(r.grpcInEndSource.PointId),
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
func (*grpcInEndSource) Topology() []typex.TopologyPoint {
	return []typex.TopologyPoint{}
}

//
// 来自外面的数据
//
func (*grpcInEndSource) DownStream([]byte) {}

//
// 上行数据
//
func (*grpcInEndSource) UpStream() {}
