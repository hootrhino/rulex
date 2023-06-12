package sensor_server

import (
	"context"

	"time"

	"github.com/hootrhino/rulex/common"
	"github.com/hootrhino/rulex/glogger"
)

/*
*
* 设备的工作进程
*
 */
type SensorWorker struct {
	Ctx     context.Context
	Cancel  context.CancelFunc
	isensor ISensor
}

func (w *SensorWorker) Run() {
	go func(ctx context.Context) {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		buffer := make([]byte, common.T_64KB)
		for {
			select {
			case <-ctx.Done():
				{
					ticker.Stop()
					w.Cancel()
					return
				}
			default:
				{
				}
			}
			if len(w.isensor.Ping()) != 0 {
				_, err := w.isensor.Session().Transport.Write(w.isensor.Ping())
				if err != nil {
					glogger.GLogger.Error(err)
					w.isensor.OnError(err)
					w.isensor.OffLine()
					return
				}
			}
			// 5S读数据
			w.isensor.Session().Transport.SetDeadline(time.Now().Add(5 * time.Second))
			n, err := w.isensor.Session().Transport.Read(buffer)
			w.isensor.Session().Transport.SetDeadline(time.Time{})
			if err != nil {
				glogger.GLogger.Error(err)
				w.isensor.OnError(err)
				w.isensor.OffLine()
				return
			}
			w.isensor.OnData(buffer[:n])
			<-ticker.C
		}

	}(w.Ctx)
}
