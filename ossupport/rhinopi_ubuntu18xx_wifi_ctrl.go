package ossupport

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/hootrhino/rulex/glogger"
)

/*
*
* WIFI 控制
nmcli device wifi rescan
nmcli device wifi list

获取WIFI网卡：iw dev | awk '$1=="Interface"{print $2}'
扫描WIFI列表: iwlist wlx0cc6551c5026 scan | grep ESSID | awk -F: '{print $2}' | sed 's/"//g'
*
*/
func ScanWIFIWithNmcli() ([]string, error) {

	wifiListReturn := []string{}
	var errReturn error
	finished := make(chan bool)

	go func() {
		{
			// 第一遍先扫描手边的 WIFI SSID
			cmd := exec.Command("sh", "-c", "nmcli device wifi rescan")
			output, err := cmd.CombinedOutput()
			if err != nil {
				errReturn = fmt.Errorf("scan WIFI error:%s,%s", string(output), err)
				return
			}
			glogger.GLogger.Debug("ScanWIFIWithNmcli:", cmd.String(), " OutPut:", string(output))

		}
		WFace := ""
		{
			cmd := exec.Command("sh", "-c", `iw dev | awk '$1=="Interface"{print $2}'`)
			output, err := cmd.CombinedOutput()
			if err != nil {
				stringWithoutNewlines := strings.Replace(string(output), "\n", "", -1)
				errReturn = fmt.Errorf("get WLAN Interface error:%s,%s", stringWithoutNewlines, err)
				return
			}
			if len(output) > 0 {
				WFace = string(output)
			} else {
				errReturn = fmt.Errorf("get WLAN Interface error:%s,%s", string(output), err)
				return
			}
			glogger.GLogger.Debug("ScanWIFIWithNmcli:", cmd.String(), " OutPut:", string(output))

		}

		{
			shell := `iwlist %s scan | grep ESSID | awk -F: '{print $2}' | sed 's/"//g'`
			stringWithoutNewlines := strings.Replace(WFace, "\n", "", -1)
			cmd := exec.Command("sh", "-c", fmt.Sprintf(shell, stringWithoutNewlines))
			output, err := cmd.CombinedOutput()
			if err != nil {
				errReturn = fmt.Errorf("scan WIFI error:%s,%s", string(output), err)
				return
			}
			// AAA\nBBB\nCCC\n
			wifiListReturn = append(wifiListReturn, strings.Split(string(output), "\n")...)
			glogger.GLogger.Debug("ScanWIFIWithNmcli:", cmd.String(), " OutPut :", string(output))
		}
		finished <- true
	}()
	select {
	case <-time.After(10 * time.Second): // 超时时间6秒
		errReturn = fmt.Errorf("scan WIFI timeout")
		return wifiListReturn, errReturn
	case <-finished:
		return wifiListReturn, errReturn
	}
}
