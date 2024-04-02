package rulexlib

import (
	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rulex/typex"
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
// var errType error = errors.New("tensor type error, must be [][]float table")

func Infer(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		return 0
	}
}
