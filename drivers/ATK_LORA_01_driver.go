package gpio

import (
	"context"

	"github.com/ngaut/log"
	"github.com/tarm/serial"
)

//------------------------------------------------------------------------
// 内部函数
//------------------------------------------------------------------------

//
// 串口数据写入
//
func write(a *ATK_LORA_01Driver, k string) (error, string) {
	_, err := a.serialPort.Write([]byte(k + "\r\n"))
	if err != nil {
		return err, ""
	}
	for {
		response := make([]byte, 4)
		size, err := a.serialPort.Read(response)
		if err != nil {
			return err, ""
		}
		if size > 0 {
			return nil, string(response)
		}
	}
}

//------------------------------------------------------------------------

//
// 正点原子的 Lora 模块封装
//
type ATK_LORA_01Driver struct {
	serialPort *serial.Port
	channel    chan bool
	ctx        context.Context
}

//
// 初始化一个驱动
//
func NewATK_LORA_01Driver(serialPort *serial.Port) *ATK_LORA_01Driver {
	m := new(ATK_LORA_01Driver)
	m.channel = make(chan bool)
	m.serialPort = serialPort
	m.ctx = context.Background()
	return m
}

//
//
//
func (a *ATK_LORA_01Driver) Init() error {
	go func(context.Context) {
		log.Debug("ATK LORA 01 Driver Start Listening...")
		for {
			select {
			case <-a.ctx.Done():
				return
			default:
				{
					response := make([]byte, 16)
					_, err := a.serialPort.Read(response)
					if err != nil {
						a.Stop()
						return
					} else {
						log.Debug("SerialPort Received:", response)
					}
				}
			}

		}
	}(a.ctx)
	return nil
}
func (a *ATK_LORA_01Driver) Stop() error {
	a.ctx.Done()
	return nil
}

// -----------------------------
// AT\r\n
// -----------------------------

func (a *ATK_LORA_01Driver) Test() (error, string) {
	return write(a, "AT\r\n")
}

// 设置命令回显
// 获取参数
func (a *ATK_LORA_01Driver) GetProperty(k string) (error, string) {
	return write(a, k)
}

// 设置命令回显
// -----------------------------

func (a *ATK_LORA_01Driver) SetEcho() (error, string) {
	return write(a, "")
}

// 重置参数
// -----------------------------

func (a *ATK_LORA_01Driver) Reset() (error, string) {
	return write(a, "")
}

// 保存参数
// -----------------------------

func (a *ATK_LORA_01Driver) SaveConfig() (error, string) {
	return write(a, "")
}

// 恢复出厂
// -----------------------------

func (a *ATK_LORA_01Driver) RevoverFactory() (error, string) {
	return write(a, "")
}

// 设置地址
// -----------------------------

func (a *ATK_LORA_01Driver) SetAddr(k string, v string) (error, string) {
	return write(a, "")
}

// 设置功率
// -----------------------------

func (a *ATK_LORA_01Driver) SetPower(power int) (error, string) {
	return write(a, "")
}

// 设置WMode
// -----------------------------

func (a *ATK_LORA_01Driver) SetCWMode(mode int) (error, string) {
	return write(a, "")
}

// 设置SetTMode
// -----------------------------

func (a *ATK_LORA_01Driver) SetTMode(mode int) (error, string) {
	return write(a, "")
}

// 设置波特率
// -----------------------------

func (a *ATK_LORA_01Driver) SetRate(rate int, channel int) (error, string) {
	return write(a, "")
}

// 设置时间
// -----------------------------

func (a *ATK_LORA_01Driver) SetTime(time int) (error, string) {
	return write(a, "")
}

// 设置串口参数
// -----------------------------
// AT+UART=<bps>,<par>
// +UART:<bps>,<par> OK
// -----------------------------

func (a *ATK_LORA_01Driver) SetUart(bps int, par int) (error, string) {
	return write(a, "")
}
