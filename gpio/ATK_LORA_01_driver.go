package gpio

type ATK_LORA_01 struct {
}

//
// 测试模块是否可用
//
func Test() string {
	return "OK"
}

// 获取参数
func GetProperty(k string) string {
	// 查询型号
	if k == "AT+MODEL?" {

	}
	// 查村软件版本号
	if k == "AT+CGMR?" {

	}
	// 查询地址
	if k == "AT+ADDR=?" {

	}
	// 查询功率
	if k == "AT+TPOWER=?" {

	}
	// 查询工作模式
	if k == "AT+CWMODE" {

	}
	// 查询发送状态
	if k == "AT+TMODE=?" {

	}
	// 查询无线速率和信道
	if k == "AT+WLRATE=?" {

	}
	// 查询休眠时间
	if k == "AT+WLTIME=?" {

	}
	// 查询串口参数
	if k == "AT+UART=?" {

	}

	return "NO_SUCH_COMMAND"
}
func SetEcho() string {
	return "OK"
}
func Reset() string {
	return "OK"
}
func SaveConfig() string {
	return "OK"
}
func RevoverFactory() string {
	return "OK"
}
func SetAddr(k string, v string) string {
	return "OK"
}
func SetPower(power int) string {
	return "OK"
}
func SetCWMode(mode int) string {
	return "OK"
}
func SetTMode(mode int) string {
	return "OK"
}
func SetRate(rate int, channel int) string {
	return "OK"
}
func SetTime(time int) string {
	return "OK"
}

// 设置串口参数
// AT+UART=<bps>,<par>
// +UART:<bps>,<par> OK
//
func SetUart(bps int, par int) string {
	return "OK"
}
