package rulexlib

import (
	lua "github.com/i4de/gopher-lua"
	"github.com/i4de/rulex/typex"
)

/*
*
* - AI 推理, 即AI的入口回调接口,入参必须是一个二维 [][]float 数组, 对应lua的Table
* - P = {
*       [0] = {11,12,13,14,15,16,17,18},
*       [1] = {21,22,23,24,25,26,27,28},
*       [2] = {31,32,33,34,35,36,37,38}
*   }
*   返回值也是一个二维矩阵, 代表预测的结果
*
 */
func Infer(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		uuid := l.ToString(2)
		// input := l.ToTable(3)
		ai := rx.GetAiBase()
		if ai != nil {
			if ai.GetAi(uuid) != nil {
				ai.Infer([][]float64{
					{10, 11, 12, 13, 14},
					{20, 21, 22, 23, 24},
				})
			}
		}
		// 简单示例
		result := lua.LTable{}
		row1 := &lua.LTable{}
		row1.Append(lua.LNumber(10))
		row1.Append(lua.LNumber(11))
		row1.Append(lua.LNumber(12))
		row1.Append(lua.LNumber(13))
		row1.Append(lua.LNumber(14))
		row2 := &lua.LTable{}
		row2.Append(lua.LNumber(20))
		row2.Append(lua.LNumber(21))
		row2.Append(lua.LNumber(22))
		row2.Append(lua.LNumber(23))
		row2.Append(lua.LNumber(24))
		result.Append(row1)
		result.Append(row2)
		l.Push(&result)  // data
		l.Push(lua.LNil) //err
		return 2
	}
}
