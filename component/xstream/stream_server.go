package xstream

import (
	"net"
	"time"

	"github.com/hootrhino/rulex/glogger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

/*
* TODO: V1.5.0
* 该模块准备被设计成一个分布式同步工具，用来将RULEX集群，节点之间可以互相感知，后期甚至可以互相调用接口
*
 */
type xStreamServer struct {
}

func (xs *xStreamServer) mustEmbedUnimplementedXStreamServer() {}

func (xs *xStreamServer) OnApproached(s XStream_OnApproachedServer) error {
	for {
		var r Response
		if err := s.RecvMsg(&r); err != nil {
			glogger.GLogger.Error(err)
			return err
		}
	}
}
func (xs *xStreamServer) SendStream(req *Request, s XStream_SendStreamServer) error {
	return nil
}
func ServerOptions() []grpc.ServerOption {
	var params = keepalive.ServerParameters{
		MaxConnectionIdle:     15 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
		MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
		MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
		Time:                  5 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
		Timeout:               1 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
	}
	return []grpc.ServerOption{
		grpc.KeepaliveParams(params),
	}
}

func StartXStreamServer() {
	server := grpc.NewServer(ServerOptions()...)
	RegisterXStreamServer(server, &xStreamServer{})
	listener, err := net.Listen("tcp", ":9999")
	if err != nil {
		glogger.GLogger.Fatal(err)
	}
	server.Serve(listener)
}
func StartXStreamClient() {
	var kacp = keepalive.ClientParameters{
		Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
		PermitWithoutStream: true,             // send pings even without active streams
	}
	conn, err := grpc.Dial("127.0.0.1:9999", grpc.WithKeepaliveParams(kacp))
	if err != nil {
		glogger.GLogger.Fatal(err)
	}
	defer conn.Close()

}
