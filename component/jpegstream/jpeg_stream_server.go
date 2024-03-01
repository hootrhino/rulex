package jpegstream

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
)

var __DefaultJpegStreamServer *JpegStreamServer

type JpegStreamServer struct {
	webServer   *gin.Engine // HTTP Server
	wsPort      int         // 端口
	locker      sync.Mutex
	JpegStreams map[string]JpegStream // 当前的JpegStream列表, Key:JpegStreamId
}

func RegisterJpegStreamSource(JpegStreamId string) error {
	__DefaultJpegStreamServer.locker.Lock()
	defer __DefaultJpegStreamServer.locker.Unlock()
	_, ok := __DefaultJpegStreamServer.JpegStreams[JpegStreamId]
	if !ok {
		__DefaultJpegStreamServer.JpegStreams[JpegStreamId] = JpegStream{}
		return nil
	}
	return fmt.Errorf("stream already exists")
}
func GetJpegStreamSource(JpegStreamId string) (JpegStream, error) {
	__DefaultJpegStreamServer.locker.Lock()
	defer __DefaultJpegStreamServer.locker.Unlock()
	JpegStream, ok := __DefaultJpegStreamServer.JpegStreams[JpegStreamId]
	if !ok {
		return JpegStream, nil
	} else {
		return JpegStream, fmt.Errorf("stream not exists")
	}
}
func DeleteJpegStreamSource(JpegStreamId string) {
	__DefaultJpegStreamServer.locker.Lock()
	defer __DefaultJpegStreamServer.locker.Unlock()
	delete(__DefaultJpegStreamServer.JpegStreams, JpegStreamId)
}
func UpdateJpegStreamSource(JpegStreamId string, frame []byte) error {
	__DefaultJpegStreamServer.locker.Lock()
	defer __DefaultJpegStreamServer.locker.Unlock()
	_, ok := __DefaultJpegStreamServer.JpegStreams[JpegStreamId]
	if ok {
		__DefaultJpegStreamServer.JpegStreams[JpegStreamId].UpdateJPEG(frame)
		return nil
	}
	return fmt.Errorf("stream not exists")

}
func JpegStreamSourceList() []JpegStream {
	List := []JpegStream{}
	for _, v := range __DefaultJpegStreamServer.JpegStreams {
		List = append(List, v)
	}
	return List
}
func JpegStreamFlush() {
	for k := range __DefaultJpegStreamServer.JpegStreams {
		delete(__DefaultJpegStreamServer.JpegStreams, k)
	}
}

/*
*
* JpegStream Server: 原理很简单,就是把流推过来以后用一个HTTP流接收
*   然后返回图片的二进制流，这样前端的<img>标签就能展示出来了
*
 */
func InitJpegStreamServer(rulex typex.RuleX) *JpegStreamServer {
	gin.SetMode(gin.ReleaseMode)
	__DefaultJpegStreamServer = &JpegStreamServer{
		webServer: gin.New(),
		locker:    sync.Mutex{},
		wsPort:    9401,
	}
	// 注册Websocket server
	__DefaultJpegStreamServer.webServer.Use(cros)
	group := __DefaultJpegStreamServer.webServer.Group("/stream")
	group.POST("/jpeg_stream", func(ctx *gin.Context) {
		JpegStreamId, _ := ctx.GetQuery("JpegStreamId")
		__DefaultJpegStreamServer.locker.Lock()
		JpegStream, ok := __DefaultJpegStreamServer.JpegStreams[JpegStreamId]
		__DefaultJpegStreamServer.locker.Unlock()
		if !ok {
			msg := "Stream not exists:" + JpegStreamId
			glogger.GLogger.Error(msg)
			ctx.Writer.Write([]byte(msg))
			ctx.Writer.Flush()
			return
		}
		glogger.GLogger.Infof("Client [%s] Start pull jpeg stream [%s]:", JpegStreamId, ctx.Request.RemoteAddr)
		ctx.Request.Header.Add("Content-Type", "multipart/x-mixed-replace;boundary=MJPEGBOUNDARY")
		for {
			_, err := ctx.Writer.Write(JpegStream.frame)
			if err != nil {
				glogger.GLogger.Error("Jpeg Stream Server Write error", err)
				break
			}
		}
		glogger.GLogger.Infof("Client [%s] pull jpeg stream[%s] finished:", ctx.Request.RemoteAddr, JpegStreamId)

	})
	go func(ctx context.Context) {
		listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", __DefaultJpegStreamServer.wsPort))
		if err != nil {
			glogger.GLogger.Fatalf("JpegStream stream server listen error: %s\n", err)
		}
		defer listener.Close()
		if err := __DefaultJpegStreamServer.webServer.RunListener(listener); err != nil {
			glogger.GLogger.Fatalf("JpegStream stream server listen error: %s\n", err)
		}
	}(context.Background())
	glogger.GLogger.Info("JpegStream stream server start success, listening at:",
		fmt.Sprintf("0.0.0.0:%d", __DefaultJpegStreamServer.wsPort))
	return __DefaultJpegStreamServer
}

func cros(c *gin.Context) {
	c.Header("Cache-Control", "private, max-age=10")
	method := c.Request.Method
	origin := c.Request.Header.Get("Origin")

	c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
	c.Header("Access-Control-Max-Age", "172800")
	c.Header("Access-Control-Allow-Credentials", "true")

	if method == http.MethodOptions {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	c.Request.Header.Del("Origin")
	c.Next()
}

/*
*
* Jpeg Stream 帧
*
 */
type JpegStream struct {
	frame         []byte
	FrameInterval time.Duration
}

const headerf = "\r\n" +
	"--MJPEGBOUNDARY\r\n" +
	"Content-Type: image/jpeg\r\n" +
	"Content-Length: %d\r\n" +
	"X-Timestamp: 0.000000\r\n" +
	"\r\n"

func (s JpegStream) UpdateJPEG(jpeg []byte) {
	header := fmt.Sprintf(headerf, len(jpeg))
	if len(s.frame) < len(jpeg)+len(header) {
		s.frame = make([]byte, (len(jpeg)+len(header))*2)
	}
	copy(s.frame, header)
	copy(s.frame[len(header):], jpeg)
}
