package rulexlib

import (
	"errors"

	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rulex/component/aibase"
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
var errType error = errors.New("tensor type error, must be [][]float table")

func Infer(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		uuid := l.ToString(2)
		input := l.ToTable(3)
		inputTensor := [][]float64{}
		var err error = nil
		input.ForEach(func(i, v lua.LValue) {
			row := []float64{}
			switch vv := v.(type) {
			case *lua.LTable:
				{
					vv.ForEach(func(ir, vr lua.LValue) {
						switch vrv := vr.(type) {
						case lua.LNumber:
							{
								row = append(row, float64(vrv))
							}
						}
					})
				}
			default:
				{
					err = errType
				}
			}
			if err == nil {
				inputTensor = append(inputTensor, row)
			}
		})
		ai := aibase.AIBaseRuntime()
		if err != nil {
			l.Push(lua.LNil)                 // data
			l.Push(lua.LString(err.Error())) //err
			return 2
		}
		InferResult := [][]float64{}
		if ai != nil {
			if xai := aibase.GetAi(uuid); xai != nil {
				InferResult = xai.XAI.Infer(inputTensor)
			}
		}
		// {
		// 	{1,2,3,4,}
		// 	{1,2,3,4,}
		// }
		returnTable := lua.LTable{}
		for _, row := range InferResult {
			rowTable := lua.LTable{}
			for _, col := range row {
				rowTable.Append(lua.LNumber(col))
			}
			returnTable.Append(&rowTable)
		}
		l.Push(&returnTable) // data
		l.Push(lua.LNil)     //err
		return 2
	}
}
