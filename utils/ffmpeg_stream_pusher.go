// Copyright (C) 2024 wwhai
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

package utils

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

/*
*
* 推本地的摄像头到远程流服务器
*ffmpeg.exe -hide_banner -r 24 -rtsp_transport tcp -re
-i video="USB2.0 PC CAMERA" -c:v mjpeg -f mjpeg
-headers "Content-Type:multipart/x-mixed-replace" "http://127.0.0.1:9401/jpeg_stream/push?liveId=123"
*/
func PushLocalCameraToJpegStreamServer(FromDevice, pushAddr string) *exec.Cmd {
	var deviceName string
	var paramsVideo []string
	if runtime.GOOS == "windows" {
		deviceName = fmt.Sprintf("video=\"%s\"", FromDevice)
	}
	if runtime.GOOS == "linux" {
		deviceName = fmt.Sprintf("video=%s", FromDevice)
	}
	if runtime.GOOS == "windows" {
		paramsVideo = []string{
			"-err_detect", "ignore_err",
			"-hide_banner",
			"-r", "24",
			"-f", "dshow", // windows下特有的DirectX加速引擎
			"-re",
			"-i", deviceName,
			"-c:v", "mjpeg",
			"-f", "mjpeg",
			"-fflags", "+genpts",
			"-preset", "veryfast",
			"-tune", "zerolatency",
			"-headers", "Content-Type:multipart/x-mixed-replace; boundary=MJPEG_BOUNDARY",
			pushAddr,
		}
	} else {
		paramsVideo = []string{
			"-err_detect", "ignore_err",
			"-hide_banner",
			"-r", "24",
			"-re",
			"-i", deviceName,
			"-c:v", "mjpeg",
			"-f", "mjpeg",
			"-fflags", "+genpts",
			"-preset", "veryfast",
			"-tune", "zerolatency",
			"-headers", "Content-Type:multipart/x-mixed-replace; boundary=MJPEG_BOUNDARY",
			pushAddr,
		}
	}
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		bat := strings.Join(paramsVideo, " ")
		cmd = exec.Command("powershell.exe", "-Command", "ffmpeg "+bat)
	} else {
		cmd = exec.Command("ffmpeg", paramsVideo...)
	}
	return cmd
}

/*
*
* 推送Rtsp流
ffmpeg.exe -hide_banner -r 24 -rtsp_transport tcp -re
-i rtsp://192.168.1.210:554/av0_0 -c:v mjpeg
-f mjpeg -headers "Content-Type:multipart/x-mixed-replace" "http://127.0.0.1:9401/jpeg_stream/push?liveId=123"
*/
func PushRtspToJpegStreamServer(From, ToUrl string) *exec.Cmd {
	params := []string{
		"-hide_banner",
		"-r", "24",
		"-rtsp_transport", "tcp",
		"-re",
		"-i", From, //"rtsp://192.168.1.210:554/av0_0"
		"-c:v", "mjpeg",
		"-f", "mjpeg",
		"-headers", "Content-Type:multipart/x-mixed-replace; boundary=MJPEG_BOUNDARY",
		ToUrl, // "http://127.0.0.1:9401/jpeg_stream/push?liveId=123",
	}
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		bat := strings.Join(params, " ")
		cmd = exec.Command("powershell.exe", "-Command", "ffmpeg "+bat)
	} else {
		cmd = exec.Command("ffmpeg", params...)
	}
	return cmd
}

/*
*
* 推送摄像头到Websocket Server
ffmpeg.exe -hide_banner -r 24 -rtsp_transport tcp -re
-i device=/dev/video0 -c:v mjpeg
-f mjpeg -headers "Content-Type:multipart/x-mixed-replace" "http://127.0.0.1:9401/jpeg_stream/push?liveId=123"
*/
func PushLocalCameraToWsServer(FromDevice, pushAddr string) *exec.Cmd {
	deviceName := ""
	if runtime.GOOS == "windows" {
		deviceName = fmt.Sprintf("video=\"%s\"", FromDevice)
	}
	if runtime.GOOS == "linux" {
		deviceName = fmt.Sprintf("video=%s", FromDevice)
	}
	var paramsVideo []string
	if runtime.GOOS == "windows" {
		paramsVideo = []string{
			"-err_detect", "ignore_err",
			"-hide_banner",
			"-f", "dshow", // windows下特有的DirectX加速引擎
			"-i", deviceName,
			"-q", "5",
			"-fflags", "+genpts",
			"-c:v", "libx264",
			"-preset", "veryfast",
			"-tune", "zerolatency",
			"-f", "mpegts",
			pushAddr,
		}
	} else {
		paramsVideo = []string{
			"-err_detect", "ignore_err",
			"-hide_banner",
			"-i", deviceName,
			"-q", "5",
			"-fflags", "+genpts",
			"-c:v", "libx264",
			"-preset", "veryfast",
			"-tune", "zerolatency",
			"-f", "mpegts",
			pushAddr,
		}
	}
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		bat := strings.Join(paramsVideo, " ")
		cmd = exec.Command("powershell.exe", "-Command", "ffmpeg "+bat)
	} else {
		cmd = exec.Command("ffmpeg", paramsVideo...)
	}
	return cmd
}

/*
*
* 推送RTSP到Websocket Server
ffmpeg.exe -hide_banner -r 24 -rtsp_transport tcp -re
-i rtsp://192.168.1.210:554/av0_0 -i .\meteorology.png -filter_complex "[0:v][1:v]overlay=10:10:alpha=0.5"
-q 5 -f mpegts -fflags nobuffer -c:v libx264 -an http://127.0.0.1:9400/h264_stream/push?liveId=123
*/
func PushRtspToWsServer(inputUrl, pushAddr string) *exec.Cmd {
	paramsRtsp := []string{
		"-hide_banner",
		"-r", "24",
		"-rtsp_transport", "tcp",
		"-re", "-i", inputUrl,
		"-q", "5",
		"-f", "mpegts",
		"-fflags", "nobuffer",
		"-c:v", "libx264",
		"-an",
		// http://127.0.0.1:9400/h264_stream/push?liveId=147a6d7ae5a785f6e3ea90f25d36c63e
		pushAddr,
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		bat := strings.Join(paramsRtsp, " ")
		cmd = exec.Command("powershell.exe", "-Command", "ffmpeg "+bat)
	} else {
		cmd = exec.Command("ffmpeg", paramsRtsp...)
	}
	return cmd
}

/*
*
* Push to rtmp
ffmpeg.exe -hide_banner -rtsp_transport tcp -re
-i rtsp://192.168.1.210:554/av0_0 -q 5 -f flv rtmp://127.0.0.1:1935
*/
func PushRtspToRTMPServer(inputUrl, pushAddr string) *exec.Cmd {
	paramsRtsp := []string{
		"-hide_banner",
		"-r", "24",
		"-rtsp_transport", "tcp",
		"-re", "-i", inputUrl,
		"-q", "5",
		"-f", "flv",
		"-fflags", "nobuffer",
		"-c:v", "libx264",
		"-an",
		pushAddr,
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		bat := strings.Join(paramsRtsp, " ")
		cmd = exec.Command("powershell.exe", "-Command", "ffmpeg "+bat)
	} else {
		cmd = exec.Command("ffmpeg", paramsRtsp...)
	}
	return cmd
}

/*
*
* 摄像头推流到RTMP服务器
*
 */
func PushLocalCameraToRTMPServer(FromDevice, pushAddr string) *exec.Cmd {
	deviceName := ""
	if runtime.GOOS == "windows" {
		deviceName = fmt.Sprintf("video=\"%s\"", FromDevice)
	}
	if runtime.GOOS == "linux" {
		deviceName = fmt.Sprintf("video=%s", FromDevice)
	}
	var paramsVideo []string
	if runtime.GOOS == "windows" {
		paramsVideo = []string{
			"-err_detect", "ignore_err",
			"-hide_banner",
			"-f", "dshow", // windows下特有的DirectX加速引擎
			"-i", deviceName,
			"-q", "5",
			"-fflags", "+genpts",
			"-c:v", "libx264",
			"-preset", "veryfast",
			"-tune", "zerolatency",
			"-f", "flv",
			pushAddr,
		}
	} else {
		paramsVideo = []string{
			"-err_detect", "ignore_err",
			"-hide_banner",
			"-i", deviceName,
			"-q", "5",
			"-fflags", "+genpts",
			"-c:v", "libx264",
			"-preset", "veryfast",
			"-tune", "zerolatency",
			"-f", "flv",
			pushAddr,
		}
	}
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		bat := strings.Join(paramsVideo, " ")
		cmd = exec.Command("powershell.exe", "-Command", "ffmpeg "+bat)
	} else {
		cmd = exec.Command("ffmpeg", paramsVideo...)
	}
	return cmd
}
