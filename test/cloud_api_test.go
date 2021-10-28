package test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	cloud "rulex/cloud"
)

// var (
// 	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:9090", "gRPC server endpoint")
// )

type AtomicCloudServiceServer struct {
	cloud.UnsafeAtomicCloudServiceServer
}

func (s *AtomicCloudServiceServer) CallCloud(ctx context.Context, cs *cloud.Service) (*cloud.CallResult, error) {
	fmt.Println("----", cs)
	return nil, nil
}
func run() error {
	server := grpc.NewServer()
	ctx := context.Background()
	_, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux()
	// opts := []grpc.DialOption{grpc.WithInsecure()}
	// cloud.RegisterAtomicCloudServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	cloud.RegisterAtomicCloudServiceServer(server, &AtomicCloudServiceServer{})
	return http.ListenAndServe(":8089", mux)
}

func TestAtomicCloudService(t *testing.T) {
	// flag.Parse()
	// defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
