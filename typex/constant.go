package typex

// Source State
type SourceState int

const (
	SOURCE_DOWN  SourceState = 0 // 此状态需要重启
	SOURCE_UP    SourceState = 1
	SOURCE_PAUSE SourceState = 2
	SOURCE_STOP  SourceState = 3
)

func (s SourceState) String() string {
	if s == 0 {
		return "DOWN"
	}
	if s == 1 {
		return "UP"
	}
	if s == 2 {
		return "PAUSE"
	}
	if s == 3 {
		return "STOP"
	}
	return "UnKnown State"

}

// Abstract driver interface
type DriverState int

const (
	// STOP 状态一般用来直接停止一个资源，监听器不需要重启
	DRIVER_STOP DriverState = 0
	// UP 工作态
	DRIVER_UP DriverState = 1
	// DOWN 状态是某个资源挂了，属于工作意外，需要重启
	DRIVER_DOWN DriverState = 2
)

type DeviceState int

const (
	// 设备故障
	DEV_DOWN DeviceState = 0
	// 设备启用
	DEV_UP DeviceState = 1
	// 暂停，这是个占位值，只为了和其他地方统一值,但是没用
	_ DeviceState = 2
	// 外部停止
	DEV_STOP DeviceState = 3
)

func (s DeviceState) String() string {
	if s == 0 {
		return "DOWN"
	}
	if s == 1 {
		return "UP"
	}
	if s == 2 {
		return "PAUSE"
	}
	if s == 3 {
		return "STOP"
	}
	return "ERROR"
}

/*
*
* 串口校验形式
*
 */
type Parity string

const (
	ODD  Parity = "O" // 奇校验
	EVEN Parity = "E" // 偶校验
	NONE Parity = "N" // 不校验
)
