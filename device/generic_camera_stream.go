package device

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"syscall"

	"github.com/hootrhino/rulex/core"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
)

type _MainConfig struct {
	InputMode   string `json:"inputMode" validate:"required"` // 视频输入模式：RTSP | LOCAL
	LocalDevice string `json:"device"`                        // 本地视频设备路径，在输入模式=LOCAL时生效
	// RTSP URL格式= rtsp://<username>:<password>@<ip>:<port>，
	RtspUrl    string `json:"rtspUrl"`                        // 远程视频设备地址，在输入模式=RTSP时生效
	OutputMode string `json:"outputMode" validate:"required"` // 输出模式：JPEG_STREAM | H264_STREAM
	outputAddr string // 输出地址, 格式为: "Ip:Port",例如127.0.0.1:7890
	playAddr   string // 输出地址, 格式为: "Ip:Port",例如127.0.0.1:7890
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
		LocalDevice: "dev/video0",
		RtspUrl:     "rtsp://127.0.0.1",
		InputMode:   "LOCAL",
		OutputMode:  "H264_STREAM",
		outputAddr:  "http://127.0.0.1:9400/stream/ffmpegPush?liveId=",
		playAddr:    "ws://127.0.0.1:9400/stream/live?liveId=",
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
	if vc.mainConfig.InputMode == "RTSP" {
		if !isValidRTSPAddress(vc.mainConfig.RtspUrl) {
			return fmt.Errorf("invalid RtspUrl Format:%s", vc.mainConfig.RtspUrl)
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
	if vc.mainConfig.InputMode == "LOCAL" { // 本地USB摄像头
		// URL1 告诉RULEX要去哪里拉流
		if vc.mainConfig.OutputMode == "H264_STREAM" {
			url1 := vc.mainConfig.outputAddr + calculateMD5(vc.mainConfig.LocalDevice)
			go vc.startFFMPEGProcess(vc.mainConfig.LocalDevice, url1)
			// 告诉用户去那里拉流
			url2 := vc.mainConfig.playAddr + calculateMD5(vc.mainConfig.LocalDevice)
			vc.mainConfig.outputAddr = url2
			vc.status = typex.DEV_UP
		}

	}
	if vc.mainConfig.InputMode == "RTSP" { // RTSP
		// URL1 告诉RULEX要去哪里拉流
		if vc.mainConfig.OutputMode == "H264_STREAM" {
			url1 := vc.mainConfig.outputAddr + calculateMD5(vc.mainConfig.RtspUrl)
			go vc.startFFMPEGProcess(vc.mainConfig.RtspUrl, url1)
			// url2 告诉用户去哪里拉流
			url2 := vc.mainConfig.playAddr + calculateMD5(vc.mainConfig.RtspUrl)
			vc.mainConfig.outputAddr = url2
			vc.status = typex.DEV_UP
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
	vc.status = typex.DEV_STOP
	if vc.CancelCTX != nil {
		vc.CancelCTX()
	}
	vc.stopFFMPEGProcess()
}

func (vc *videoCamera) Property() []typex.DeviceProperty {
	return []typex.DeviceProperty{}
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

func (vc *videoCamera) StartPullRtsp(rtspUrl, pushAddr string) {
	vc.startFFMPEGProcess(rtspUrl, pushAddr)
}
func (vc *videoCamera) StartLocalVideo(localDevice, pushAddr string) {
	vc.startFFMPEGProcess(localDevice, pushAddr)
}
func (vc *videoCamera) startFFMPEGProcess(rtspUrl, pushAddr string) {
	defer func() {
		vc.status = typex.DEV_DOWN
	}()
	paramsVideo := []string{
		"-f", "dshow",
		"-i", fmt.Sprintf("video=%s", rtspUrl),
		"-c:v", "libx264",
		"-preset", "veryfast",
		"-tune", "zerolatency",
		"-f", "flv",
		pushAddr,
	}

	paramsRtsp := []string{
		// rtsp://192.168.199.243:554/av0_0
		"-rtsp_transport",
		"tcp",
		"-re",
		"-i",
		fmt.Sprintf("'%s'", rtspUrl),
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
	var cmd *exec.Cmd
	if vc.mainConfig.InputMode == "LOCAL" {
		cmd = exec.Command("ffmpeg", paramsVideo...)
	}
	if vc.mainConfig.InputMode == "RTSP" {
		cmd = exec.Command("ffmpeg", paramsRtsp...)
	}
	if cmd == nil {
		glogger.GLogger.Error(fmt.Errorf("no supported InputMode:" + vc.mainConfig.InputMode))
		return
	}
	glogger.GLogger.Info("Start FFMPEG ffmpegProcess with:", cmd.String())
	// 启动 FFmpeg 推流
	if err := cmd.Start(); err != nil {
		glogger.GLogger.Error(err)
		return
	}
	if core.GlobalConfig.AppDebugMode {
		inOut := wsInOut{}
		cmd.Stdin = nil
		cmd.Stdout = inOut
		cmd.Stderr = inOut
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	vc.ffmpegProcess = cmd
	// 等待 FFmpeg 进程完成
	if err := vc.ffmpegProcess.Wait(); err != nil {
		output, _ := cmd.CombinedOutput()
		glogger.GLogger.Error(err, string(output))
		return
	}
	glogger.GLogger.Info("stop Video Stream Endpoint:", rtspUrl)
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
* MD5 URL
*
 */
func calculateMD5(inputString string) string {
	hasher := md5.New()
	io.WriteString(hasher, inputString)
	hashBytes := hasher.Sum(nil)
	md5String := fmt.Sprintf("%x", hashBytes)
	return md5String
}
