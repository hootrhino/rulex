package device

import (
	"github.com/hootrhino/rulex/typex"
)

/*
*
* 一般来说不会有太多摄像头，默认都是0、1，到4已经能覆盖绝大多数设备
*
 */
var videoDevMap = map[string]int{
	"video0": 0,
}

// RTSP URL格式= rtsp://<username>:<password>@<ip>:<port>，
type _MainConfig struct {
	MaxThread  int    `json:"maxThread"`  // 最大连接数, 防止连接过多导致摄像头拉流失败
	InputMode  string `json:"inputMode"`  // 视频输入模式：RTSP | LOCAL
	Device     string `json:"device"`     // 本地视频设备路径，在输入模式=LOCAL时生效
	RtspUrl    string `json:"rtspUrl"`    // 远程视频设备地址，在输入模式=RTSP时生效
	OutputMode string `json:"outputMode"` // 输出模式：JPEG_STREAM | RTSP_STREAM
	OutputAddr string `json:"outputAddr"` // 输出地址, 格式为: "Ip:Port",例如127.0.0.1:7890
}

// 摄像头
type videoCamera struct {
	typex.XStatus
	status     typex.DeviceState
	mainConfig _MainConfig
}

func NewVideoCamera(e typex.RuleX) typex.XDevice {
	videoCamera := new(videoCamera)
	videoCamera.RuleEngine = e
	videoCamera.status = typex.DEV_DOWN
	videoCamera.mainConfig = _MainConfig{
		Device:     "video0",
		RtspUrl:    "rtsp://127.0.0.1",
		InputMode:  "LOCAL",
		OutputMode: "JPEG_STREAM",
	}
	return videoCamera
}

// 初始化 通常用来获取设备的配置
func (vc *videoCamera) Init(devId string, configMap map[string]interface{}) error {

	return nil
}

// 启动, 设备的工作进程
func (vc *videoCamera) Start(cctx typex.CCTX) error {
	vc.status = typex.DEV_UP
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
* 外部指令，未来可以实现一些对摄像头的操作，例如拍照等
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
	vc.CancelCTX()
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
