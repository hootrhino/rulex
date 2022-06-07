package core

import (
	"rulex/typex"
	"time"

	"github.com/ngaut/log"
)

func StartLogWatcher(path string) {
	typex.GLOBAL_LOGGER = typex.NewLogWriter("./"+time.Now().Format("2006-01-02_15-04-05-")+path, 1000)
	log.SetRotateByDay()
	log.SetOutput(typex.GLOBAL_LOGGER)
}
