package glogger

import (
	"os"
)

type LogMsg struct {
	Sn      string `json:"sn"`
	Uid     string `json:"uid"`
	Content string `json:"content"`
}

/*
*
* 日志记录器，未来会移除这个slot
*
 */
type LogWriter struct {
	file *os.File
}
