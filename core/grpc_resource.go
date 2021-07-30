package core

import (
	"context"
	"net"
	"rulex/rulexrpc"

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
	r.grpcInEndResource.ruleEngine.Work(
		r.grpcInEndResource.ruleEngine.GetInEnd(r.grpcInEndResource.PointId),
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
	XStatus
	rulexServer *RulexRpcServer
	rpcServer   *grpc.Server
}

//
func NewGrpcInEndResource(inEndId string, e *RuleEngine) *GrpcInEndResource {
	h := GrpcInEndResource{}
	h.PointId = inEndId
	h.ruleEngine = e
	return &h
}

//
func (g *GrpcInEndResource) Start() error {
	config := g.ruleEngine.GetInEnd(g.PointId).Config
	port := ":" + (*config)["port"].(string)
	// transport
	// TCP SSL HTTP HTTPS
	// transPort := (*config)["port"].(string)
	var err error
	var listener net.Listener
	if port == "" {
		listener, err = net.Listen(DEFAULT_TRANSPORT, DEFAULT_PORT)

	} else {
		listener, err = net.Listen(DEFAULT_TRANSPORT, ":"+(*config)["port"].(string))
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
	go func(context.Context) {
		log.Info("GrpcInEndResource Started At:", listener.Addr())
		g.rpcServer.Serve(listener)
	}(context.Background())
	return nil
}

//
func (g *GrpcInEndResource) DataModels() *map[string]XDataModel {
	return &map[string]XDataModel{}
}

//
func (g *GrpcInEndResource) Stop() {

}
func (g *GrpcInEndResource) Reload() {

}
func (g *GrpcInEndResource) Pause() {

}
func (g *GrpcInEndResource) Status() State {
	return g.ruleEngine.GetInEnd(g.PointId).State
}

func (g *GrpcInEndResource) Register(inEndId string) error {
	g.PointId = inEndId
	return nil
}

func (g *GrpcInEndResource) Test(inEndId string) bool {
	return true
}

func (g *GrpcInEndResource) Enabled() bool {
	return g.Enable
}
