package device

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
	"github.com/hootrhino/rulex/utils"
	"gocv.io/x/gocv"
)

/*
*
* 一般来说不会有太多摄像头，默认都是0、1，到4已经能覆盖绝大多数设备
*
 */
var videoDevMap = map[string]int{
	"video0": 0,
	"video1": 1,
	"video2": 2,
	"video3": 3,
	"video4": 4,
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
	video      *gocv.VideoCapture
}

func NewVideoCamera(e typex.RuleX) typex.XDevice {
	videoCamera := new(videoCamera)
	videoCamera.RuleEngine = e
	videoCamera.status = typex.DEV_DOWN
	videoCamera.mainConfig = _MainConfig{
		MaxThread:  10,
		Device:     "video0",
		RtspUrl:    "rtsp://127.0.0.1",
		InputMode:  "LOCAL",
		OutputMode: "JPEG_STREAM",
		OutputAddr: "127.0.0.1:2599",
	}
	return videoCamera
}

// 初始化 通常用来获取设备的配置
func (vc *videoCamera) Init(devId string, configMap map[string]interface{}) error {
	vc.PointId = devId
	if err := utils.BindSourceConfig(configMap, &vc.mainConfig); err != nil {
		return err
	}

	// 拉取IP摄像头
	if vc.mainConfig.InputMode == "RTSP" {
		video, err := gocv.OpenVideoCapture(vc.mainConfig.RtspUrl)
		if err != nil {
			errMsg := fmt.Errorf("RtspUrl %v, open error: %v",
				vc.mainConfig.RtspUrl, err.Error())
			glogger.GLogger.Error(errMsg)
			return errMsg
		}
		defer video.Close()
	}
	// 从本地摄像头获取
	if vc.mainConfig.InputMode == "LOCAL" {
		if _, ok := videoDevMap[vc.mainConfig.Device]; !ok {
			errMsg := fmt.Errorf("video device: %v not exists", vc.mainConfig.Device)
			glogger.GLogger.Error(errMsg)
			return errMsg
		}
		video, err := gocv.OpenVideoCapture(videoDevMap[vc.mainConfig.Device])
		if err != nil {
			errMsg := fmt.Errorf("video device %v, open error: %v",
				vc.mainConfig.Device, err.Error())
			glogger.GLogger.Error(errMsg)
			return errMsg
		}
		defer video.Close()
	}
	return nil
}

// 启动, 设备的工作进程
func (vc *videoCamera) Start(cctx typex.CCTX) error {
	vc.Ctx = cctx.Ctx
	vc.CancelCTX = cctx.CancelCTX
	var err error
	//
	// 从远程摄像头拉流
	//
	if vc.mainConfig.InputMode == "RTSP" {
		vc.video, err = gocv.OpenVideoCapture(vc.mainConfig.RtspUrl)
	}
	//
	// 从本地摄像头拉流
	//
	if vc.mainConfig.InputMode == "LOCAL" {
		if _, ok := videoDevMap[vc.mainConfig.Device]; !ok {
			errMsg := fmt.Errorf("video device: %v not exists", vc.mainConfig.Device)
			glogger.GLogger.Error(errMsg)
			return errMsg
		}
		vc.video, err = gocv.OpenVideoCapture(videoDevMap[vc.mainConfig.Device])
	}
	if err != nil {
		errMsg := fmt.Errorf("video device %v, open error: %v",
			vc.mainConfig.Device, err.Error())
		glogger.GLogger.Error(errMsg)
		return errMsg
	}
	if vc.mainConfig.OutputMode == "JPEG_STREAM" {
		go serveHttpJPEGStream(vc)
	}
	if vc.mainConfig.OutputMode == "RTSP_STREAM" {
		go serveRtspStream(vc)
	}
	vc.status = typex.DEV_UP
	return nil
}

/*
*
* 提供RTSP接口来访问摄像头
*
 */
func serveRtspStream(videoCamera *videoCamera) {

}

/*
*
* 提供HTTP接口来访问摄像头
*
 */
func serveHttpJPEGStream(videoCamera *videoCamera) {
	defer videoCamera.video.Close()
	cvMat := gocv.NewMat()
	defer cvMat.Close()
	stream := NewMJPEGStream(videoCamera.mainConfig.MaxThread)
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", stream.ServeHTTP)
		serverOne := &http.Server{
			Addr:    videoCamera.mainConfig.OutputAddr,
			Handler: mux,
			BaseContext: func(l net.Listener) context.Context {
				return videoCamera.Ctx
			},
		}
		serverOne.ListenAndServe()
	}()
	errTimes := 0
	for {
		select {
		case <-videoCamera.Ctx.Done():
			return
		default:
		}
		if ok := videoCamera.video.Read(&cvMat); !ok {
			// 如果连续30帧失败，直接重启
			errTimes++
			if errTimes > 30 {
				videoCamera.status = typex.DEV_DOWN
				continue
			} else {
				continue
			}
		}
		if cvMat.Empty() {
			continue
		}
		cvBuf, err := gocv.IMEncode(".png", cvMat)
		if err != nil {
			glogger.GLogger.Error(err)
			continue
		}
		stream.UpdateJPEG(cvBuf.GetBytes())
	}

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

func (vc *videoCamera) SetState(s typex.DeviceState) {
	vc.status = s
}

func (vc *videoCamera) Driver() typex.XExternalDriver {
	return nil
}

func (vc *videoCamera) OnDCACall(_ string, _ string, _ interface{}) typex.DCAResult {
	return typex.DCAResult{}
}

//--------------------------------------------------------------
// HTTP 图片流
//--------------------------------------------------------------

type mJPEGStream struct {
	m             map[chan []byte]bool
	frame         []byte
	lock          sync.Mutex
	FrameInterval time.Duration
	maxThread     int
	currentThread int
}

// multipart/x-mixed-replace 一种固定写法
const boundaryWord = "MJPEGBOUNDARY"
const header = "\r\n" +
	"--" + boundaryWord + "\r\n" +
	"Content-Type: image/JPEG\r\n" +
	"Content-Length: %d\r\n" +
	"X-Timestamp: 0.000000\r\n" +
	"\r\n"

// serveHttpJPEGStream responds to HTTP requests with the MJPEG stream, implementing the http.Handler interface.
func (s *mJPEGStream) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.currentThread++
	if s.currentThread > s.maxThread {
		w.Write([]byte("Exceed MaxThread"))
		return
	}
	w.Header().Add("Content-Type", "multipart/x-mixed-replace;boundary="+boundaryWord)
	c := make(chan []byte)
	s.lock.Lock()
	s.m[c] = true
	s.lock.Unlock()
	for {
		time.Sleep(s.FrameInterval)
		b := <-c
		_, err := w.Write(b)
		if err != nil {
			break
		}
	}
	s.lock.Lock()
	delete(s.m, c)
	s.lock.Unlock()
	s.currentThread--
}

func (s *mJPEGStream) UpdateJPEG(JPEG []byte) {
	header := fmt.Sprintf(header, len(JPEG))
	if len(s.frame) < len(JPEG)+len(header) {
		s.frame = make([]byte, (len(JPEG)+len(header))*2)
	}
	copy(s.frame, header)
	copy(s.frame[len(header):], JPEG)
	s.lock.Lock()
	for c := range s.m {
		select {
		case c <- s.frame:
		default:
		}
	}
	s.lock.Unlock()
}

func NewMJPEGStream(mt int) *mJPEGStream {
	return &mJPEGStream{
		m:             make(map[chan []byte]bool),
		frame:         make([]byte, len(header)),
		FrameInterval: 50 * time.Millisecond,
		currentThread: 0,
		maxThread:     mt,
	}
}
