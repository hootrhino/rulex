package typex

import (
	"log"
	"os"
)

type LogWriter struct {
	file         *os.File
	logSlot      []string
	maxSlotCount int
}

func NewLogWriter(filepath string, maxSlotCount int) *LogWriter {
	logFile, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Fatalf("Fail to read log file: %v", err)
		os.Exit(1)
	}

	return &LogWriter{file: logFile,
		logSlot: make([]string, maxSlotCount),
	}
}
func (lw *LogWriter) Write(b []byte) (n int, err error) {
	if len(lw.logSlot) > lw.maxSlotCount {
		lw.logSlot = append(lw.logSlot[1:], string(b))
	} else {
		lw.logSlot = append(lw.logSlot, string(b))
	}

	return lw.file.Write(b)
}

func (lw *LogWriter) Slot() []string {
	return lw.logSlot
}
func (lw *LogWriter) Close() error {
	if lw.file != nil {
		return lw.file.Close()
	} else {
		return nil
	}

}
