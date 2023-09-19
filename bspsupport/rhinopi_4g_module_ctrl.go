// Copyright (C) 2023 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package archsupport

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hootrhino/rulex/glogger"
)

/*
*
  - 初始化4G模组
    echo -e "AT+QCFG=\"usbnet\",1\r\n" >/dev/ttyUSB1  //驱动模式
    echo -e "AT+QNETDEVCTL=3,1,1\r\n" >/dev/ttyUSB1   //自动拨号
    echo -e "AT+QCFG=\"nat\",1 \r\n" >/dev/ttyUSB1    //网卡模式
    echo -e "AT+CFUN=1,1\r\n" >/dev/ttyUSB1           //重启
*/
const (
	__DRIVER_MODE_AT_CMD = "AT+QCFG=\"usbnet\",1\r\n"
	__DIAL_AT_CMD        = "AT+QNETDEVCTL=3,1,1\r\n"
	__NET_MODE_AT_CMD    = "AT+QCFG=\"nat\",1 \r\n"
	__RESET_AT_CMD       = "AT+CFUN=1,1\r\n"
	__SAVE_CONFIG        = "AT&W\r\n"
	__TURN_OFF_ECHO      = "ATE0\r\n"
	__ATTimeout          = 300 //ms
)

func init() {
	env := os.Getenv("ARCHSUPPORT")
	if env == "EEKITH3" {
		__RhinoPiInit4G()
	}
}

/*
*
  - 获取信号: +CSQ: 39,99
  - 0：没有信号。
  - 1-9：非常弱的信号，可能无法建立连接。
  - 10-14：较弱的信号，但可能可以建立连接。
  - 15-19：中等强度的信号。
  - 20-31：非常强的信号，信号质量非常好。
    RhinoPiGet4GCSQ: 返回值代表信号格
*/
func RhinoPiGet4GCSQ() int {
	csq := __Get4GCSQ()
	if csq == 0 {
		return 0
	}
	if csq > 0 && csq <= 9 {
		return 1
	}
	if csq > 9 && csq <= 14 {
		return 2

	}
	if csq > 15 && csq <= 19 {
		return 3

	}
	if csq > 19 && csq <= 31 {
		return 4
	}
	return 0
}
func __Get4GCSQ() int {
	result, err := __EC200A_AT("AT+CSQ\r\n", __ATTimeout)
	if err != nil {
		glogger.GLogger.Error(err)
		return 0
	}

	for _, v := range result {
		if v[:6] == "+CSQ: " {
			parts := strings.Split(string(v[6:]), ",")
			if len(parts) == 2 {
				v, err := strconv.Atoi(parts[0])
				if err != nil {
					return 0
				}
				return v
			}
		}
	}
	return 0
}

/*
*
* 初始化4G模组
*
 */
func __RhinoPiInit4G() {
	if err := turnOffEcho(); err != nil {
		glogger.GLogger.Error("RhinoPiInit4G turnOffEcho error:", err)
		return
	}
	glogger.GLogger.Info("RhinoPiInit4G turnOffEcho ok.")
	if err := setDriverMode(); err != nil {
		glogger.GLogger.Error("RhinoPiInit4G setDriverMode error:", err)
		return
	}
	glogger.GLogger.Info("RhinoPiInit4G setDriverMode ok.")
	if err := setDial(); err != nil {
		glogger.GLogger.Error("RhinoPiInit4G setDial error:", err)
		return
	}
	glogger.GLogger.Info("RhinoPiInit4G setDial ok.")
	if err := setNetMode(); err != nil {
		glogger.GLogger.Error("RhinoPiInit4G setNetMode error:", err)
		return
	}
	glogger.GLogger.Info("RhinoPiInit4G setNetMode ok.")
	if err := resetCard(); err != nil {
		glogger.GLogger.Error("RhinoPiInit4G resetCard error:", err)
		return
	}
	glogger.GLogger.Info("RhinoPiInit4G resetCard ok.")

}
func turnOffEcho() error {
	return __ExecuteAT(__TURN_OFF_ECHO)
}
func setDriverMode() error {
	return __ExecuteAT(__DRIVER_MODE_AT_CMD)
}
func setDial() error {
	return __ExecuteAT(__DIAL_AT_CMD)
}
func setNetMode() error {
	return __ExecuteAT(__NET_MODE_AT_CMD)
}
func resetCard() error {
	return __ExecuteAT(__RESET_AT_CMD)
}
func __ExecuteAT(cmd string) error {
	result, err := __EC200A_AT(cmd, __ATTimeout)
	if err != nil {
		return err
	}
	_, err1 := __EC200A_AT(__SAVE_CONFIG, __ATTimeout)
	if err1 != nil {
		return err
	}
	if !atOK(result) {
		return err
	}
	return nil
}

/*
*
  - EC200A 系列AT指令封装
    指令格式：AT+<?>\r\n
    指令返回值：\r\nCMD\r\n\r\nOK\r\n

解析结果

	[

		"",
		"CMD",
		"",
		"OK",
		""

	]

*
*/
func __EC200A_AT(command string, timeout time.Duration) ([]string, error) {
	// 打开设备文件以供读写
	device := "/dev/ttyUSB1"
	file, err := os.OpenFile(device, os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 写入AT指令
	_, err = file.WriteString(command + "\r\n")
	if err != nil {
		return nil, err
	}
	file.Sync()

	// 设置读取超时时间
	readTimeout := time.After(timeout)

	// 读取响应并以"\r\n"分割
	scanner := bufio.NewScanner(file)
	scanner.Split(scanCRLF)

	var responseLines []string
	for {
		select {
		case <-readTimeout:
			return responseLines, nil
		default:
			if scanner.Scan() {
				line := scanner.Text()
				// 过滤掉换行
				if line == "\n\r" {
					continue
				}
				responseLines = append(responseLines, line)
			} else if scanner.Err() != nil {
				return nil, scanner.Err()
			} else {
				// 所有行已读取完毕
				return responseLines, nil
			}
		}
	}
}

// 自定义分割函数，用于将输入按"\r\n"分割
func scanCRLF(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := strings.Index(string(data), "\r\n"); i >= 0 {
		return i + 2, data[:i+2], nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}

/*
*
* 判断是否成功
*
 */
func atOK(parts []string) bool {
	for _, part := range parts {
		if part == "OK" {
			return true
		}
		if part == "ERROR" {
			return false
		}
	}
	return false
}
