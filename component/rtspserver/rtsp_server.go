package rtspserver

import (
	"bufio"
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"

	"os/exec"

	"github.com/gin-gonic/gin"
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
	PullAddr string    `json:"pullAddr"`
	PushAddr string    `json:"pushAddr"`
	Process  *exec.Cmd `json:"-"`
}
type rtspServer struct {
	webServer   *gin.Engine
	rtspCameras map[string]RtspCameraInfo
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
		webServer:   gin.New(),
		rtspCameras: map[string]RtspCameraInfo{},
	}
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
		go StartFfmpegProcess(form.PullAddr, url1)
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
		StopFfmpegProcess((form.PullAddr))
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

func StartFfmpegProcess(rtspUrl, pushAddr string) error {
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
	glogger.GLogger.Info("Start FFMPEG process with:", cmd.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	info := RtspCameraInfo{
		PullAddr: rtspUrl,
		PushAddr: pushAddr,
		Process:  cmd,
	}
	// Run and Wait
	AddVideoStreamEndpoint(rtspUrl, info)
	defer DeleteVideoStreamEndpoint(rtspUrl)
	glogger.GLogger.Info("start Video Stream Endpoint:", rtspUrl)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(err.Error() + ":" + string(out))
	}
	glogger.GLogger.Info("stop Video Stream Endpoint:", rtspUrl)
	return nil
}

/*
*
* 停止进程
*
 */
func StopFfmpegProcess(rtspUrl string) error {
	if p := GetVideoStreamEndpoint(rtspUrl); p.Process != nil {
		return p.Process.Process.Kill()
	}
	return fmt.Errorf("no such process:" + rtspUrl)
}
