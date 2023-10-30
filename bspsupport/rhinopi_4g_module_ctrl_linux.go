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
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
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
	__CSQ                = "AT+CSQ\r\n"
	__SAVE_CONFIG        = "AT&W\r\n"
	__TURN_OFF_ECHO      = "ATE0\r\n"
	__CURRENT_COPS_CMD   = "AT+COPS?\r\n"
	__GET_ICCID_CMD      = "AT+QCCID\r\n"
	__ATTimeout          = 300 * time.Millisecond //ms
	__USB_4GDEV          = "/dev/ttyUSB1"
)

func init() {
	env := os.Getenv("ARCHSUPPORT")
	if env == "EEKITH3" {
		fmt.Println("RhinoPi Init 4G")
		__RhinoPiInit4G()
		fmt.Println("RhinoPi Init 4G Ok.")
	}
}

/*
*
* APN 配置, 参考文档: Quectel_LTE_Standard(A)系列_TCP(IP)_应用指导_V1.4.pdf-2.3.2章节
--

AT+QICSGP=<contextID>[,<context_
type>,<APN>[,<username>,<passwo
rd>)[,<authentication>[,<CDMA_pw
d>]]]]

AT+QICSGP=1 //查询场景 1 配置。
+QICSGP: 1,"","","",0
OK
AT+QICSGP=1,1,"UNINET","","",1 //配置场景 1，APN 配置为"UNINET"（中国联通）。
OK\ERROR
*/
func RhinoPiGetAPN() (string, error) {
	return __EC200A_AT("AT+QICSGP=1", __ATTimeout)
}

// 场景恒等于1
func RhinoPiSetAPN(ptype int, apn, username, password string, auth, cdmaPwd int) (string, error) {
	return __EC200A_AT(fmt.Sprintf(`AT+QICSGP=1,%d,"%s","%s","%s",%d,%d`,
		ptype, apn, username, password, auth, cdmaPwd), __ATTimeout)
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
	return __Get4GCSQ()
}

/*
*
* 获取运营商
+COPS:
(2,"CHINA MOBILE","CMCC","46000",7),
(1,"CHINA MOBILE","CMCC","46000",0),
(3,"CHN-UNICOM","UNICOM","46001",7),
(3,"CHN-CT","CT","46011",7),
(1,"460 15","460 15","46015",7),,(0-4),(0-2)
*/
func RhinoPiGetCOPS() (string, error) {
	// +COPS: 0,0,"CHINA MOBILE",7
	// +COPS: 0,0,"CHIN-UNICOM",7
	return __EC200A_AT(__CURRENT_COPS_CMD, __ATTimeout)
}
func RhinoPiRestart4G() (string, error) {
	return __EC200A_AT(__RESET_AT_CMD, __ATTimeout)
}

/*
*
* 获取ICCID, 用户查询电话卡号
* +QCCID: 89860025128306012474
 */
func RhinoPiGetICCID() (string, error) {
	return __EC200A_AT(__GET_ICCID_CMD, __ATTimeout)

}
func __Get4GCSQ() int {
	csq := 0
	file, err := os.OpenFile(__USB_4GDEV, os.O_RDWR, os.ModePerm)
	if err != nil {
		return csq
	}
	defer file.Close()
	_, err = file.WriteString(__CSQ)
	if err != nil {
		return csq
	}
	// 4G 模组的绝大多数指令都是100毫秒
	timeout := 300 * time.Millisecond
	buffer := [1]byte{}
	var responseData []byte
	b1 := 0
	for {
		if b1 == 4 {
			break
		}
		deadline := time.Now().Add(timeout)
		file.SetReadDeadline(deadline)
		n, err := file.Read(buffer[:])
		if err != nil {
			if err == io.EOF {
				break
			} else {
				break
			}
		}
		if n > 0 {
			if buffer[0] == 10 {
				b1++
			}
			if buffer[0] != 10 {
				responseData = append(responseData, buffer[0])
			}
		}
	}
	if len(responseData) > 6 {
		// +CSQ: 30,99
		response := string(responseData[6:])
		parts := strings.Split(response, ",")
		if len(parts) == 2 {
			v, err := strconv.Atoi(parts[0])
			if err == nil {
				csq = v
			}
		}
	}

	return csq
}

/*
*
* 初始化4G模组
*
 */
func __RhinoPiInit4G() {
	if err := turnOffEcho(); err != nil {
		fmt.Println("RhinoPiInit4G turnOffEcho error:", err)
		return
	}
	fmt.Println("RhinoPiInit4G turnOffEcho ok.")
	if err := setDriverMode(); err != nil {
		fmt.Println("RhinoPiInit4G setDriverMode error:", err)
		return
	}
	fmt.Println("RhinoPiInit4G setDriverMode ok.")
	if err := setDial(); err != nil {
		fmt.Println("RhinoPiInit4G setDial error:", err)
		return
	}
	fmt.Println("RhinoPiInit4G setDial ok.")
	if err := setNetMode(); err != nil {
		fmt.Println("RhinoPiInit4G setNetMode error:", err)
		return
	}
	fmt.Println("RhinoPiInit4G setNetMode ok.")
	if err := resetCard(); err != nil {
		fmt.Println("RhinoPiInit4G resetCard error:", err)
		return
	}
	fmt.Println("RhinoPiInit4G resetCard ok.")

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
	_, err0 := __EC200A_AT(cmd, __ATTimeout)
	if err0 != nil {
		return err0
	}
	_, err1 := __EC200A_AT(__SAVE_CONFIG, __ATTimeout)
	if err1 != nil {
		return err1
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
func __EC200A_AT(command string, timeout time.Duration) (string, error) {
	// 打开设备文件以供读写
	file, err := os.OpenFile(__USB_4GDEV, os.O_RDWR, os.ModePerm)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 写入AT指令
	_, err = file.WriteString(command)
	if err != nil {
		return "", err
	}
	buffer := [1]byte{}
	var responseData []byte
	b1 := 0
	for {
		if b1 == 4 {
			break
		}
		deadline := time.Now().Add(timeout)
		file.SetReadDeadline(deadline)
		n, err := file.Read(buffer[:])
		if err != nil {
			return "", err
		}
		if n > 0 {
			if buffer[0] == 10 {
				b1++
			}
			if buffer[0] != 10 {
				responseData = append(responseData, buffer[0])
			}
		}
	}
	return string(responseData), nil
}
