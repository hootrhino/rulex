package glogger

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var GRealtimeLogger *RealTimeLogger
var lock sync.Mutex = sync.Mutex{}

/*
*
* 这是用来给外部输出日志的websocket服务器，其功能非常简单，就是单纯的对外输出实时日志，方便调试使用
* 注意：该功能需要配合HttpApiServer使用
*
 */
type RealTimeLogger struct {
	WsServer websocket.Upgrader
	Clients  map[string]*websocket.Conn
}

func (w *RealTimeLogger) Write(p []byte) (n int, err error) {
	for _, c := range w.Clients {
		c.WriteMessage(websocket.BinaryMessage, p)
	}
	return 0, nil
}

func StartNewRealTimeLogger(s string) *RealTimeLogger {
	GRealtimeLogger = &RealTimeLogger{
		WsServer: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		Clients: make(map[string]*websocket.Conn),
	}
	return GRealtimeLogger
}

/*
*
* 启动服务
*
 */

func WsLogger(c *gin.Context) {
	//upgrade get request to websocket protocol
	wsConn, err := GRealtimeLogger.WsServer.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	// 首先读第一个包是不是: WsTerminal
	wsConn.SetReadDeadline(time.Now().Add(5 * time.Second))
	_, b, err := wsConn.ReadMessage()
	if err != nil {
		//
		return
	}
	wsConn.SetReadDeadline(time.Time{})
	token := string(b)
	if token != "WsTerminal" {
		wsConn.WriteMessage(1, []byte("Invalid client token"))
		wsConn.Close()
		return
	}
	// 最多允许连接10个客户端，实际情况下根本用不了那么多
	if len(GRealtimeLogger.Clients) >= 10 {
		wsConn.WriteMessage(1, []byte("Reached max connections"))
		wsConn.Close()
		return
	}
	GRealtimeLogger.Clients[wsConn.RemoteAddr().String()] = wsConn
	wsConn.WriteMessage(1, []byte("Connected"))
	go func(ctx context.Context, wsConn *websocket.Conn) {
		for {
			select {
			case <-ctx.Done():
				{
					return
				}
			default:
				{
				}
			}
			// 当前不需要相互交互，单向给Websocket发送日志就行
			// 因此这里只需要判断下是否掉线即可

			_, _, err := wsConn.ReadMessage()
			// wsConn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				wsConn.Close()
				lock.Lock()
				delete(GRealtimeLogger.Clients, wsConn.RemoteAddr().String())
				lock.Unlock()
				return
			}

		}

	}(context.Background(), wsConn)
}
