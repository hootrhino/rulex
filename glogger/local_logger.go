package glogger

import (
	"os"
)

/*
*
* 日志记录器，未来会移除这个slot
*
 */
type LogWriter struct {
	file *os.File
}

func NewLogWriter(filepath string, maxSlotCount int) *LogWriter {
	logFile, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		GLogger.Fatalf("Fail to read log file: %v", err)
		os.Exit(1)
	}

	return &LogWriter{file: logFile}
}
func (lw *LogWriter) Write(b []byte) (n int, err error) {
	return lw.file.Write(b)
}

func (lw *LogWriter) Close() error {
	if lw.file != nil {
		return lw.file.Close()
	} else {
		return nil
	}

}
