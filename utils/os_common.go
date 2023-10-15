package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/hootrhino/rulex/glogger"
)

/*
*
* GetPwd
*
 */
func GetPwd() string {
	dir, err := os.Getwd()
	if err != nil {
		glogger.GLogger.Error(err)
		return ""
	}
	return dir
}

/*
*
* DEBUG使用
*
 */
func TraceMemStats() {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	var info [7]float64
	info[0] = float64(ms.HeapObjects)
	info[1] = BtoMB(ms.HeapAlloc)
	info[2] = BtoMB(ms.TotalAlloc)
	info[3] = BtoMB(ms.HeapSys)
	info[4] = BtoMB(ms.HeapIdle)
	info[5] = BtoMB(ms.HeapReleased)
	info[6] = BtoMB(ms.HeapIdle - ms.HeapReleased)

	for _, v := range info {
		fmt.Printf("%v,\t", v)
	}
	fmt.Println()
}
func BtoMB(bytes uint64) float64 {
	return float64(bytes) / 1024 / 1024
}

/*
*
* Byte to Mbyte
*
 */
func BToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

/*
*
* 获取操作系统发行版版本
runtime.GOARCH:

	386: 32-bit Intel/AMD x86 architecture
	amd64: 64-bit Intel/AMD x86 architecture
	arm: ARM architecture (32-bit)
	arm64: ARM architecture (64-bit)
	ppc64: 64-bit PowerPC architecture
	ppc64le: 64-bit little-endian PowerPC architecture
	mips: MIPS architecture (32-bit)
	mips64: MIPS architecture (64-bit)
	s390x: IBM System z architecture (64-bit)
	wasm: WebAssembly architecture

runtime.GOOS:

	darwin: macOS
	freebsd: FreeBSD
	linux: Linux
	windows: Windows
	netbsd: NetBSD
	openbsd: OpenBSD
	plan9: Plan 9
	dragonfly: DragonFly BSD

*
*/
func GetOSDistribution() (string, error) {
	if runtime.GOOS == "windows" {
		return runtime.GOOS, nil
	}
	// Linux 有很多发行版, 目前特别要识别一下Openwrt
	if runtime.GOOS == "linux" {
		cmd := exec.Command("cat", "/etc/os-release")
		output, err := cmd.Output()
		if err != nil {
			return runtime.GOOS, err
		}
		osIssue := strings.ToLower(string(output))
		if strings.Contains((osIssue), "openwrt") {
			return "openwrt", nil
		}
		if strings.Contains((osIssue), "ubuntu") {
			return "ubuntu", nil
		}
		if strings.Contains((osIssue), "debian") {
			return "debian", nil
		}
		if strings.Contains((osIssue), "armbian") {
			return "armbian", nil
		}
		if strings.Contains((osIssue), "deepin") {
			return "deepin", nil
		}
	}
	return runtime.GOOS, nil
}

/*
*
* 获取Ubuntu的版本
*
 */
func GetUbuntuVersion() (string, error) {
	// lsb_release -ds -> Ubuntu 22.04.1 LTS
	cmd := exec.Command("lsb_release", "-ds")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	info := strings.ToLower(strings.TrimSpace(string(output)))
	if strings.Contains(info, "ubuntu") {
		if strings.Contains(info, "16.") {
			return "ubuntu16", nil
		}
		if strings.Contains(info, "18.") {
			return "ubuntu18", nil
		}
		if strings.Contains(info, "20.") {
			return "ubuntu20", nil
		}
		if strings.Contains(info, "22.") {
			return "ubuntu22", nil
		}
		if strings.Contains(info, "24.") {
			return "ubuntu24", nil
		}
	}
	return "", fmt.Errorf("unsupported OS:%s", info)
}

/*
*
* 检查命令是否存在
*
 */

func CommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}
