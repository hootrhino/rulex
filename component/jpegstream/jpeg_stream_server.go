package jpegstream

import (
	"bufio"
	"context"
	"fmt"

	"net"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rulex/component"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
)

var __DefaultJpegStreamServer *JpegStreamServer

type JpegStreamServer struct {
	webServer   *gin.Engine            // HTTP Server
	wsPort      int                    // 端口
	JpegStreams map[string]*JpegStream // 当前的JpegStream列表, Key:liveId
	locker      sync.Mutex
}

/*
*
* 初始化
*
 */
func InitJpegStreamServer(rulex typex.RuleX) {
	gin.SetMode(gin.ReleaseMode)
	__DefaultJpegStreamServer = &JpegStreamServer{
		webServer:   gin.New(),
		JpegStreams: map[string]*JpegStream{},
		locker:      sync.Mutex{},
		wsPort:      9401,
	}
	__DefaultJpegStreamServer.Init(map[string]any{})
	__DefaultJpegStreamServer.Start(rulex)
}

/*
*
* JpegStream Server: 原理很简单,就是把流推过来以后用一个HTTP流接收
*   然后返回图片的二进制流，这样前端的<img>标签就能展示出来了
*
 */

func (s *JpegStreamServer) Init(cfg map[string]any) error {

	// 注册Websocket server
	__DefaultJpegStreamServer.webServer.Use(utils.AllowCros)
	group := s.webServer.Group("/jpeg_stream")
	// 注意：这个接口是给FFMPEG请求的
	//    ffmpeg -stream_loop -1 -re -i 1.mp4 -c:v mjpeg -f mjpeg
	// -headers "Content-Type:multipart/x-mixed-replace" "http://127.0.0.1:9401/jpeg_stream/push?liveId=123"
	group.POST("/push", func(ctx *gin.Context) {
		defer ctx.Writer.Flush()
		// 从请求头中获取Content-Type，并解析出boundary
		LiveId := ctx.Query("liveId")
		if LiveId == "" {
			msg := "Missing required 'liveId', Example: http://host:9401/jpeg_stream/push?liveId=${your-liveId}"
			ctx.Writer.Write([]byte(msg))
			glogger.GLogger.Error(msg)
			return
		}
		if s.Exists(LiveId) {
			msg := "Jpeg Stream Already Exists:" + LiveId
			ctx.Writer.Write([]byte(msg))
			glogger.GLogger.Error(msg)
			return
		}
		glogger.GLogger.Info("Receive stream push From:", LiveId, ", jpeg stream play url is: ",
			fmt.Sprintf(`http://127.0.0.1:9401/jpeg_stream/pull?liveId=%s`, LiveId))
		// http://127.0.0.1:9400 :后期通过参数传进
		// 启动一个FFMPEG开始从摄像头拉流
		s.RegisterJpegStreamSource(LiveId)
		defer s.DeleteJpegStreamSource(LiveId)
		//
		bodyReader := bufio.NewReader(ctx.Request.Body)
		FrameStart1 := false                  // 标志位：FF
		FrameStart2 := false                  // 标志位：D8
		FrameEnd1 := false                    // 标志位：FF
		FrameEnd2 := false                    // 标志位：D9
		Offset := 0                           // 记录读取的字节
		var FrameBuffer = [1250000 / 2]byte{} // 默认5MB数据

		timeoutSignal := make(chan bool)
		// defer close(timeoutSignal)
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		go func() {
			<-ticker.C
			timeoutSignal <- true
		}()
		for {
			select {
			case <-ctx.Done():
				goto END
			case <-timeoutSignal:
				glogger.GLogger.Error("Read Jpeg Frame magic Number Timeout")
				goto END
			default:
				{
				}
			}
			// 逐行读取数据，直到找到 JPEG 图像文件的结束标识 FF D9
			AByte, err := bodyReader.ReadByte()
			if err != nil {
				glogger.GLogger.Error(err)
				break
			}
			if Offset == 0 && AByte == '\xFF' {
				FrameStart1 = true
				goto PARSE
			}
			if Offset == 1 && AByte == '\xD8' {
				FrameStart2 = true
				goto PARSE
			}
			if Offset > 2 && AByte == '\xFF' {
				FrameEnd1 = true
				goto PARSE
			}
			if Offset > 2 && AByte == '\xD9' && FrameBuffer[Offset-1] == '\xFF' {
				FrameEnd2 = true
				goto PARSE
			}
		PARSE:
			FrameBuffer[Offset] = AByte
			// 当取完整一帧以后就开始输出
			if (FrameStart1 && FrameStart2) && (FrameEnd1 && FrameEnd2) {
				JpegStream, err1 := s.GetJpegStreamSource(LiveId)
				if err1 != nil {
					ctx.Writer.Write([]byte(err1.Error()))
					glogger.GLogger.Error(err1.Error())
					return
				}
				if !JpegStream.Pulled && !JpegStream.GetFirstFrame {
					_, Resolution, err2 := utils.CvMatToImageBytes(FrameBuffer[:Offset])
					if err2 != nil {
						glogger.GLogger.Error(err1)
						return
					}
					JpegStream.Type = "JPEG_STREAM"    // 流
					JpegStream.LiveId = LiveId         // 设置ID
					JpegStream.Resolution = Resolution // 设置分辨率
					JpegStream.GetFirstFrame = true    // 获取到第一帧
				}
				// 当没有被拉流的时候不需要推流, 减少OpenCV消耗
				if JpegStream.Pulled {
					ImageBytes, _, err2 := utils.CvMatToImageBytes(FrameBuffer[:Offset])
					if err2 != nil {
						glogger.GLogger.Error(err1)
						break
					}
					defer ImageBytes.Close()
					JpegStream.UpdateJPEG(ImageBytes.GetBytes()) // 刷新帧
					for i := range FrameBuffer[:Offset] {
						FrameBuffer[i] = 0
					}
				}
				FrameStart1 = false // 标志位：FF
				FrameStart2 = false // 标志位：D8
				FrameEnd1 = false   // 标志位：FF
				FrameEnd2 = false   // 标志位：D9
				Offset = 0          // 初始化游标
				go func() {
					ticker.Reset(5 * time.Second) // 重置计时器
					<-ticker.C
					timeoutSignal <- true
				}()
			} else {
				if Offset > 625000 {
					glogger.GLogger.Errorf("Jpeg Frame Size too large:%d", Offset)
					break
				}
				Offset++
			}
		}
	END:
		glogger.GLogger.Info("Stream push stop:", LiveId)
	})
	group.GET("/pull", func(ctx *gin.Context) {
		defer ctx.Writer.Flush()
		liveId, _ := ctx.GetQuery("liveId")
		JpegStream, err := s.GetJpegStreamSource(liveId)
		if err != nil {
			glogger.GLogger.Error(err)
			return
		}
		JpegStream.Pulled = true
		defer func() {
			JpegStream.Pulled = false
		}()
		glogger.GLogger.Infof("Client [%s] Start pull jpeg stream [%s]:",
			liveId, ctx.Request.RemoteAddr)
		ctx.Header("Content-Type", "multipart/x-mixed-replace; boundary=MJPEG_BOUNDARY")
		ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		ctx.Header("Pragma", "no-cache")
		ctx.Header("Expires", "0")
		ctx.Header("Age", "0")
		ctx.Writer.Write([]byte("HTTP/1.1 200 OK\r\n"))
		for {
			_, err := ctx.Writer.Write(JpegStream.GetWebJpegFrame())
			if err != nil {
				glogger.GLogger.Error("Jpeg Stream Server Write error", err)
				break
			}
			time.Sleep(50 * time.Millisecond) // 1000/50=20帧
		}
		glogger.GLogger.Infof("Client [%s] pull jpeg stream[%s] finished:",
			ctx.Request.RemoteAddr, liveId)

	})
	return nil
}
func (s *JpegStreamServer) Start(r typex.RuleX) error {
	go func(ctx context.Context) {
		listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", __DefaultJpegStreamServer.wsPort))
		if err != nil {
			glogger.GLogger.Fatalf("JpegStream stream server listen error: %s\n", err)
			return
		}
		defer listener.Close()
		if err := __DefaultJpegStreamServer.webServer.RunListener(listener); err != nil {
			glogger.GLogger.Fatalf("JpegStream stream server listen error: %s\n", err)
			return
		}
	}(context.Background())
	glogger.GLogger.Info("JpegStream stream server start success, listening at:",
		fmt.Sprintf("0.0.0.0:%d", __DefaultJpegStreamServer.wsPort))
	return nil
}
func (s *JpegStreamServer) Stop() error {
	return nil
}
func (s *JpegStreamServer) PluginMetaInfo() component.XComponentMetaInfo {
	return component.XComponentMetaInfo{}
}

/*
*
* Manage API
*
 */

func (s *JpegStreamServer) RegisterJpegStreamSource(liveId string) error {
	s.locker.Lock()
	defer s.locker.Unlock()
	_, ok := s.JpegStreams[liveId]
	if !ok {
		s.JpegStreams[liveId] = &JpegStream{
			GetFirstFrame: false,
			LiveId:        liveId,
			Pulled:        false,
			Resolution:    utils.Resolution{Width: 0, Height: 0},
			frame:         []byte{},
		}
		return nil
	}
	return fmt.Errorf("stream already exists")
}

func (s *JpegStreamServer) GetJpegStreamSource(liveId string) (*JpegStream, error) {
	s.locker.Lock()
	defer s.locker.Unlock()
	JpegStream, ok := s.JpegStreams[liveId]
	if ok {
		return JpegStream, nil
	} else {
		return JpegStream, fmt.Errorf("stream not exists")
	}
}

func (s *JpegStreamServer) Exists(liveId string) bool {
	s.locker.Lock()
	defer s.locker.Unlock()
	_, ok := s.JpegStreams[liveId]
	return ok
}
func (s *JpegStreamServer) DeleteJpegStreamSource(liveId string) {
	s.locker.Lock()
	defer s.locker.Unlock()
	delete(s.JpegStreams, liveId)
}

func (s *JpegStreamServer) JpegStreamSourceList() []JpegStream {
	List := []JpegStream{}
	for _, v := range s.JpegStreams {
		List = append(List, *v)
	}
	return List
}
func (s *JpegStreamServer) JpegStreamFlush() {
	for k := range s.JpegStreams {
		delete(s.JpegStreams, k)
	}
}

type JpegStream struct {
	frame         []byte
	frameSize     int
	headerSize    int
	GetFirstFrame bool
	Type          string
	LiveId        string
	Pulled        bool
	Resolution    utils.Resolution
}

func (S JpegStream) String() string {
	return fmt.Sprintf(`{"liveId":%s,"pulled":%v,"resolution":%s}`,
		S.LiveId, S.Pulled, S.Resolution.String())
}
func (s *JpegStream) GetWebJpegFrame() []byte {
	b := s.frame[:s.frameSize]
	return b
}
func (s *JpegStream) GetRawFrame() []byte {
	return s.frame[s.headerSize:s.frameSize]
}
func (s *JpegStream) UpdateJPEG(jpeg []byte) {
	__headerReal := "--MJPEG_BOUNDARY\r\nContent-Type: image/jpeg\r\nContent-Length: %d\r\nX-Timestamp: 0.000000\r\n\r\n"
	size := len(jpeg)
	frameSize := size
	headerReal := fmt.Sprintf(__headerReal, size)
	if len(s.frame) < size+len(headerReal) {
		s.frame = make([]byte, (size+len(headerReal))*2)
	}
	copy(s.frame, headerReal)
	copy(s.frame[len(headerReal):], jpeg[:size])
	s.headerSize = len(headerReal) + 10 // \r\n出现了5次 计算为10个字符
	s.frameSize = frameSize

}
