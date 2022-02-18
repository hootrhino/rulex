package core

import (
	"rulex/typex"
	"time"

	"github.com/ngaut/log"
)

var GLOBAL_LOGGER *typex.LogWriter

func StartLogWatcher(path string) {
	GLOBAL_LOGGER = typex.NewLogWriter("./"+time.Now().Format("2006-01-02_15-04-05-")+path, 1000)
	log.SetRotateByDay()
	log.SetOutput(GLOBAL_LOGGER)
}
