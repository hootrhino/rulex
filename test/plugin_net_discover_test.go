package test

import (
	"fmt"
	"net"
	"time"

	netdiscover "github.com/hootrhino/rulex/plugin/net_discover"

	"testing"
)

// go test -timeout 30s -run ^Test_net_discover github.com/hootrhino/rulex/test -v -count=1

func Test_net_discover(t *testing.T) {
	engine := RunTestEngine()
	engine.Start()

	hh := netdiscover.NewNetDiscover()
	// HttpApiServer loaded default
	if err := engine.LoadPlugin("plugin.netdiscover", hh); err != nil {
		t.Fatal("Rule load failed:", err)
	}
	for i := 0; i < 5; i++ {
		test_client()
		time.Sleep(1 * time.Second)
	}
	time.Sleep(3 * time.Second)
	engine.Stop()
}

func test_client() {

	addr, _ := net.ResolveUDPAddr("udp4", "0.0.0.0:1994")
	socket, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		panic(err)
	}
	defer socket.Close()
	sendData := []byte("NODE_INFO")
	_, err = socket.Write(sendData)
	if err != nil {
		panic(err)
	}
	data := make([]byte, 1024)
	// socket.SetReadDeadline(time.Now().Add(time.Duration(time.Second * 5)))
	n, remoteAddr, err := socket.ReadFromUDP(data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("recv:%v addr:%v count:%v\n", string(data[:n]), remoteAddr, n)
}
