package glogger

import (
	"encoding/json"
	"os"
)

func NewLogWriter(filepath string) *LogWriter {
	logFile, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		GLogger.Fatalf("Fail to read log file: %v", err)
		os.Exit(1)
	}

	return &LogWriter{file: logFile}
}

/*
*
* 日志记录的本地的同时,可能会记录到远程UDP Server, 该功能主要用来远程诊断
*
 */
func (lw *LogWriter) Write(b []byte) (n int, err error) {
	// UDP 日志及录默认不开启, 当端口配置成非0才会开
	if private_remote_logger.udpPort != 0 {
		logmsg := LogMsg{
			Sn:      private_remote_logger.Sn,
			Uid:     private_remote_logger.Uid,
			Content: string(b),
		}
		bytes, _ := json.Marshal(logmsg)
		private_remote_logger.Write(bytes)
	}
	return lw.file.Write(b)
}

func (lw *LogWriter) Close() error {
	if lw.file != nil {
		return lw.file.Close()
	} else {
		return nil
	}

}
