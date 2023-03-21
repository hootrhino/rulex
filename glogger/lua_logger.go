package glogger

import "time"

var private_lua_logger *LogWriter

/*
*
* StartLuaLogger
*
 */
func StartLuaLogger(path string) {
	private_lua_logger = NewLogWriter("./" + time.Now().Format("2006-01-02-") + path)
}

/*
*
* LUA 脚本的日志接口
*
 */
func LuaLog(b []byte) {
	private_lua_logger.Write(b)
}
