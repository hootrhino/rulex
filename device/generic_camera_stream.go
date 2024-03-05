package device

import (
	"crypto/md5"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"syscall"

	"github.com/hootrhino/rulex/component/iotschema"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
)

const (
	__INPUT_REMOTE_STREAM_RTSP        = "REMOTE_STREAM_RTSP"       // 远程RTSP拉流
	__INPUT_LOCAL_CAMERA              = "LOCAL_CAMERA"             // 本地摄像头
	__OUTPUT_LOCAL_H264_STREAM_SERVER = "LOCAL_H264_STREAM_SERVER" // RULEX 自带的FLVServer
	__OUTPUT_LOCAL_JPEG_STREAM_SERVER = "LOCAL_JPEG_STREAM_SERVER" // RULEX 自带的FLVServer
	__OUTPUT_REMOTE_STREAM_SERVER     = "REMOTE_STREAM_SERVER"     // 远程地址
	__internal_ws_server_url          = "http://127.0.0.1:9400/h264_stream/push?liveId="
	__internal_jpeg_stream_server_url = "http://127.0.0.1:9401/jpeg_stream/push?liveId="
	// 输出模式
	__OUTPUT_MODE_H264_STREAM = "H264_STREAM"
	__OUTPUT_MODE_JPEG_STREAM = "JPEG_STREAM"
)

// RTSP URL格式= rtsp://<username>:<password>@<ip>:<port>，
// 默认为本地拉流，推向Rulex自带的Stream server
type _MainConfig struct {
	InputMode    string `json:"inputMode" validate:"required"`    // 输入模式: REMOTE_STREAM_RTSP | LOCAL_CAMERA
	InputAddr    string `json:"inputAddr" validate:"required"`    // 视频路径，REMOTE_STREAM_RTSP是远程地址，LOCAL_CAMERA是本地设备
	OutputMode   string `json:"outputMode" validate:"required"`   // 输出模式: LOCAL_H264_STREAM_SERVER | LOCAL_JPEG_STREAM_SERVER | REMOTE_STREAM_SERVER
	OutputEncode string `json:"outputEncode" validate:"required"` // 输出编码: H264_STREAM | JPEG_STREAM
	OutputAddr   string `json:"outputAddr"  validate:"required"`  // 输出地址, 格式为: "Ip:Port",例如127.0.0.1:7890
}

// 摄像头
type videoCamera struct {
	typex.XStatus
	status        typex.DeviceState
	mainConfig    _MainConfig
	ffmpegProcess *exec.Cmd
}

func NewVideoCamera(e typex.RuleX) typex.XDevice {
	videoCamera := new(videoCamera)
	videoCamera.RuleEngine = e
	videoCamera.status = typex.DEV_DOWN
	videoCamera.mainConfig = _MainConfig{
		InputMode:    __INPUT_LOCAL_CAMERA,              // 默认是本地摄像头
		InputAddr:    "0",                               // 默认值
		OutputMode:   __OUTPUT_LOCAL_JPEG_STREAM_SERVER, // 默认输出到JpegServer
		OutputEncode: __OUTPUT_MODE_JPEG_STREAM,         // 编码为JpegStream
		OutputAddr:   __internal_jpeg_stream_server_url, // 默认输出地址到本地Jpeg stream Server
	}
	return videoCamera
}

// 初始化 通常用来获取设备的配置
func (vc *videoCamera) Init(devId string, configMap map[string]interface{}) error {
	vc.PointId = devId
	if err := utils.BindSourceConfig(configMap, &vc.mainConfig); err != nil {
		return err
	}
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	if vc.mainConfig.InputMode == __INPUT_REMOTE_STREAM_RTSP {
		if !isValidRTSPAddress(vc.mainConfig.InputAddr) {
			return fmt.Errorf("invalid RtspUrl Format:%s", vc.mainConfig.InputAddr)
		}
	}
	return nil
}
func isValidRTSPAddress(address string) bool {
	rtspPattern := `rtsp://[a-zA-Z0-9.-]+(:[0-9]+)?(/[a-zA-Z0-9/._-]*)?`
	matched, err := regexp.MatchString(rtspPattern, address)
	if err != nil {
		return false
	}
	return matched
}

/*
*
## 3种情况
1. 向本地FLV服务器推H264流
2. 向本地Jpeg Stream推JpegStream流
3. 向远端推流，当前默认只能推RTMP，原计划配合SRS流媒体服务器使用，其他的后期再建设
*
*/
func (vc *videoCamera) Start(cctx typex.CCTX) error {
	// 本地JPEG只能推JPEG格式
	if vc.mainConfig.OutputMode == __OUTPUT_LOCAL_JPEG_STREAM_SERVER {
		if vc.mainConfig.OutputEncode == __OUTPUT_MODE_JPEG_STREAM {
			// Jpeg Stream 推流地址
			pushUrl := __internal_jpeg_stream_server_url + calculateMD5(vc.mainConfig.InputAddr)
			go vc.startFFMPEGProcess(vc.mainConfig.InputAddr, pushUrl)
		}
		vc.status = typex.DEV_UP
		return nil
	}
	// 本地FLV只能推FLV格式
	if vc.mainConfig.OutputMode == __OUTPUT_LOCAL_H264_STREAM_SERVER {
		if vc.mainConfig.OutputEncode == __OUTPUT_MODE_H264_STREAM {
			// Websocket 推流地址
			pushUrl := __internal_ws_server_url + calculateMD5(vc.mainConfig.InputAddr)
			go vc.startFFMPEGProcess(vc.mainConfig.InputAddr, pushUrl)
		}
		vc.status = typex.DEV_UP
		return nil
	}
	// 输出到远端只能H264
	if vc.mainConfig.OutputMode == __OUTPUT_REMOTE_STREAM_SERVER {
		if vc.mainConfig.OutputEncode == __OUTPUT_MODE_H264_STREAM {
			// 直接使用远程地址
			go vc.startFFMPEGProcess(vc.mainConfig.InputAddr, vc.mainConfig.OutputAddr)
		}
		vc.status = typex.DEV_UP
		return nil
	}
	glogger.GLogger.Errorf("Start failed, Output Mode Invalid:%s", vc.mainConfig.OutputMode)
	return nil
}

func (vc *videoCamera) OnRead(cmd []byte, data []byte) (int, error) {
	return 0, nil
}

func (vc *videoCamera) OnWrite(cmd []byte, data []byte) (int, error) {
	return 0, nil
}

/*
*
* 可以在这里进行PTZ控制
*
 */
func (vc *videoCamera) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	return nil, nil
}

// 设备当前状态
func (vc *videoCamera) Status() typex.DeviceState {
	return vc.status
}

func (vc *videoCamera) Stop() {
	vc.status = typex.DEV_DOWN
	if vc.CancelCTX != nil {
		vc.CancelCTX()
	}
	vc.stopFFMPEGProcess()
}

func (vc *videoCamera) Property() []iotschema.IoTSchema {
	return []iotschema.IoTSchema{}
}

func (vc *videoCamera) Details() *typex.Device {
	return vc.RuleEngine.GetDevice(vc.PointId)

}

func (vc *videoCamera) SetState(_ typex.DeviceState) {

}

func (vc *videoCamera) Driver() typex.XExternalDriver {
	return nil
}

func (vc *videoCamera) OnDCACall(_ string, _ string, _ interface{}) typex.DCAResult {
	return typex.DCAResult{}
}

/*
* 下面是v0.6.7规划的功能
* ## 当前支持的模式:
√ 本地摄像头      --> 本地JpegStream
√ 本地摄像头      --> 本地WebsocketFLV
√ RTSP网络摄像头  --> 本地JpegStream
√ RTSP网络摄像头  --> 本地WebsocketFLV
* ## 下个版本
... 本地摄像头      --> 远程RTMP
... RTSP网络摄像头  --> 远程RTMP
*
*/

func (vc *videoCamera) startFFMPEGProcess(inputUrl, pushAddr string) {
	defer func() {
		vc.status = typex.DEV_DOWN
	}()

	// 本地摄像头推向本地JpegStream服务器
	if vc.mainConfig.OutputMode == __OUTPUT_LOCAL_JPEG_STREAM_SERVER {
		glogger.GLogger.Info("Start FFMPEG ffmpegProcess with: LOCAL_JPEG_STREAM_SERVER Mode")
		var cmd *exec.Cmd
		// 本地摄像头
		if vc.mainConfig.InputMode == __INPUT_LOCAL_CAMERA {
			var deviceName string
			var paramsVideo []string
			if runtime.GOOS == "windows" {
				deviceName = fmt.Sprintf("video=\"%s\"", inputUrl)
			}
			if runtime.GOOS == "linux" {
				deviceName = fmt.Sprintf("video=%s", inputUrl)
			}
			if runtime.GOOS == "windows" {
				paramsVideo = []string{
					"-stream_loop", "-1",
					"-re",
					"-f", "dshow", // windows下特有的DirectX加速引擎
					"-i", deviceName,
					"-c:v", "mjpeg",
					"-f", "mjpeg",
					"-headers", `"Content-Type:multipart/x-mixed-replace; boundary=MJPEG_BOUNDARY"`,
					pushAddr,
				}
			} else {
				paramsVideo = []string{
					"-stream_loop", "-1",
					"-re",
					"-i", deviceName,
					"-c:v", "mjpeg",
					"-f", "mjpeg",
					"-headers", `"Content-Type:multipart/x-mixed-replace; boundary=MJPEG_BOUNDARY"`,
					pushAddr,
				}
			}
			cmd = concatFFmpegCommand(paramsVideo)

		}
		// RTSP摄像头
		if vc.mainConfig.InputMode == __INPUT_REMOTE_STREAM_RTSP {
			params := []string{
				"-hide_banner",
				"-r", "24",
				"-rtsp_transport", "tcp",
				"-re",
				"-i", inputUrl, //"rtsp://192.168.1.210:554/av0_0"
				"-c:v", "mjpeg",
				"-f", "mjpeg",
				"-preset", "veryfast",
				"-tune", "zerolatency",
				"-headers", `"Content-Type:multipart/x-mixed-replace; boundary=MJPEG_BOUNDARY"`,
				pushAddr, // "http://127.0.0.1:9401/jpeg_stream/push?liveId=123",
			}
			cmd = concatFFmpegCommand(params)
		}
		if cmd == nil {
			glogger.GLogger.Error(fmt.Errorf("InputMode Not Supported:" + vc.mainConfig.InputMode))
			return
		}
		glogger.GLogger.Debug("Start FFMPEG with:", cmd.String())
		// 启动 FFmpeg 推流
		vc.ffmpegProcess = cmd
		if output, err1 := cmd.CombinedOutput(); err1 != nil {
			glogger.GLogger.Error("Combined Output error: ", err1)
			fmt.Println(string(output))
			return
		}
	}
	// RTSP摄像头推向本地WebsocketFLV服务器
	// ffmpeg -i rtsp://IP/av0_0 -c:v h264 -c:a aac -f rtsp rtsp://IP/live/test001
	if vc.mainConfig.OutputMode == __OUTPUT_LOCAL_H264_STREAM_SERVER {
		glogger.GLogger.Info("Start FFMPEG ffmpegProcess with: LOCAL_H264_STREAM_SERVER Mode")
		var cmd *exec.Cmd
		if vc.mainConfig.InputMode == __INPUT_LOCAL_CAMERA {
			var paramsVideo []string
			var deviceName string
			if runtime.GOOS == "windows" {
				deviceName = fmt.Sprintf("video=\"%s\"", inputUrl)
			}
			if runtime.GOOS == "linux" {
				deviceName = fmt.Sprintf("video=%s", inputUrl)
			}
			if runtime.GOOS == "windows" {
				paramsVideo = []string{
					"-err_detect",
					"ignore_err",
					"-hide_banner",
					"-f", "dshow", // windows下特有的DirectX加速引擎
					"-i", deviceName,
					"-q", "5",
					"-fflags", "+genpts",
					"-preset", "veryfast",
					"-tune", "zerolatency",
					"-c:v", "libx264",
					"-f", "mpegts",
					pushAddr,
				}
			} else {
				paramsVideo = []string{
					"-err_detect",
					"ignore_err",
					"-hide_banner",
					"-i", deviceName,
					"-q", "5",
					"-fflags", "+genpts",
					"-preset", "veryfast",
					"-tune", "zerolatency",
					"-c:v", "libx264",
					"-f", "mpegts",
					pushAddr,
				}
			}
			cmd = concatFFmpegCommand(paramsVideo)
		}
		// RTSP转地址
		if vc.mainConfig.InputMode == __INPUT_REMOTE_STREAM_RTSP {
			params := []string{
				"-hide_banner",
				"-r", "24",
				"-rtsp_transport", "tcp",
				"-re",
				"-i", inputUrl, //"rtsp://192.168.1.210:554/av0_0"
				"-c:v", "mjpeg",
				"-f", "mjpeg",
				"-preset", "veryfast",
				"-tune", "zerolatency",
				"-headers", `"Content-Type:multipart/x-mixed-replace; boundary=MJPEG_BOUNDARY"`,
				pushAddr, // "http://127.0.0.1:9401/jpeg_stream/push?liveId=123",
			}
			cmd = concatFFmpegCommand(params)
		}
		if cmd == nil {
			glogger.GLogger.Error(fmt.Errorf("InputMode Not Supported:" + vc.mainConfig.InputMode))
			return
		}
		glogger.GLogger.Debug("Start FFMPEG with command line:", cmd.String())
		// 启动 FFmpeg 推流
		vc.ffmpegProcess = cmd
		if output, err1 := cmd.CombinedOutput(); err1 != nil {
			glogger.GLogger.Error("Combined Output error: ", err1, ", output: ", string(output))
		}
		return
	}
	// 远程推送只支持RTSP推RTMP格式
	// 指令 ffmpeg -rtsp_transport tcp -i rtsp://192.168.1.210:554/av0_0 -c:v libx264 \
	// -preset fast -tune zerolatency -ar 44100 -f flv rtmp://112.5.155.64:10110/live/123
	if vc.mainConfig.OutputMode == __OUTPUT_REMOTE_STREAM_SERVER {
		glogger.GLogger.Info("Start FFMPEG ffmpegProcess with: REMOTE_STREAM_SERVER Mode")
		var cmd *exec.Cmd
		// 本地摄像头推向远程服务器
		if vc.mainConfig.InputMode == __INPUT_LOCAL_CAMERA {
			var paramsVideo []string
			var deviceName string
			if runtime.GOOS == "windows" {
				deviceName = fmt.Sprintf("video=\"%s\"", inputUrl)
			}
			if runtime.GOOS == "linux" {
				deviceName = fmt.Sprintf("video=%s", inputUrl)
			}
			if runtime.GOOS == "windows" {
				paramsVideo = []string{
					"-err_detect",
					"ignore_err",
					"-hide_banner",
					"-f", "dshow", // windows下特有的DirectX加速引擎
					"-i", deviceName,
					"-q", "5",
					"-c:v", "libx264",
					"-preset", "veryfast",
					"-tune", "zerolatency",
					"-f", "flv",
					pushAddr,
				}
			} else {
				paramsVideo = []string{
					"-err_detect",
					"ignore_err",
					"-hide_banner",
					"-i", deviceName,
					"-q", "5",
					"-c:v", "libx264",
					"-preset", "veryfast",
					"-tune", "zerolatency",
					"-f", "flv",
					pushAddr,
				}
			}
			cmd = concatFFmpegCommand(paramsVideo)
		}
		// 远程拉RTSP推向RTMP
		if vc.mainConfig.InputMode == __INPUT_REMOTE_STREAM_RTSP {
			params := []string{
				"-rtsp_transport", "tcp",
				"-i", inputUrl,
				"-c:v", "libx264",
				"-preset", "fast",
				"-tune", "zerolatency",
				"-ar", "44100",
				"-f", "flv",
				pushAddr,
			}
			cmd = concatFFmpegCommand(params)
		}
		if cmd == nil {
			glogger.GLogger.Error(fmt.Errorf("InputMode Not Supported:" + vc.mainConfig.InputMode))
			return
		}
		glogger.GLogger.Debug("Start FFMPEG with command line: ", cmd.String())
		// 启动 FFmpeg 推流
		vc.ffmpegProcess = cmd
		if output, err1 := cmd.CombinedOutput(); err1 != nil {
			glogger.GLogger.Error("Combined Output error: ", err1, ", output: ", string(output))
		}
		return
	}
	glogger.GLogger.Info("stop Video Stream Endpoint:", inputUrl)
}
func concatFFmpegCommand(params []string) *exec.Cmd {
	var cmd *exec.Cmd
	// params = append([]string{"ffmpeg"}, params...)
	if runtime.GOOS == "windows" {
		bat := strings.Join(params, " ")
		cmd = exec.Command("powershell.exe", "-Command", "ffmpeg "+bat)
	} else {
		cmd = exec.Command("ffmpeg", params...)
	}
	return cmd
}

/*
*
* 停止进程
*
 */
func (vc *videoCamera) stopFFMPEGProcess() error {
	if vc.ffmpegProcess != nil {
		vc.ffmpegProcess.Process.Kill()
		vc.ffmpegProcess.Process.Signal(syscall.SIGTERM)
		glogger.GLogger.Info("FFMPEG Process stopped:", vc.ffmpegProcess.Process.Pid)
	}
	return nil
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

/*
*
  - MD5 URL= 播放源进行Hash后的字符串
    如果是RTSP，则Hash(URL)
    如果是Local，则Hash(deviceName)

*
*/
func calculateMD5(inputString string) string {
	hasher := md5.New()
	io.WriteString(hasher, inputString)
	hashBytes := hasher.Sum(nil)
	md5String := fmt.Sprintf("%x", hashBytes)
	return md5String
}
