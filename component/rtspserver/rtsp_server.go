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
	"github.com/hootrhino/rulex/utils"
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

type websocketPlayerManager struct {
	WsServer websocket.Upgrader
	Clients  map[string]StreamPlayer
	lock     sync.Mutex
}

type rtspServer struct {
	webServer              *gin.Engine
	rtspCameras            map[string]RtspCameraInfo
	websocketPlayerManager *websocketPlayerManager
	wsPort                 int
}

// NewRouter Gin 路由配置
func InitRtspServer(rulex typex.RuleX) *rtspServer {
	gin.SetMode(gin.ReleaseMode)
	__DefaultRtspServer = &rtspServer{
		webServer:              gin.New(),
		rtspCameras:            map[string]RtspCameraInfo{},
		websocketPlayerManager: NewPlayerManager(),
		wsPort:                 9400,
	}
	// 注册Websocket server
	__DefaultRtspServer.webServer.Use(utils.AllowCros)
	__DefaultRtspServer.webServer.GET("/ws", wsServerEndpoint)
	// http://127.0.0.1:3000/h264_stream/live/001
	group := __DefaultRtspServer.webServer.Group("/h264_stream")
	// 注意：这个接口是给FFMPEG请求的
	//    ffmpeg -hide_banner -r 24 -rtsp_transport tcp -re -i rtsp://192.168.1.210:554/av0_0 -q 5 -f mpegts
	// -fflags nobuffer -c:v libx264 -an -s 1920x1080 http://?:9400/h264_stream/push?liveId=xx
	group.POST("/push", func(ctx *gin.Context) {
		LiveId := ctx.Query("liveId")
		if LiveId == "" {
			ctx.Writer.Write([]byte("missing required 'liveId'"))
			glogger.GLogger.Error("missing required 'liveId'")
			ctx.Writer.Flush()
			return
		}
		glogger.GLogger.Info("Receive stream push From:", LiveId,
			", stream play url is:", fmt.Sprintf("ws://127.0.0.1:9400/ws?token=WebRtspPlayer&liveId=%s", LiveId))
		// http://127.0.0.1:9400 :后期通过参数传进
		// 启动一个FFMPEG开始从摄像头拉流
		bodyReader := bufio.NewReader(ctx.Request.Body)
		// 每个ffmpeg都给起一个进程，监控是否有websocket来拉流
		for {
			// data 就是 RTSP 帧
			// 只需将其转发给websocket即可
			data, err := bodyReader.ReadBytes('\n')
			if err != nil {
				glogger.GLogger.Error("ReadBytes from ffmpeg error:", err)
				break
			}
			if len(__DefaultRtspServer.websocketPlayerManager.Clients) > 0 {
				pushToWebsocket(LiveId, data)
			}
		}
		ctx.Writer.Flush()
		glogger.GLogger.Info("Stream push stop:", LiveId)
	})
	go func(ctx context.Context) {
		listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", __DefaultRtspServer.wsPort))
		if err != nil {
			glogger.GLogger.Fatalf("Rtsp stream server listen error: %s\n", err)
		}
		defer listener.Close()
		if err := __DefaultRtspServer.webServer.RunListener(listener); err != nil {
			glogger.GLogger.Fatalf("Rtsp stream server listen error: %s\n", err)
		}
	}(context.Background())
	glogger.GLogger.Info("Rtsp stream server start success, listening at:",
		fmt.Sprintf("0.0.0.0:%d", __DefaultRtspServer.wsPort))
	return __DefaultRtspServer
}

/*
*
* 推流
*
 */
func pushToWebsocket(liveId string, data []byte) {
	__DefaultRtspServer.websocketPlayerManager.lock.Lock()
	defer __DefaultRtspServer.websocketPlayerManager.lock.Unlock()
	// 检查到底有没有拉流的,如果有就推给他
	for _, C := range __DefaultRtspServer.websocketPlayerManager.Clients {
		if C.RequestLiveId == liveId {
			if C.wsConn != nil {
				C.wsConn.WriteMessage(websocket.BinaryMessage, data)
			}
		}
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
	__DefaultRtspServer.websocketPlayerManager.lock.Lock()
	defer __DefaultRtspServer.websocketPlayerManager.lock.Unlock()
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
		err := c.wsConn.WriteMessage(websocket.TextMessage, p)
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
		Clients: make(map[string]StreamPlayer),
		lock:    sync.Mutex{},
	}

}

/*
*
* websocket播放器
*
 */
type StreamPlayer struct {
	ClientId      string
	RequestLiveId string
	Token         string
	wsConn        *websocket.Conn
	Resolution    utils.Resolution
}

/*
*
* 启动服务
*
 */
func wsServerEndpoint(c *gin.Context) {
	wsConn, err := __DefaultRtspServer.websocketPlayerManager.WsServer.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	LiveId := c.Query("liveId")
	ClientId := c.Query("clientId")
	Token := c.Query("token")
	glogger.GLogger.Debugf("Request live:%s, Token is :%s", LiveId, Token)
	// 最多允许连接10个客户端，实际情况下根本用不了那么多
	if len(__DefaultRtspServer.websocketPlayerManager.Clients) >= 10 {
		wsConn.WriteMessage(websocket.CloseMessage, []byte{})
		wsConn.Close()
		return
	}

	StreamPlayer := StreamPlayer{
		RequestLiveId: LiveId,
		ClientId:      ClientId,
		Token:         Token,
		wsConn:        wsConn,
	}
	if Token != "WebRtspPlayer" {
		wsConn.WriteMessage(websocket.CloseMessage, []byte("Invalid client token"))
		wsConn.Close()
		return
	}
	if C, ok := __DefaultRtspServer.websocketPlayerManager.Clients[ClientId]; ok {
		wsConn.WriteMessage(websocket.CloseMessage, []byte("already exists a client:"+C.ClientId))
		wsConn.Close()
		return
	}
	// 每个LiveId只能有1路播放
	__DefaultRtspServer.websocketPlayerManager.Clients[ClientId] = StreamPlayer
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
				glogger.GLogger.Warn("wsConn Close Handler:", wsConn.RemoteAddr().String())
				__DefaultRtspServer.websocketPlayerManager.lock.Lock()
				delete(__DefaultRtspServer.websocketPlayerManager.Clients, wsConn.RemoteAddr().String())
				__DefaultRtspServer.websocketPlayerManager.lock.Unlock()
				break
			}
			err = wsConn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				break
			}
		}
	}(wsConn)
}
