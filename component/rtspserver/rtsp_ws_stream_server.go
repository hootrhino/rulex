// Copyright (C) 2024 wwhai
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

package rtspserver

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
)

/*
*
* websocket播放器记录
*
 */
type streamPlayer struct {
	wsConn        *websocket.Conn
	GetFirstFrame bool
	Token         string
	Type          string           `json:"type"` // push | pull
	LiveId        string           `json:"liveId"`
	Pulled        bool             `json:"pulled"`
	Resolution    utils.Resolution `json:"resolution"`
}

/*
*
* 启动服务
*
 */
func wsServerEndpoint(c *gin.Context) {
	wsConn, err := __DefaultRtspServer.WsServer.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	LiveId := c.Query("liveId")
	Token := c.Query("token")
	glogger.GLogger.Debugf("Request live:%s, Token is :%s", LiveId, Token)
	// 最多允许连接10个客户端，实际情况下根本用不了那么多
	if len(__DefaultRtspServer.Clients) >= 10 {
		wsConn.WriteMessage(websocket.CloseMessage, []byte{})
		wsConn.Close()
		return
	}

	streamPlayer := streamPlayer{
		LiveId: LiveId,
		Token:  Token,
		wsConn: wsConn,
	}
	if Token != "WebRtspPlayer" {
		msg := "Invalid client token"
		glogger.GLogger.Error(msg)
		wsConn.WriteMessage(websocket.CloseMessage, []byte(msg))
		wsConn.Close()
		return
	}
	if C, ok := __DefaultRtspServer.Clients[wsConn.RemoteAddr().String()]; ok {
		msg := "already exists a client:" + C.LiveId
		glogger.GLogger.Error(msg)
		wsConn.WriteMessage(websocket.CloseMessage, []byte(msg))
		wsConn.Close()
		return
	}
	// 考虑到性能，每个LiveId只能有1路播放, 后期优化成支持多路播放
	__DefaultRtspServer.Clients[wsConn.RemoteAddr().String()] = streamPlayer
	glogger.GLogger.Info("WebSocket Player connected:" + wsConn.RemoteAddr().String())
	wsConn.SetCloseHandler(func(code int, text string) error {
		glogger.GLogger.Info("WebSocket Close Handler:", wsConn.RemoteAddr().String())
		__DefaultRtspServer.locker.Lock()
		delete(__DefaultRtspServer.Clients, wsConn.RemoteAddr().String())
		__DefaultRtspServer.locker.Unlock()
		return nil
	})
	wsConn.SetPingHandler(func(appData string) error {
		return nil
	})
	wsConn.SetPongHandler(func(appData string) error {
		return nil
	})
	go func(wsConn *websocket.Conn) {
		for {
			select {
			case <-typex.GCTX.Done():
				{
					return
				}
			default:
				{
				}
			}
			_, _, err := wsConn.ReadMessage()
			if err != nil {
				glogger.GLogger.Warn("WebSocket Close Handler:", wsConn.RemoteAddr().String())
				__DefaultRtspServer.locker.Lock()
				delete(__DefaultRtspServer.Clients, wsConn.RemoteAddr().String())
				__DefaultRtspServer.locker.Unlock()
				break
			}
			err = wsConn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				break
			}
			time.Sleep(5 * time.Second)
		}
	}(wsConn)
}
