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
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package pusher

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// StreamPusher 结构体包含了推流器的相关信息
type StreamPusher struct {
	Cmd   *exec.Cmd
	Stdin io.WriteCloser
}

// NewStreamPusher 创建一个新的推流器实例
func NewStreamPusher(rtmpURL string) (*StreamPusher, error) {
	// 构建 FFmpeg 推流命令
	cmd := exec.Command(
		"./ffmpeg.exe",
		"-re",
		"-f", "image2pipe",
		"-vcodec", "png",
		"-i", "-",
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"-preset", "fast",
		"-r", "25",
		"-f", "flv",
		rtmpURL,
	)
	fmt.Println("NewStreamPusher", cmd.String())
	// 获取 FFmpeg 的标准输入
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("error creating stdin pipe: %v", err)
	}
	cmd.Stdout = os.Stdout
	return &StreamPusher{
		Cmd:   cmd,
		Stdin: stdin,
	}, nil
}
func (p *StreamPusher) WritePNG(pngData []byte) error {
	if _, err := p.Stdin.Write(pngData); err != nil {
		return fmt.Errorf("error writing PNG data to stdin: %v", err)
	}
	return nil
}

func (p *StreamPusher) StartPush() error {
	defer p.Close()
	if err := p.Cmd.Run(); err != nil {
		return fmt.Errorf("error starting FFmpeg command: %v", err)
	}
	return nil
}

// Close 关闭推流器，确保所有资源被释放
func (p *StreamPusher) Close() error {

	// 关闭标准输入，以确保 FFmpeg 正常结束
	if err := p.Stdin.Close(); err != nil {
		return fmt.Errorf("error closing stdin pipe: %v", err)
	}

	return nil
}

func test() {
	// 假设 RTMP 服务器的 URL 为 "rtmp://example.com/live/stream"
	rtmpURL := "rtmp://example.com/live/stream"

	// 创建一个新的推流器
	pusher, err := NewStreamPusher(rtmpURL)
	if err != nil {
		fmt.Printf("Error creating stream pusher: %v\n", err)
		return
	}
	defer pusher.Close()

	// 向推流器的标准输入写入 H264 流数据
	// 这里只是一个示例，实际应用中应该从视频源获取 H264 数据
	// 比如从摄像头、视频文件或网络流中读取
	h264Data := []byte{} // 假设的 H264 数据
	if _, err := pusher.Stdin.Write(h264Data); err != nil {
		fmt.Printf("Error writing to stream pusher: %v\n", err)
		return
	}
}
