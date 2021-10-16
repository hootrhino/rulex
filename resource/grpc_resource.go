package resource

import (
	"context"
	"fmt"
	"net"
	"rulex/rulexrpc"
	"rulex/typex"

	"github.com/ngaut/log"
	"google.golang.org/grpc"
)

const (
	DEFAULT_PORT      = ":2583"
	DEFAULT_TRANSPORT = "tcp"
)

type RulexRpcServer struct {
	grpcInEndResource *GrpcInEndResource
	rulexrpc.UnimplementedRulexRpcServer
}

//
func (r *RulexRpcServer) Work(ctx context.Context, in *rulexrpc.Data) (*rulexrpc.Response, error) {
	r.grpcInEndResource.RuleEngine.Work(
		r.grpcInEndResource.RuleEngine.GetInEnd(r.grpcInEndResource.PointId),
		in.Value,
	)
	return &rulexrpc.Response{
		Code:    0,
		Message: "OK",
	}, nil
}

//
// Resource interface
//
type GrpcInEndResource struct {
	typex.XStatus
	rulexServer *RulexRpcServer
	rpcServer   *grpc.Server
}

//
func NewGrpcInEndResource(inEndId string, e typex.RuleX) *GrpcInEndResource {
	h := GrpcInEndResource{}
	h.PointId = inEndId
	h.RuleEngine = e
	return &h
}

//
func (g *GrpcInEndResource) Start() error {
	config := g.RuleEngine.GetInEnd(g.PointId).Config
	var port = ""
	switch (*config)["port"].(type) {
	case string:
		port = ":" + (*config)["port"].(string)
		break
	case int:
		port = fmt.Sprintf(":%v", (*config)["port"].(int))
		break
	case int64:
		port = fmt.Sprintf(":%v", (*config)["port"].(int))
		break
	case float64:
		port = fmt.Sprintf(":%v", (*config)["port"].(int))
		break
	}
	// transport
	// TCP SSL HTTP HTTPS
	// transPort := (*config)["port"].(string)
	var err error
	var listener net.Listener
	if port == "" {
		listener, err = net.Listen(DEFAULT_TRANSPORT, DEFAULT_PORT)

	} else {
		listener, err = net.Listen(DEFAULT_TRANSPORT, port)
	}
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
func (g *GrpcInEndResource) DataModels() *map[string]typex.XDataModel {
	return &map[string]typex.XDataModel{}
}

//
func (g *GrpcInEndResource) Stop() {
	g.rpcServer.Stop()

}
func (g *GrpcInEndResource) Reload() {

}
func (g *GrpcInEndResource) Pause() {

}
func (g *GrpcInEndResource) Status() typex.ResourceState {
	return typex.UP
}

func (g *GrpcInEndResource) Register(inEndId string) error {
	g.PointId = inEndId
	return nil
}

func (g *GrpcInEndResource) Test(inEndId string) bool {
	return true
}

func (g *GrpcInEndResource) Enabled() bool {
	return true
}

func (g *GrpcInEndResource) Details() *typex.InEnd {
	return g.RuleEngine.GetInEnd(g.PointId)
}
func (m *GrpcInEndResource) OnStreamApproached(data string) error {
	return nil
}