package driver

import (
	"context"
	"rulex/typex"
	"strings"
	"time"

	"github.com/goburrow/serial"
	"github.com/ngaut/log"
)

// 数据缓冲区,单位: 字节
const max_BUFFER_SIZE = 1024 * 2 // 4KB

var buffer = [max_BUFFER_SIZE]byte{}

//------------------------------------------------------------------------
// 内部函数
//------------------------------------------------------------------------

//
// 正点原子的 Lora 模块封装
//
type UartDriver struct {
	state      typex.DriverState
	serialPort serial.Port
	ctx        context.Context
	In         *typex.InEnd
	RuleEngine typex.RuleX
	Separator  uint8
}

//
// 初始化一个驱动
//
func NewUartDriver(
	serialPort serial.Port,
	in *typex.InEnd,
	e typex.RuleX,
	separator uint8) typex.XExternalDriver {
	return &UartDriver{
		In:         in,
		RuleEngine: e,
		serialPort: serialPort,
		ctx:        context.Background(),
	}
}

//
//
//
func (a *UartDriver) Init() error {
	a.state = typex.RUNNING
	return nil
}

func (a *UartDriver) SetState(state typex.DriverState) {
	a.state = state

}
func (a *UartDriver) Work() error {

	go func(ctx context.Context) {
		acc := 0
		data := make([]byte, 1)
		ticker := time.NewTicker(time.Duration(time.Microsecond * 400))
		for a.state == typex.RUNNING {
			<-ticker.C
			if _, err0 := a.serialPort.Read(data); err0 != nil {
				// 有的串口因为CPU频率原因 超时属于正常情况，所以不计为错误
				if !strings.Contains(err0.Error(), "timeout") {
					log.Error("error:", err0)
					a.Stop()
					return
				} else {
					continue
				}
			}
			// # 分隔符
			if data[0] == '#' {
				// log.Info("bytes => ", string(buffer[:acc]), buffer[:acc], acc)
				a.RuleEngine.PushQueue(typex.QueueData{
					In:   a.In,
					Out:  nil,
					E:    a.RuleEngine,
					Data: string(buffer[1:acc]),
				})
				// 重新初始化缓冲区
				for i := 0; i < acc-1; i++ {
					buffer[i] = 0
				}
				data[0] = 0
				acc = 0
			}

			if (data[0] != 0) && (data[0] != '\r') && (data[0] != '\n') {
				buffer[acc] = data[0]
				acc += 1
			}
		}
	}(a.ctx)
	return nil

}
func (a *UartDriver) State() typex.DriverState {
	return a.state

}
func (a *UartDriver) Stop() error {
	a.state = typex.STOP
	return a.serialPort.Close()
}

func (a *UartDriver) Test() error {
	return nil
}

//
func (a *UartDriver) Read(b []byte) (int, error) {
	return a.serialPort.Read(b)
}

//
func (a *UartDriver) Write(b []byte) (int, error) {
	n, err := a.serialPort.Write(b)
	if err != nil {
		log.Error(err)
		return 0, err
	} else {
		return n, nil
	}

}
func (a *UartDriver) DriverDetail() *typex.DriverDetail {
	return &typex.DriverDetail{
		Name:        "Generic Uart Driver",
		Type:        "UART",
		Description: "A Generic Uart Driver",
	}
}
