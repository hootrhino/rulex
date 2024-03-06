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
var __DefaultRtspServer *RtspServer

/*
*
* RTSP 设备(rtsp://192.168.199.243:554/av0_0)
*
 */
type FlvStream struct {
	Type          string `json:"type"`     // 1-RTSP,2-Local
	LocalId       string `json:"local_id"` // 本地ID
	PullAddr      string `json:"pullAddr"`
	PushAddr      string `json:"pushAddr"`
	GetFirstFrame bool
	LiveId        string
	Pulled        bool
	Resolution    utils.Resolution
}

type RtspServer struct {
	webServer   *gin.Engine
	wsPort      int
	locker      sync.Mutex
	RtspStreams map[string]*FlvStream
	WsServer    websocket.Upgrader
	Clients     map[string]streamPlayer
}

// NewRouter Gin 路由配置
func InitRtspServer(rulex typex.RuleX) *RtspServer {
	gin.SetMode(gin.ReleaseMode)
	__DefaultRtspServer = &RtspServer{
		webServer:   gin.New(),
		RtspStreams: map[string]*FlvStream{},
		WsServer: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		Clients: make(map[string]streamPlayer),
		locker:  sync.Mutex{},
		wsPort:  9400,
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
			if len(__DefaultRtspServer.Clients) > 0 {
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
	__DefaultRtspServer.locker.Lock()
	defer __DefaultRtspServer.locker.Unlock()
	// 检查到底有没有拉流的,如果有就推给他
	for _, C := range __DefaultRtspServer.Clients {
		if C.LiveId == liveId {
			if C.wsConn != nil {
				C.wsConn.WriteMessage(websocket.BinaryMessage, data)
			}
		}
	}
}

/*
*
* Manage API
*
 */

func (s *RtspServer) RegisterFlvStreamSource(liveId string) error {
	s.locker.Lock()
	defer s.locker.Unlock()
	_, ok := s.RtspStreams[liveId]
	if !ok {
		s.RtspStreams[liveId] = &FlvStream{
			GetFirstFrame: false,
			LiveId:        liveId,
			Pulled:        false,
			Resolution:    utils.Resolution{Width: 0, Height: 0},
		}
		return nil
	}
	return fmt.Errorf("stream already exists")
}

func (s *RtspServer) GetFlvStreamSource(liveId string) (*FlvStream, error) {
	s.locker.Lock()
	defer s.locker.Unlock()
	FlvStream, ok := s.RtspStreams[liveId]
	if ok {
		return FlvStream, nil
	} else {
		return FlvStream, fmt.Errorf("stream not exists")
	}
}

func (s *RtspServer) Exists(liveId string) bool {
	s.locker.Lock()
	defer s.locker.Unlock()
	_, ok := s.RtspStreams[liveId]
	return ok
}
func (s *RtspServer) DeleteFlvStreamSource(liveId string) {
	s.locker.Lock()
	defer s.locker.Unlock()
	delete(s.RtspStreams, liveId)
}

func (s *RtspServer) FlvStreamSourceList() []FlvStream {
	List := []FlvStream{}
	for _, v := range s.RtspStreams {
		List = append(List, *v)
	}
	return List
}
func (s *RtspServer) FlvStreamFlush() {
	for k := range s.RtspStreams {
		delete(s.RtspStreams, k)
	}
}
