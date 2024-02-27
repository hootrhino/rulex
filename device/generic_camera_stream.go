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
	__INPUT_REMOTE_STREAM_RTSP               = "REMOTE_STREAM_RTSP"   // 远程RTSP拉流
	__INPUT_LOCAL_CAMERA                     = "LOCAL_CAMERA"         // 本地摄像头
	__OUTPUT_LOCAL_STREAM_SERVER             = "LOCAL_STREAM_SERVER"  // RULEX 自带的FLVServer
	__OUTPUT_REMOTE_STREAM_SERVER            = "REMOTE_STREAM_SERVER" // 远程地址
	__default_push_to_internal_ws_server_url = "http://127.0.0.1:9400/stream/ffmpegPush?liveId="
)

// RTSP URL格式= rtsp://<username>:<password>@<ip>:<port>，
// 默认为本地拉流，推向Rulex自带的Stream server
type _MainConfig struct {
	InputMode    string `json:"inputMode" validate:"required"`    // 视频输入模式: REMOTE_STREAM_RTSP | LOCAL_CAMERA
	InputAddr    string `json:"inputAddr"`                        // 本地视频设备路径，在输入模式=LOCAL时生效
	OutputMode   string `json:"outputMode" validate:"required"`   // 输出模式: LOCAL_STREAM_SERVER | REMOTE_STREAM_SERVER
	OutputEncode string `json:"outputEncode" validate:"required"` // 输出编码: H264_STREAM
	OutputAddr   string `json:"outputAddr"`                       // 输出地址, 格式为: "Ip:Port",例如127.0.0.1:7890
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
		InputMode:    __INPUT_LOCAL_CAMERA,
		InputAddr:    "0",
		OutputMode:   __OUTPUT_LOCAL_STREAM_SERVER,
		OutputEncode: "H264_STREAM",
		OutputAddr:   __default_push_to_internal_ws_server_url,
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

// 启动, 设备的工作进程
// http://127.0.0.1:9400 RULEX自带的RTSP服务
func (vc *videoCamera) Start(cctx typex.CCTX) error {
	// 输出模式: LOCAL_STREAM_SERVER | REMOTE_STREAM_SERVER

	if vc.mainConfig.InputMode == __INPUT_LOCAL_CAMERA { // 本地USB摄像头
		// URL1 告诉RULEX要去哪里拉流
		if vc.mainConfig.OutputEncode == "H264_STREAM" {
			if vc.mainConfig.OutputMode == __OUTPUT_LOCAL_STREAM_SERVER {
				outputAddr := "http://127.0.0.1:9400/stream/ffmpegPush?liveId="
				pushUrl := outputAddr + calculateMD5(vc.mainConfig.InputAddr)
				go vc.startFFMPEGProcess(vc.mainConfig.InputAddr, pushUrl)
			}
			if vc.mainConfig.OutputMode == __OUTPUT_REMOTE_STREAM_SERVER {
				if !isValidRTSPAddress(vc.mainConfig.OutputAddr) {
					return fmt.Errorf("invalid OutputAddr Format:%s", vc.mainConfig.OutputAddr)
				}
				go vc.startFFMPEGProcess(vc.mainConfig.InputAddr, vc.mainConfig.OutputAddr)
			}
			// 告诉用户去那里拉流
			vc.status = typex.DEV_UP
			return nil
		}

	}
	if vc.mainConfig.InputMode == __INPUT_REMOTE_STREAM_RTSP { // RTSP
		// URL1 告诉RULEX要去哪里拉流
		if vc.mainConfig.OutputEncode == "H264_STREAM" {
			if vc.mainConfig.OutputMode == __OUTPUT_LOCAL_STREAM_SERVER {
				outputAddr := "http://127.0.0.1:9400/stream/ffmpegPush?liveId="
				pushUrl := outputAddr + calculateMD5(vc.mainConfig.InputAddr)
				go vc.startFFMPEGProcess(vc.mainConfig.InputAddr, pushUrl)
			}
			if vc.mainConfig.OutputMode == __OUTPUT_REMOTE_STREAM_SERVER {
				if !isValidRTSPAddress(vc.mainConfig.OutputAddr) {
					return fmt.Errorf("invalid OutputAddr Format:%s", vc.mainConfig.OutputAddr)
				}
				go vc.startFFMPEGProcess(vc.mainConfig.InputAddr, vc.mainConfig.OutputAddr)
			}
			// url2 告诉用户去哪里拉流
			vc.status = typex.DEV_UP
			return nil
		}
	}
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

func (vc *videoCamera) startFFMPEGProcess(inputUrl, pushAddr string) {
	defer func() {
		vc.status = typex.DEV_DOWN
	}()

	// LOCAL_STREAM_SERVER 向前端websocket输出
	// REMOTE_STREAM_SERVER 向远端推流(比如Monibuca)
	// ffmpeg -i rtsp://IP/av0_0 -c:v h264 -c:a aac -f rtsp rtsp://IP/live/test001
	if vc.mainConfig.OutputMode == __OUTPUT_LOCAL_STREAM_SERVER {
		glogger.GLogger.Info("Start FFMPEG ffmpegProcess with: LOCAL_STREAM_SERVER Mode")
		var cmd *exec.Cmd
		if vc.mainConfig.InputMode == __INPUT_LOCAL_CAMERA {
			var paramsVideo []string

			if runtime.GOOS == "windows" {
				paramsVideo = []string{
					"-err_detect", "ignore_err",
					"-hide_banner",
					"-f", "dshow", // windows下特有的DirectX加速引擎
					"-i", fmt.Sprintf("video=\"%s\"", inputUrl),
					"-q", "5",
					"-fflags", "nobuffer",
					"-c:v", "libx264",
					"-preset", "veryfast",
					"-tune", "zerolatency",
					"-f", "mpegts",
					"-an", pushAddr,
				}
			} else {
				paramsVideo = []string{
					"-err_detect", "ignore_err",
					"-hide_banner",
					"-i", fmt.Sprintf("video=%s", inputUrl),
					"-q", "5",
					"-fflags", "nobuffer",
					"-c:v", "libx264",
					"-preset", "veryfast",
					"-tune", "zerolatency",
					"-f", "mpegts",
					"-an", pushAddr,
				}
			}

			if runtime.GOOS == "windows" {
				bat := strings.Join(paramsVideo, " ")
				cmd = exec.Command("powershell.exe", "-Command", "ffmpeg "+bat)
			} else {
				cmd = exec.Command("ffmpeg", paramsVideo...)
			}
		}
		if vc.mainConfig.InputMode == __INPUT_REMOTE_STREAM_RTSP {
			paramsRtsp := []string{
				"-hide_banner",
				"-r", "24",
				"-rtsp_transport",
				"tcp",
				"-re",
				"-i",
				// rtsp://192.168.199.243:554/av0_0
				inputUrl,
				"-q",
				"5",
				"-f",
				"mpegts",
				"-fflags",
				"nobuffer",
				"-c:v",
				"libx264",
				"-an",
				// "-s",
				// "1920x1080",
				// http://127.0.0.1:9400/stream/ffmpegPush?liveId=147a6d7ae5a785f6e3ea90f25d36c63e
				pushAddr,
			}
			cmd = exec.Command("ffmpeg", paramsRtsp...)
		}
		if cmd == nil {
			glogger.GLogger.Error(fmt.Errorf("no supported InputMode:" + vc.mainConfig.InputMode))
			return
		}
		glogger.GLogger.Info("Start FFMPEG with:", cmd.String())
		// 启动 FFmpeg 推流
		vc.ffmpegProcess = cmd
		if output, err1 := cmd.CombinedOutput(); err1 != nil {
			glogger.GLogger.Error("Combined Output error: ", err1, ", output: ", string(output))
		}

	}

	glogger.GLogger.Info("stop Video Stream Endpoint:", inputUrl)
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
