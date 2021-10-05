package driver

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
func write(a *ATK_LORA_01Driver, k string) (string, error) {
	_, err := a.serialPort.Write([]byte(k + "\r\n"))
	if err != nil {
		return "", err
	}
	for {
		response := make([]byte, 4)
		size, err := a.serialPort.Read(response)
		if err != nil {
			return "", err
		}
		if size > 0 {
			return string(response), nil
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
	// 初始化配置参数
	// GPIO set 22 High 模块配置模式需要拉高一个脚
	// AT+UART=7,0 配置串口波特率: 115200, 无校验
	write(a, "AT+UART=7,0")
	// AT+WLRATE=23,5  信道: 433Hz 功率：19.2kbps
	write(a, "AT+WLRATE=23,5")
	// AT+ADDR=0000000000000001 配置地址 1
	write(a, "AT+ADDR=0000000000000001")
	// GPIO set 22 Low 退出配置模式

	return nil
}
func (a *ATK_LORA_01Driver) Work() error {
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
						log.Debug("SerialPort Received:", string(response))
					}
				}
			}

		}
	}(a.ctx)
	return nil

}
func (a *ATK_LORA_01Driver) State() DriverState {
	return RUNNING

}
func (a *ATK_LORA_01Driver) Stop() error {
	a.ctx.Done()
	return nil
}

// -----------------------------
// AT\r\n
// -----------------------------

func (a *ATK_LORA_01Driver) Test() (string, error) {
	return write(a, "AT\r\n")
}

// 设置命令回显
// 获取参数
func (a *ATK_LORA_01Driver) GetProperty(k string) (string, error) {
	return write(a, k)
}

// 设置命令回显
// -----------------------------

func (a *ATK_LORA_01Driver) SetEcho() (string, error) {
	return write(a, "")
}

// 重置参数
// -----------------------------

func (a *ATK_LORA_01Driver) Reset() (string, error) {
	return write(a, "")
}

// 保存参数
// -----------------------------

func (a *ATK_LORA_01Driver) SaveConfig() (string, error) {
	return write(a, "")
}

// 恢复出厂
// -----------------------------

func (a *ATK_LORA_01Driver) RevoverFactory() (string, error) {
	return write(a, "")
}

// 设置地址
// -----------------------------

func (a *ATK_LORA_01Driver) SetAddr(k string, v string) (string, error) {
	return write(a, "")
}

// 设置功率
// -----------------------------

func (a *ATK_LORA_01Driver) SetPower(power int) (string, error) {
	return write(a, "")
}

// 设置WMode
// -----------------------------

func (a *ATK_LORA_01Driver) SetCWMode(mode int) (string, error) {
	return write(a, "")
}

// 设置SetTMode
// -----------------------------

func (a *ATK_LORA_01Driver) SetTMode(mode int) (string, error) {
	return write(a, "")
}

// 设置波特率
// -----------------------------

func (a *ATK_LORA_01Driver) SetRate(rate int, channel int) (string, error) {
	return write(a, "")
}

// 设置时间
// -----------------------------

func (a *ATK_LORA_01Driver) SetTime(time int) (string, error) {
	return write(a, "")
}

// 设置串口参数
// -----------------------------
// AT+UART=<bps>,<par>
// +UART:<bps>,<par> OK
// -----------------------------

func (a *ATK_LORA_01Driver) SetUart(bps int, par int) (string, error) {
	return write(a, "")
}
