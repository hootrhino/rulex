package test

import (
	"net/url"
	"os"
	"os/signal"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hootrhino/rulex/glogger"
)

func Test_WS(t *testing.T) {
	ws()
}
func ws() {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "localhost:2580", Path: "/ws"}
	glogger.GLogger.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		glogger.GLogger.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				glogger.GLogger.Println("read:", err)
				return
			}
			glogger.GLogger.Printf("Client recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			ticker.Stop()
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				glogger.GLogger.Println("write:", err)
				return
			}
		case <-interrupt:
			glogger.GLogger.Println("interrupt")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				glogger.GLogger.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return

		}
		<-ticker.C
	}
}
