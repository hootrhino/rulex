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

package archsupport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// nvidia-smi --query-gpu="index,name,temperature.gpu,memory.total,memory.used,memory.free,utilization.gpu,utilization.memory" --format="csv,noheader,nounits"
// GPUInfo 表示单个GPU的信息
type GPUInfo struct {
	Index    string `json:"index"`
	Name     string `json:"name"`
	GPUTemp  int    `json:"gpu_temp"`
	MemTotal int64  `json:"mem_total"`
	MemUsed  int64  `json:"mem_used"`
	MemFree  int64  `json:"mem_free"`
	GPUUtil  int    `json:"gpu_util"`
	MemUtil  int    `json:"mem_util"`
}

func (O GPUInfo) String() string {
	if bytes, err := json.Marshal(O); err != nil {
		return ""
	} else {
		return string(bytes)
	}
}

// GetGpuInfoWithNvidiaSmi 执行nvidia-smi命令并返回GPU信息
func GetGpuInfoWithNvidiaSmi() ([]GPUInfo, error) {
	// 创建一个*exec.Cmd对象，并设置nvidia-smi命令
	cmd := exec.Command("nvidia-smi",
		"--query-gpu=index,name,temperature.gpu,memory.total,memory.used,memory.free,utilization.gpu,utilization.memory",
		"--format=csv,noheader,nounits")

	// 创建一个缓冲区来捕获命令的输出
	var out bytes.Buffer
	cmd.Stdout = &out
	// 执行命令
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("error running nvidia-smi: %w, %s", err, out.String())
	}
	// 解析输出
	gpus := strings.Split(out.String(), "\n")
	var result []GPUInfo
	for _, gpu := range gpus {
		if gpu != "" {
			// 使用strings.Split将每行拆分为字段
			fields := strings.Split(strings.TrimSpace(gpu), ",")
			if len(fields) != 8 {
				continue // 跳过不符合格式的行
			}
			gpuInfo := GPUInfo{
				Index:    fields[0],
				Name:     fields[1],
				GPUTemp:  atoi(strings.TrimSpace(fields[2])),
				MemTotal: atoi64(strings.TrimSpace(fields[3])),
				MemUsed:  atoi64(strings.TrimSpace(fields[4])),
				MemFree:  atoi64(strings.TrimSpace(fields[5])),
				GPUUtil:  atoi(strings.TrimSpace(fields[6])),
				MemUtil:  atoi(strings.TrimSpace(fields[7])),
			}
			result = append(result, gpuInfo)
		}
	}
	return result, nil
}

// atoi 将字符串转换为整数
func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

// atoi64 将字符串转换为int64
func atoi64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return i
}
