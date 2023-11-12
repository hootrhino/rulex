package rtspserver

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
)

/*
*
* 默认服务
*
 */
var __DefaultRtspServer *rtspServer

/*
*
* RTSP 设备(rtsp://192.168.199.243:554/av0_0)
*
 */
type RtspCameraInfo struct {
	Type     string `json:"type,omitempty"`     // 1-RTSP,2-Local
	LocalId  string `json:"local_id,omitempty"` // 本地ID
	PullAddr string `json:"pullAddr,omitempty"`
	PushAddr string `json:"pushAddr,omitempty"`
}

/*
*
* 这是用来给外部输出日志的websocket服务器，其功能非常简单，就是单纯的对外输出实时日志，方便调试使用。
* 注意：该功能需要配合HttpApiServer使用, 客户端连上以后必须在5s内发送一个 ‘WsPlayer’ 的固定字符
*       串到服务器才能过认证。
*
 */
type websocketPlayerManager struct {
	WsServer websocket.Upgrader
	Clients  map[string]*websocket.Conn
	lock     sync.Mutex
}

type rtspServer struct {
	webServer              *gin.Engine
	rtspCameras            map[string]RtspCameraInfo
	websocketPlayerManager *websocketPlayerManager
}

// NewRouter Gin 路由配置
func InitRtspServer() *rtspServer {
	gin.SetMode(gin.ReleaseMode)
	__DefaultRtspServer = &rtspServer{
		webServer:              gin.New(),
		rtspCameras:            map[string]RtspCameraInfo{},
		websocketPlayerManager: NewPlayerManager(),
	}
	// 注册Websocket server
	__DefaultRtspServer.webServer.GET("/ws", wsServerEndpoint)
	// http://127.0.0.1:3000/stream/live/001
	group := __DefaultRtspServer.webServer.Group("/stream")
	// 注意：这个接口是给FFMPEG请求的
	group.POST("/ffmpegPush", func(ctx *gin.Context) {
		LiveId := ctx.Query("liveId")
		// Token := ctx.Query("token")
		glogger.GLogger.Info("Try to load RTSP From:", LiveId)
		// http://127.0.0.1:9400 :后期通过参数传进
		// 启动一个FFMPEG开始从摄像头拉流
		bodyReader := bufio.NewReader(ctx.Request.Body)
		for {
			// data 就是 RTSP 帧
			// 只需将其转发给websocket即可
			data, err := bodyReader.ReadBytes('\n')
			if err != nil {
				break
			}
			pushToWebsocket(LiveId, data)
		}
		ctx.Writer.Flush()
	})
	go func(ctx context.Context) {
		listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", 9400))
		if err != nil {
			glogger.GLogger.Fatalf("Rtsp stream server listen error: %s\n", err)
		}
		defer listener.Close()
		if err := __DefaultRtspServer.webServer.RunListener(listener); err != nil {
			glogger.GLogger.Fatalf("Rtsp stream server listen error: %s\n", err)
		}
	}(context.Background())
	glogger.GLogger.Info("Rtsp stream server start success")
	return __DefaultRtspServer
}
func pushToWebsocket(liveId string, data []byte) {
	// fmt.Println(liveId, data)
	if C, Ok := __DefaultRtspServer.websocketPlayerManager.Clients[liveId]; Ok {
		C.WriteMessage(2, data)
	}
}

/*
*
* 远程摄像头列表
*
 */
func AllVideoStreamEndpoints() map[string]RtspCameraInfo {
	return __DefaultRtspServer.rtspCameras
}
func AddVideoStreamEndpoint(k string, v RtspCameraInfo) {
	if GetVideoStreamEndpoint(k).PullAddr == "" {
		__DefaultRtspServer.rtspCameras[k] = v
	}
}
func GetVideoStreamEndpoint(k string) RtspCameraInfo {
	return __DefaultRtspServer.rtspCameras[k]
}
func DeleteVideoStreamEndpoint(k string) {
	delete(__DefaultRtspServer.rtspCameras, k)
}

type wsInOut struct {
}

func NewWSStdInOut() wsInOut {
	return wsInOut{}
}

func (hk wsInOut) Write(p []byte) (n int, err error) {
	glogger.Logrus.Info(string(p))
	return len(p), nil
}
func (hk wsInOut) Read(p []byte) (n int, err error) {
	return len(p), nil
}

func (w *websocketPlayerManager) Write(p []byte) (n int, err error) {
	for _, c := range w.Clients {
		w.lock.Lock()
		err := c.WriteMessage(websocket.TextMessage, p)
		w.lock.Unlock()
		if err != nil {
			return 0, err
		}
	}
	return 0, nil
}

func NewPlayerManager() *websocketPlayerManager {
	return &websocketPlayerManager{
		WsServer: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		Clients: make(map[string]*websocket.Conn),
		lock:    sync.Mutex{},
	}

}

/*
*
* 启动服务
*
 */
type wsToken struct {
	Token  string `json:"token"`
	LiveId string `json:"live_id"`
}

func wsServerEndpoint(c *gin.Context) {
	//upgrade get request to websocket protocol
	wsConn, err := __DefaultRtspServer.websocketPlayerManager.WsServer.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	LiveId := c.Query("liveId")
	Token := c.Query("token")

	if Token != "WebRtspPlayer" {
		wsConn.WriteMessage(1, []byte("Invalid client token"))
		wsConn.Close()
		return
	}
	glogger.GLogger.Debugf("Request live:%s, Token is :%s", LiveId, Token)
	// 最多允许连接10个客户端，实际情况下根本用不了那么多
	if len(__DefaultRtspServer.websocketPlayerManager.Clients) >= 2 {
		wsConn.WriteMessage(websocket.CloseMessage, []byte{})
		wsConn.Close()
		return
	}
	__DefaultRtspServer.websocketPlayerManager.Clients[LiveId] = wsConn
	wsConn.WriteMessage(websocket.TextMessage, []byte("Connected"))
	glogger.GLogger.Info("WebSocket Player connected:" + wsConn.RemoteAddr().String())
	wsConn.SetCloseHandler(func(code int, text string) error {
		glogger.GLogger.Info("wsConn CloseHandler:", wsConn.RemoteAddr().String())
		__DefaultRtspServer.websocketPlayerManager.lock.Lock()
		delete(__DefaultRtspServer.websocketPlayerManager.Clients, wsConn.RemoteAddr().String())
		__DefaultRtspServer.websocketPlayerManager.lock.Unlock()
		return nil
	})
	wsConn.SetPingHandler(func(appData string) error {
		return nil
	})
	wsConn.SetPongHandler(func(appData string) error {
		return nil
	})
	go func(wsConn *websocket.Conn) {
		defer func() {
			if wsConn != nil {
				glogger.GLogger.Info("wsConn Disconnect By accident:", wsConn.RemoteAddr().String())
				__DefaultRtspServer.websocketPlayerManager.lock.Lock()
				delete(__DefaultRtspServer.websocketPlayerManager.Clients, wsConn.RemoteAddr().String())
				__DefaultRtspServer.websocketPlayerManager.lock.Unlock()
			}
		}()
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
				break
			}
			err = wsConn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				break
			}
		}
	}(wsConn)
}
