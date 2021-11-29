package rulexlib

import (
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
		// TODO 明天搞
		log.Info("[CmdSuccessLib ::: finishCmd, will write to emqx for sync]" + cmdId + " ==> " + stateTargetId)
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
		// TODO 明天搞
		log.Info("[CmdSuccessLib ::: finishCmd, will write to emqx for sync]" + cmdId + " ==> " + stateTargetId)
		return 0
	}
}
