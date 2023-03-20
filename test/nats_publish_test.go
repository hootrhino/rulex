package test

import (
	"testing"
	"time"

	"github.com/nats-io/nats.go"
)

func TestOk(t *testing.T) {
	// Connect to a server
	connection, err := nats.Connect("127.0.0.1:4222", func(o *nats.Options) error {
		o.User = "nats_client"
		o.Password = "******"
		return nil
	})
	// connection.Subscribe("downstream.services.publish", func(m *nats.Msg) {
	// 	t.Logf("Received a message =<<<<<< %s\n", string(m.Data))
	// })
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Publish =>>>>>>> Hello World")
	connection.Publish("downstream.services.publish", []byte("Hello World"))
	t.Log("Publish Ok.")
	time.Sleep(5 * time.Second)
	connection.Drain()
	connection.Close()
}
