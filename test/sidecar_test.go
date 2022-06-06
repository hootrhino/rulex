package test

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"rulex/sidecar"
	"testing"
)

func fork() {

}
func Test_sidecar_client(t *testing.T) {
	// Unix domain
	conn, err := grpc.Dial("/usr/sidecar.sock",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Errorf("grpc.Dial err: %v", err)
	}
	client := sidecar.NewSidecarClient(conn)

	client.Init(context.Background(), &sidecar.Config{
		Kv: map[string]string{
			"K1": "V1",
			"K2": "V2",
			"K3": "V3",
		},
	})
	defer conn.Close()
}
