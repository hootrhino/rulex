package rulexlib

import (
	"encoding/json"
	"rulex/typex"

	"github.com/ngaut/log"
	lua "github.com/yuin/gopher-lua"
)

/*
*
* 指令执行成功
*
 */
type CmdSuccessLib struct {
}

func NewCmdSuccessLib() typex.XLib {
	return &LogLib{}
}
func (l *CmdSuccessLib) Name() string {
	return "finishCmd"
}
func (l *CmdSuccessLib) LibFun(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		cmdId := l.ToString(2)
		stateTargetId := l.ToString(3)
		bytes, _ := json.Marshal(map[string]interface{}{
			"type":  "finishCmd",
			"cmdId": cmdId,
		})
		write(rx, stateTargetId, string(bytes))
		return 0
	}
}

/*
*
* 指令执行失败
*
 */
type CmdFailedLib struct {
}

func NewCmdFailedLib() typex.XLib {
	return &LogLib{}
}
func (l *CmdFailedLib) Name() string {
	return "failedCmd"
}
func (l *CmdFailedLib) LibFun(rx typex.RuleX) func(*lua.LState) int {
	return func(l *lua.LState) int {
		cmdId := l.ToString(2)
		stateTargetId := l.ToString(3)
		bytes, _ := json.Marshal(map[string]interface{}{
			"type":  "failedCmd",
			"cmdId": cmdId,
		})
		write(rx, stateTargetId, string(bytes))
		return 0
	}
}
func write(e typex.RuleX, uuid string, incoming string) {
	outEnd, exists := e.AllOutEnd().Load(uuid)
	if exists {
		(outEnd.(*typex.OutEnd)).Target.OnStreamApproached(incoming)
	} else {
		log.Error("OutEnd: " + uuid + " not exists")
	}
}
