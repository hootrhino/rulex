package rtspserver

import (
	"bufio"
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"

	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/hootrhino/rulex/core"
	"github.com/hootrhino/rulex/glogger"
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
	Type          string `json:"type,omitempty"`     // 1-RTSP,2-Local
	LocalId       string `json:"local_id,omitempty"` // 本地ID
	PullAddr      string `json:"pullAddr,omitempty"`
	PushAddr      string `json:"pushAddr,omitempty"`
	ffmpegProcess *exec.Cmd
}

/*
*
* 这是用来给外部输出日志的websocket服务器，其功能非常简单，就是单纯的对外输出实时日志，方便调试使用。
* 注意：该功能需要配合HttpApiServer使用, 客户端连上以后必须在5s内发送一个 ‘WsTerminal’ 的固定字符
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

func calculateMD5(inputString string) string {
	hasher := md5.New()
	io.WriteString(hasher, inputString)
	hashBytes := hasher.Sum(nil)
	md5String := fmt.Sprintf("%x", hashBytes)
	return md5String
}
func isValidRTSPAddress(address string) bool {
	rtspPattern := `rtsp://[a-zA-Z0-9.-]+(:[0-9]+)?(/[a-zA-Z0-9/._-]*)?`
	matched, err := regexp.MatchString(rtspPattern, address)
	if err != nil {
		return false
	}
	return matched
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
	group.POST("/registerLive", func(ctx *gin.Context) {
		type Form struct {
			PullAddr string `json:"pull_addr,omitempty"`
		}
		form := Form{}
		if err := ctx.ShouldBindJSON(&form); err != nil {
			ctx.JSON(400, map[string]interface{}{
				"code": 4001,
				"msg":  err,
			})
			return
		}
		if !isValidRTSPAddress(form.PullAddr) {
			glogger.GLogger.Info("InValid RTSP Address:", form.PullAddr)
			return
		}
		url1 := "http://127.0.0.1:9400/stream/ffmpegPush?liveId=" + calculateMD5(form.PullAddr)
		url2 := "ws://127.0.0.1:9400/stream/live?liveId=" + calculateMD5(form.PullAddr)
		go StartFFMPEGProcess(form.PullAddr, url1)
		ctx.JSON(200, map[string]interface{}{
			"code": 200,
			"msg":  "Success",
			"data": url2,
		})
	})
	group.POST("/stopLive", func(ctx *gin.Context) {
		type Form struct {
			PullAddr string `json:"pull_addr,omitempty"`
		}
		form := Form{}
		if err := ctx.ShouldBindJSON(&form); err != nil {
			ctx.JSON(400, map[string]interface{}{
				"code": 4001,
				"msg":  err,
			})
			return
		}
		StopFFMPEGProcess((form.PullAddr))
		ctx.JSON(200, map[string]interface{}{
			"code": 200,
			"msg":  "Success",
		})
	})
	// 注意：这个接口是给FFMPEG请求的
	group.POST("/ffmpegPush", func(ctx *gin.Context) {
		LiveId := ctx.Query("liveId")
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
	fmt.Println(liveId, data)
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

/*
*
把外部RTSP流拉下来推给Go实现的RTSPServer
ffmpeg -rtsp_transport tcp -re

	-i 'RTSP-URL' -q 0 -f mpegts -c:v mpeg1video -an -s 1920x1080 http://127.0.0.1:%s/stream/%s

*
*/
/*
*
* Fork 一个进程来执行
*
 */

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

func StartFFMPEGProcess(rtspUrl, pushAddr string) error {
	params := []string{
		"-rtsp_transport",
		"tcp",
		"-re",
		"-i",
		// rtsp://192.168.199.243:554/av0_0
		rtspUrl,
		"-q",
		"5",
		"-f",
		"mpegts",
		"-fflags",
		"nobuffer",
		"-c:v",
		"mpeg1video",
		"-an",
		"-s",
		"1920x1080",
		// http://127.0.0.1:9400/stream/ffmpegPush?liveId=147a6d7ae5a785f6e3ea90f25d36c63e
		pushAddr,
	}

	cmd := exec.Command("ffmpeg", params...)
	glogger.GLogger.Info("Start FFMPEG ffmpegProcess with:", cmd.String())

	glogger.GLogger.Info("start Video Stream Endpoint:", rtspUrl)
	// 启动 FFmpeg 推流
	if err := cmd.Start(); err != nil {
		fmt.Printf("无法启动 FFmpeg: %v\n", err)
		return err
	}
	if core.GlobalConfig.AppDebugMode {
		inOut := wsInOut{}
		cmd.Stdin = nil
		cmd.Stdout = inOut
		cmd.Stderr = inOut
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	info := RtspCameraInfo{
		PullAddr:      rtspUrl,
		PushAddr:      pushAddr,
		ffmpegProcess: cmd,
	}

	// Run and Wait
	AddVideoStreamEndpoint(rtspUrl, info)
	defer DeleteVideoStreamEndpoint(rtspUrl)
	// 等待 FFmpeg 进程完成
	if err := cmd.Wait(); err != nil {
		return err
	}
	glogger.GLogger.Info("stop Video Stream Endpoint:", rtspUrl)
	return nil
}

/*
*
* 停止进程
*
 */
func StopFFMPEGProcess(rtspUrl string) error {
	if p := GetVideoStreamEndpoint(rtspUrl); p.ffmpegProcess != nil {
		return p.ffmpegProcess.Process.Kill()
	}
	return fmt.Errorf("no such ffmpegProcess:" + rtspUrl)
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

func wsServerEndpoint(c *gin.Context) {
	//upgrade get request to websocket protocol
	wsConn, err := __DefaultRtspServer.websocketPlayerManager.WsServer.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	// 首先读第一个包是不是: WsTerminal
	wsConn.SetReadDeadline(time.Now().Add(5 * time.Second))
	_, b, err := wsConn.ReadMessage()
	if err != nil {
		glogger.GLogger.Error(err)
		return
	}
	wsConn.SetReadDeadline(time.Time{})
	token := string(b)
	if token != "WebRtspPlayer" {
		wsConn.WriteMessage(1, []byte("Invalid client token"))
		wsConn.Close()
		return
	}
	// 最多允许连接10个客户端，实际情况下根本用不了那么多
	if len(__DefaultRtspServer.websocketPlayerManager.Clients) >= 10 {
		wsConn.WriteMessage(websocket.TextMessage, []byte("Reached max connections"))
		wsConn.Close()
		return
	}
	__DefaultRtspServer.websocketPlayerManager.Clients[wsConn.RemoteAddr().String()] = wsConn
	wsConn.WriteMessage(websocket.TextMessage, []byte("Connected"))
	glogger.GLogger.Info("WebSocket Terminal connected:" + wsConn.RemoteAddr().String())
	wsConn.SetCloseHandler(func(code int, text string) error {
		glogger.GLogger.Info("wsConn CloseHandler:", wsConn.RemoteAddr().String())
		return nil
	})
	go func(ctx context.Context, wsConn *websocket.Conn) {
		defer func() {
			wsConn.Close()
			__DefaultRtspServer.websocketPlayerManager.lock.Lock()
			delete(__DefaultRtspServer.websocketPlayerManager.Clients, wsConn.RemoteAddr().String())
			__DefaultRtspServer.websocketPlayerManager.lock.Unlock()
		}()
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
			_, _, err := wsConn.ReadMessage()
			wsConn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				glogger.GLogger.Info("WsConn error:", wsConn.RemoteAddr().String(), ", Error:", err)
				return
			}

		}

	}(context.Background(), wsConn)
}
