package test

import (
	"context"
	"testing"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
	"github.com/gopcua/opcua/uacp"
	"github.com/hootrhino/rulex/glogger"
)

func Test_opcua_read(t *testing.T) {
	ctx := context.Background()
	go startServer(ctx)
	startClient(ctx)
}
func startClient(ctx context.Context) {
	c := opcua.NewClient("opc.tcp://localhost:4840/foo/bar", opcua.SecurityMode(ua.MessageSecurityModeNone))
	if err := c.Connect(ctx); err != nil {
		glogger.GLogger.Fatal(err)
	}
	defer c.CloseSessionWithContext(ctx)

	req := &ua.ReadRequest{
		MaxAge:             2000,
		NodesToRead:        []*ua.ReadValueID{{NodeID: &ua.NodeID{}}},
		TimestampsToReturn: ua.TimestampsToReturnBoth,
	}

	resp, err := c.ReadWithContext(ctx, req)
	if err != nil {
		glogger.GLogger.Fatalf("Read failed: %s", err)
	}
	if resp.Results[0].Status != ua.StatusOK {
		glogger.GLogger.Fatalf("Status not OK: %v", resp.Results[0].Status)
	}
	glogger.GLogger.Printf("%#v", resp.Results[0].Value.Value())
}
func startServer(ctx context.Context) {
	endpoint := "opc.tcp://localhost:4840/foo/bar"
	glogger.GLogger.Printf("Listening on %s", endpoint)
	l, err := uacp.Listen(endpoint, nil)
	if err != nil {
		glogger.GLogger.Fatal(err)
	}
	c, err := l.Accept(ctx)
	if err != nil {
		glogger.GLogger.Fatal(err)
	}
	glogger.GLogger.Printf("conn %d: connection from %s", c.ID(), c.RemoteAddr())
}
