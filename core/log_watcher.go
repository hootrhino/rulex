package core

import (
	"io"
	"os"

	"github.com/ngaut/log"
)

//
// 默认日志槽大小: 1000条
//
const max_LOG_COUNT int = 1000
var LogSlot []string = make([]string, max_LOG_COUNT)
type LogWriter struct {
	file *os.File
}

func NewLogWriter(filepath string) io.Writer {
	logFile, err2 := os.OpenFile(GlobalConfig.LogPath, os.O_WRONLY|os.O_CREATE, 0777)
	if err2 != nil {
		log.Fatalf("Fail to read log file: %v", err2)
		os.Exit(1)
	}
	return &LogWriter{file: logFile}
}
func (lw *LogWriter) Write(b []byte) (n int, err error) {
	if len(LogSlot) > max_LOG_COUNT {
		LogSlot = append(LogSlot[1:], string(b))
	} else {
		LogSlot = append(LogSlot, string(b))
	}
	return lw.file.Write(b)
}
func (lw *LogWriter) Close(p []byte) (err string) {
	return lw.file.Close().Error()
}

func StartLogWatcher() {
	log.SetOutput(NewLogWriter(GlobalConfig.LogPath))
	log.SetRotateByDay()
}
