package driver

//
// `#`` 分隔符: 注意该驱动的消息内容不要包含 `#`, 因为已经将其作为数据结尾提交符号
//
import (
	"context"
	"errors"
	"rulex/typex"
	"strings"
	"time"

	"github.com/goburrow/serial"
	"github.com/ngaut/log"
)

// 数据缓冲区,单位: 字节
const max_BUFFER_SIZE = 1024 * 4 // 4KB

var buffer = [max_BUFFER_SIZE]byte{}

//------------------------------------------------------------------------
// 内部函数
//------------------------------------------------------------------------

//
// 正点原子的 Lora 模块封装
//
type uartDriver struct {
	state      typex.DriverState
	serialPort serial.Port
	ctx        context.Context
	In         *typex.InEnd
	RuleEngine typex.RuleX
}

//
// 初始化一个驱动
//
func NewUartDriver(
	serialPort serial.Port,
	in *typex.InEnd,
	e typex.RuleX) typex.XExternalDriver {
	return &uartDriver{
		In:         in,
		RuleEngine: e,
		serialPort: serialPort,
		ctx:        context.Background(),
	}
}

//
//
//
func (a *uartDriver) Init() error {
	a.state = typex.RUNNING
	return nil
}

func (a *uartDriver) SetState(state typex.DriverState) {
	a.state = state

}
func (a *uartDriver) Work() error {

	go func(ctx context.Context) {
		acc := 0
		data := make([]byte, 1)
		ticker := time.NewTicker(time.Duration(time.Microsecond * 400))
		for a.state == typex.RUNNING {
			<-ticker.C
			if _, err0 := a.serialPort.Read(data); err0 != nil {
				// 有的串口因为CPU频率原因 超时属于正常情况, 所以不计为错误
				if !strings.Contains(err0.Error(), "timeout") {
					log.Error("error:", err0)
					a.Stop()
					return
				} else {
					continue
				}
			}
			//
			// # 分隔符: 注意该驱动的消息内容不要包含 #, 因为已经将其作为数据结尾提交符号
			//
			if data[0] == '#' {
				// log.Info("bytes => ", string(buffer[:acc]), buffer[:acc], acc)
				a.RuleEngine.Work(a.In, string(buffer[1:acc]))
				// 重新初始化缓冲区
				for i := 0; i < acc-1; i++ {
					buffer[i] = 0
				}
				data[0] = 0
				acc = 0
			}
			// 此处是为了过滤空行以及制表符
			if (data[0] != 0) && (data[0] != '\r') && (data[0] != '\n') {
				if acc <= max_BUFFER_SIZE {
					buffer[acc] = data[0]
					acc += 1
				} else {
					log.Error("max buffer reached!")
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
	a.state = typex.STOP
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
		log.Error(err)
		return 0, err
	} else {
		return n, nil
	}

}
func (a *uartDriver) DriverDetail() *typex.DriverDetail {
	return &typex.DriverDetail{
		Name:        "Generic Uart Driver",
		Type:        "UART",
		Description: "A Generic Uart Driver",
	}
}
