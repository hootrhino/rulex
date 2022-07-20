package driver

// 这是个通用串口驱动，主要是用来做主动采集和控制用。
// `#` 分隔符: 注意该驱动的消息内容不要包含 `#`, 因为已经将其作为数据结尾提交符号
//
import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"

	"github.com/goburrow/serial"
)

// 数据缓冲区,单位: 字节
// const max_BUFFER_SIZE = 1024 * 4 // 4KB

// var buffer = [max_BUFFER_SIZE]byte{}

//------------------------------------------------------------------------
// 内部函数
//------------------------------------------------------------------------
type uartDriver struct {
	state      typex.DriverState
	serialPort serial.Port
	ctx        context.Context
	In         *typex.InEnd
	RuleEngine typex.RuleX
	onRead     func([]byte)
	bufferSize int
	buffer     []byte
}

//
// 初始化一个驱动
//
func NewUartDriver(
	ctx context.Context,
	config serial.Config,
	in *typex.InEnd,
	e typex.RuleX,
	bufferSize int,
	onRead func([]byte)) (typex.XExternalDriver, error) {
	serialPort, err := serial.Open(&config)
	if err != nil {
		glogger.GLogger.Error("uartModuleSource start failed:", err)
		return nil, err
	}
	return &uartDriver{
		In:         in,
		RuleEngine: e,
		serialPort: serialPort,
		ctx:        ctx,
		buffer:     make([]byte, bufferSize),
		onRead:     onRead,
	}, nil
}

//
//
//
func (a *uartDriver) Init(map[string]string) error {
	a.state = typex.DRIVER_RUNNING
	return nil
}

func (a *uartDriver) Work() error {
	ticker := time.NewTicker(time.Duration(time.Microsecond * 400))
	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			return
		default:
			{
			}
		}
		acc := 0
		data := make([]byte, 1)
		for a.state == typex.DRIVER_RUNNING {
			<-ticker.C
			if _, err0 := a.serialPort.Read(data); err0 != nil {
				//
				// 有的串口因为CPU频率原因 超时属于正常情况, 所以不计为错误
				// 只需要重启一下就可
				if !strings.Contains(err0.Error(), "timeout") {
					glogger.GLogger.Error("error:", err0)
					a.Stop()
					return
				} else {
					continue
				}
			}
			//---------------------------------------------------------------------------
			// 如果配置了自定义回调，则启用, 并且跳过默认协议，否则自动执行 '#' 结束符协议
			//---------------------------------------------------------------------------
			if a.onRead != nil {
				a.onRead(data)
				continue
			}
			//---------------------------------------------------------------------------
			// # 分隔符: 注意该驱动的消息内容不要包含 #, 因为已经将其作为数据结尾提交符号
			//---------------------------------------------------------------------------
			if data[0] == '#' {
				// glogger.GLogger.Info("bytes => ", string(buffer[:acc]), buffer[:acc], acc)
				a.RuleEngine.WorkInEnd(a.In, string(a.buffer[1:acc]))
				// 重新初始化缓冲区
				for i := 0; i < acc-1; i++ {
					a.buffer[i] = 0
				}
				data[0] = 0
				acc = 0
			}
			// 此处是为了过滤空行以及制表符
			if (data[0] != 0) && (data[0] != '\r') && (data[0] != '\n') {
				if acc <= a.bufferSize {
					a.buffer[acc] = data[0]
					acc += 1
				} else {
					glogger.GLogger.Errorf("data length exceed maximum buffer size limit: %v", a.bufferSize)
				}

			}
		}
	}(a.ctx)
	return nil

}
func (a *uartDriver) State() typex.DriverState {
	return a.state

}
func (a *uartDriver) Stop() error {
	a.state = typex.DRIVER_STOP
	return a.serialPort.Close()
}

func (a *uartDriver) Test() error {
	if a.serialPort == nil {
		return errors.New("serialPort is nil")
	}
	_, err := a.serialPort.Write([]byte("\r\n"))
	return err

}

//
func (a *uartDriver) Read(b []byte) (int, error) {
	return a.serialPort.Read(b)
}

//
func (a *uartDriver) Write(b []byte) (int, error) {
	n, err := a.serialPort.Write(b)
	if err != nil {
		glogger.GLogger.Error(err)
		return 0, err
	} else {
		return n, nil
	}

}
func (a *uartDriver) DriverDetail() typex.DriverDetail {
	return typex.DriverDetail{
		Name:        "General Uart Driver",
		Type:        "UART",
		Description: "A General Uart Driver Can Be Used For Common UART Device",
	}
}
