package device

import (
	"encoding/json"
	"fmt"
	"syscall"
	"time"
	"unsafe"

	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
)

const __IR_DEV = "/dev/input/event1"

/*
*
* 红外线接收模块
$ ir-keytable
Found /sys/class/rc/rc0/ (/dev/input/event1) with:

	Name: sunxi-ir
	Driver: sunxi-ir, table: rc-empty
	lirc device: /dev/lirc0
	Supported protocols: other lirc rc-5 rc-5-sz jvc sony nec sanyo mce_kbd rc-6 sharp xmp
	Enabled protocols: lirc nec
	bus: 25, vendor/product: 0001:0001, version: 0x0100
	Repeat delay = 500 ms, repeat period = 125 ms
*/
type IR struct {
	typex.XStatus
	status typex.DeviceState
	irFd   int
	// irFd       syscall.Handle windows
	RuleEngine typex.RuleX
}

func NewIRDevice(e typex.RuleX) typex.XDevice {
	uart := new(IR)
	uart.RuleEngine = e
	return uart
}

//  初始化
func (ird *IR) Init(devId string, configMap map[string]interface{}) error {
	ird.PointId = devId

	return nil
}

type timeval struct {
	Second  int32 `json:"second,omitempty"`
	USecond int32 `json:"uSecond,omitempty"`
}
type irInputEvent struct {
	Time  timeval `json:"-"`
	Type  uint16  `json:"-"`
	Code  uint16  `json:"code,omitempty"`
	Value int32   `json:"value,omitempty"`
}

func (v irInputEvent) String() string {
	b, _ := json.Marshal(v)
	return string(b)
}

// 启动
func (ird *IR) Start(cctx typex.CCTX) error {
	ird.Ctx = cctx.Ctx
	ird.CancelCTX = cctx.CancelCTX

	fd, err := syscall.Open("/dev/input/event1", syscall.O_RDONLY, 0777)
	if err != nil {
		fmt.Printf("device open failed\r\n")
		syscall.Close(fd)
		return err
	}
	ird.irFd = fd
	go func(ird *IR) {
		defer func() {
			syscall.Close(fd)
		}()
		buf := make([]byte, 1024)
		for {
			select {
			case <-ird.Ctx.Done():
				return
			default:
				{
				}
			}
			n, e := syscall.Read(fd, buf)
			if e != nil {
				glogger.GLogger.Error(e)
				continue
			}
			if n > 0 {
				event := irInputEvent{}
				_, err := syscall.Read(fd, (*[24]byte)(unsafe.Pointer(&event))[:])
				if err != nil {
					glogger.GLogger.Error(err)
				} else {
					ird.RuleEngine.WorkDevice(ird.Details(), event.String())
				}
			}
			time.Sleep(125 * time.Millisecond)
		}
	}(ird)
	ird.status = typex.DEV_UP
	return nil
}

/*
*
* 不支持读, 仅仅是个数据透传
*
 */
func (ird *IR) OnRead(cmd []byte, data []byte) (int, error) {
	return 0, fmt.Errorf("IR not support read data")
}

func (ird *IR) OnWrite(cmd []byte, b []byte) (int, error) {
	return 0, nil
}

// 设备当前状态
func (ird *IR) Status() typex.DeviceState {

	return typex.DEV_UP
}

// 停止设备
func (ird *IR) Stop() {
	ird.status = typex.DEV_DOWN
	if ird.CancelCTX != nil {
		ird.CancelCTX()
	}
	if ird.irFd != 0 {
		syscall.Close(ird.irFd)
	}

}

// 设备属性，是一系列属性描述
func (ird *IR) Property() []typex.DeviceProperty {
	return []typex.DeviceProperty{}
}

// 真实设备
func (ird *IR) Details() *typex.Device {
	return ird.RuleEngine.GetDevice(ird.PointId)
}

// 状态
func (ird *IR) SetState(status typex.DeviceState) {
	ird.status = status

}

// 驱动
func (ird *IR) Driver() typex.XExternalDriver {
	return nil
}

func (ird *IR) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}
func (ird *IR) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	return []byte{}, nil
}
