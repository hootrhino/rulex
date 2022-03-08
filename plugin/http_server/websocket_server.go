package httpserver

import (
	"net/http"
	"rulex/typex"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ngaut/log"
)

var upgrader = websocket.Upgrader{}

func WS(c *gin.Context, hh *HttpApiServer, e typex.RuleX) {
	ws(c.Writer, c.Request)
}
func ws(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}
	defer conn.Close()
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Error(err)
			break
		}
		log.Debugf("Server recv: %s", message)
		err = conn.WriteMessage(mt, []byte("test"))
		if err != nil {
			log.Error(err)
			break
		}
	}
}
