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

package aibase

/*
*
* 算法模型
*
 */
type Algorithm struct {
	UUID        string // UUID
	Type        string // 模型类型: ANN_APP1 RNN_APP2 CNN_APP3 ....
	Name        string // 名称
	State       int    // 0开启;1关闭
	Document    string // 文档连接
	Description string // 概述
}

/*
*
* AI 接口
*
 */
type AlgorithmResource interface {
	Init(map[string]interface{}) error // 初始化环境
	// Type , Sample, ExpectOut
	Train(string, [][]float64, [][]float64) error      // 训练模型
	Load() error                                       // 加载模型
	OnCall(string, [][]float64) map[string]interface{} // 用数据去执行
	Unload() error                                     // 卸载模型
	AiDetail() Algorithm                               // 获取信息
}
