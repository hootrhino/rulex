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
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/hootrhino/rulex/component/interqueue"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
)

/*
*
* WebsocketDataPipe 主要用来解决大屏的数据交互以及数据推送问题，不涉及别的业务
*
 */
var __DefaultWebDataPipe *WebsocketDataPipe

/*
*
* Websocket 管道
*
 */
type WebsocketDataPipe struct {
	rulex    typex.RuleX
	WsServer websocket.Upgrader
	Clients  map[string]*websocket.Conn
	lock     sync.Mutex
}

/*
*
* 初始化缓冲器
*
 */
func InitWebDataPipe(rulex typex.RuleX) *WebsocketDataPipe {
	__DefaultWebDataPipe = &WebsocketDataPipe{
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
	return __DefaultWebDataPipe
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

	/*
	*
	*从管道里面拿写到前端的数据
	*
	 */
	go func(ctx context.Context, WebsocketDataPipe *WebsocketDataPipe) {
		for {
			select {
			case <-ctx.Done():
				return
			case d := <-interqueue.OutQueue():
				{
					glogger.GLogger.Debug("DefaultInteractQueue OutQueue:", d.String())
					for _, wsClient := range WebsocketDataPipe.Clients {
						wsClient.WriteMessage(websocket.TextMessage, []byte(d.String()))
					}
				}
			}
		}
	}(typex.GCTX, __DefaultWebDataPipe)
	go func(ctx context.Context, WebsocketDataPipe *WebsocketDataPipe) {
		for {
			select {
			case <-ctx.Done():
				return
			case d := <-interqueue.InQueue():
				{
					//
					// TODO 交互事件数据
					// v0.8 规划
					glogger.GLogger.Debug("DefaultInteractQueue InQueue:", d.String())
				}
			}
		}
	}(typex.GCTX, __DefaultWebDataPipe)
	glogger.GLogger.Info("WebsocketDataPipe started on: 0.0.0.0:2579")
	router.Run(":2579")
	return nil
}

/*
*
* Websocket
*
 */
func dataPipLoop(c *gin.Context) {
	wsConn, err := __DefaultWebDataPipe.WsServer.Upgrade(c.Writer, c.Request, nil)
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
	if token != "WebsocketDataPipe" {
		wsConn.WriteMessage(1, []byte("Invalid client token"))
		wsConn.Close()
		return
	}
	// 最多允许连接10个客户端，实际情况下根本用不了那么多
	if len(__DefaultWebDataPipe.Clients) > 5 {
		wsConn.WriteMessage(websocket.TextMessage, []byte("Reached max client connections"))
		wsConn.Close()
		return
	}
	__DefaultWebDataPipe.Clients[wsConn.RemoteAddr().String()] = wsConn
	wsConn.WriteMessage(websocket.TextMessage, []byte("Connected"))
	glogger.GLogger.Info("__DefaultWebDataPipe Terminal connected:" + wsConn.RemoteAddr().String())
	wsConn.SetCloseHandler(func(code int, text string) error {
		glogger.GLogger.Info("wsConn CloseHandler:", wsConn.RemoteAddr().String())
		__DefaultWebDataPipe.lock.Lock()
		delete(__DefaultWebDataPipe.Clients, wsConn.RemoteAddr().String())
		__DefaultWebDataPipe.lock.Unlock()
		return nil
	})
	wsConn.SetPingHandler(func(appData string) error {
		return nil
	})
	wsConn.SetPongHandler(func(appData string) error {
		return nil
	})
	/*
	*
	* 来自前端的事件
	*
	 */
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
				glogger.GLogger.Error("__DefaultWebDataPipe error:",
					wsConn.RemoteAddr().String(), ", Error:", err)
				wsConn.Close()
				__DefaultWebDataPipe.lock.Lock()
				delete(__DefaultWebDataPipe.Clients, wsConn.RemoteAddr().String())
				__DefaultWebDataPipe.lock.Unlock()
				return
			}
			if Type == websocket.BinaryMessage {
				uiData := interqueue.InteractQueueData{}
				if err := json.Unmarshal(Data, &uiData); err != nil {
					glogger.GLogger.Error(err)
					continue
				}
				interqueue.ReceiveData(uiData)
			} else {
				glogger.GLogger.Error(fmt.Errorf("Message type not support:%v", Type))
			}
		}
	}(context.Background(), wsConn)

}
