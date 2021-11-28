package test

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
)

/*
*
* 监听器
*
 */
func Start_nats_listener() {
	connection, err := nats.Connect("127.0.0.1:4222", func(o *nats.Options) error {
		o.User = "nats_client"
		o.Password = "123456"
		return nil
	})
	if err != nil {
		fmt.Printf("Error:%v", err)
		return
	}
	//
	// {
	//     "type": "finishCmd",
	//     "cmdId":"112233...."
	// }
	//
	connection.Subscribe("upstream.devices.state", func(m *nats.Msg) {
		fmt.Printf("state: %s", string(m.Data))
	})

	defer connection.Close()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGABRT)
	signal := <-c
	fmt.Printf("Received stop signal:%v", signal)
	connection.Drain()
	os.Exit(1)
}
