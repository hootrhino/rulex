package gpio

import "github.com/tarm/serial"

//------------------------------------------------------------------------
// 内部函数
//------------------------------------------------------------------------

//
//
//
func write(a *ATK_LORA_01Driver, k string) (error, string) {
	_, err := a.serialPort.Write([]byte(k + "\r\n"))
	if err != nil {
		return err, ""
	}
	for {
		response := make([]byte, 0)
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
}

//
// 初始化一个驱动
//
func NewATK_LORA_01Driver(serialPort *serial.Port) *ATK_LORA_01Driver {
	m := new(ATK_LORA_01Driver)
	m.serialPort = serialPort
	return m
}

func (a *ATK_LORA_01Driver) Test() (error, string) {
	return write(a, "AT\r\n")
}

// 获取参数
func (a *ATK_LORA_01Driver) GetProperty(k string) (error, string) {
	return write(a, k)
}
func (a *ATK_LORA_01Driver) SetEcho() (error, string) {
	return write(a, "")
}
func (a *ATK_LORA_01Driver) Reset() (error, string) {
	return write(a, "")
}
func (a *ATK_LORA_01Driver) SaveConfig() (error, string) {
	return write(a, "")
}
func (a *ATK_LORA_01Driver) RevoverFactory() (error, string) {
	return write(a, "")
}
func (a *ATK_LORA_01Driver) SetAddr(k string, v string) (error, string) {
	return write(a, "")
}
func (a *ATK_LORA_01Driver) SetPower(power int) (error, string) {
	return write(a, "")
}
func (a *ATK_LORA_01Driver) SetCWMode(mode int) (error, string) {
	return write(a, "")
}
func (a *ATK_LORA_01Driver) SetTMode(mode int) (error, string) {
	return write(a, "")
}
func (a *ATK_LORA_01Driver) SetRate(rate int, channel int) (error, string) {
	return write(a, "")
}
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
