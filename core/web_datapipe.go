// Copyright (C) 2023 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
)

// 前端管道
var DefaultWebDataPipe *WebDataPipe

/*
*
* Websocket 管道
*
 */
type WebDataPipe struct {
	rulex    typex.RuleX
	WsServer websocket.Upgrader
	Clients  map[string]*websocket.Conn
	lock     sync.Mutex
}

/*
*
* 初始化
*
 */
func InitWebDataPipe(rulex typex.RuleX) *WebDataPipe {
	DefaultWebDataPipe = &WebDataPipe{
		rulex:   rulex,
		Clients: map[string]*websocket.Conn{},
		lock:    sync.Mutex{},
		WsServer: websocket.Upgrader{
			EnableCompression: true,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
	return DefaultWebDataPipe
}

/*
*
* 启动, 默认在3580端口
*
 */
func StartWebDataPipe() error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.GET("/ws", dataPipLoop)
	glogger.GLogger.Info("WebDataPipe started on: 0.0.0.0:3580")
	router.Run(":3580")
	return nil
}

/*
*
* Websocket
*
 */
func dataPipLoop(c *gin.Context) {
	wsConn, err := DefaultWebDataPipe.WsServer.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	wsConn.SetReadDeadline(time.Now().Add(5 * time.Second))
	_, b, err := wsConn.ReadMessage()
	if err != nil {
		return
	}
	wsConn.SetReadDeadline(time.Time{})
	token := string(b)
	if token != "WebDataPipe" {
		wsConn.WriteMessage(1, []byte("Invalid client token"))
		wsConn.Close()
		return
	}
	// 最多允许连接10个客户端，实际情况下根本用不了那么多
	if len(DefaultWebDataPipe.Clients) >= 10 {
		wsConn.WriteMessage(websocket.TextMessage, []byte("Reached max client connections"))
		wsConn.Close()
		return
	}
	DefaultWebDataPipe.Clients[wsConn.RemoteAddr().String()] = wsConn
	wsConn.WriteMessage(websocket.TextMessage, []byte("Connected"))
	glogger.GLogger.Info("DefaultWebDataPipe Terminal connected:" + wsConn.RemoteAddr().String())
	wsConn.SetCloseHandler(func(code int, text string) error {
		glogger.GLogger.Info("DefaultWebDataPipe CloseHandler:", wsConn.RemoteAddr().String())
		return nil
	})
	// ping
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
			_, _, err1 := wsConn.ReadMessage()
			err2 := wsConn.WriteMessage(websocket.PingMessage, []byte{})
			if err1 != nil || err2 != nil {
				glogger.GLogger.Error("DefaultWebDataPipe error:",
					wsConn.RemoteAddr().String(), ", Error:", func(e1, e2 error) error {
						if e1 != nil {
							return e1
						}
						if e2 != nil {
							return e2
						}
						return nil
					}(err1, err2))
				wsConn.Close()
				DefaultWebDataPipe.lock.Lock()
				delete(DefaultWebDataPipe.Clients, wsConn.RemoteAddr().String())
				DefaultWebDataPipe.lock.Unlock()
				return
			}
		}

	}(context.Background(), wsConn)
	// Read
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
			Type, Data, err := wsConn.ReadMessage()
			if err != nil {
				glogger.GLogger.Error("DefaultWebDataPipe error:",
					wsConn.RemoteAddr().String(), ", Error:", err)
				wsConn.Close()
				DefaultWebDataPipe.lock.Lock()
				delete(DefaultWebDataPipe.Clients, wsConn.RemoteAddr().String())
				DefaultWebDataPipe.lock.Unlock()
				return
			}

			// TODO
			// 对前端的要求是Text数据,解码成一套事件系统,处理交互组件 但是这块估计到0.9以后再考虑了
			//
			glogger.GLogger.Info("DefaultWebDataPipe Receive Data From UI:",
				wsConn.RemoteAddr().String(), Type, Data)
		}
	}(context.Background(), wsConn)
}
