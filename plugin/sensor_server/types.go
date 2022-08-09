package sensor_server

import (
	"context"
	"net"

	"github.com/google/uuid"
)

/*
*
* 服务器
*
 */
type ISServer interface {
	Start()
	Stop()
	AddSensor(Sensor)
	RemoveSensor(Sensor)
	Write(Sensor, []byte)
}

/*
*
* 传感器接口
*
 */
type ISensor interface {
	Addr() string
	Session() Session
	Ping() []byte
	OnRegister() error
	OnLine()
	OffLine()
	OnError(error)
	OnData([]byte)
}

/*
*
* 设备会话层
*
 */
type Session struct {
	Id        string
	Transport net.Conn
}

func NewSession(Transport net.Conn) Session {
	return Session{
		Id:        uuid.NewString(),
		Transport: Transport,
	}
}

/*
*
* 设备表示层、应用层
*
 */

type Sensor struct {
	session Session
	addr    string
	Authed  bool
}

func (s *Sensor) Addr() string {
	return s.addr
}
func (s *Sensor) Session() Session {

}
func (s *Sensor) Ping() []byte {

}
func (s *Sensor) OnRegister() error {

}
func (s *Sensor) OnLine() {

}
func (s *Sensor) OffLine() {

}
func (s *Sensor) OnError(error) {

}
func (s *Sensor) OnData([]byte) {

}

func NewSensor(addr string, session Session) Sensor {
	return Sensor{addr: addr, session: session}
}

/*
*
* 设备的工作进程
*
 */
type SensorWorker struct {
	Ctx    context.Context
	Cancel context.CancelFunc
	Sensor ISensor
}

func (w *SensorWorker) Run() {
	go func(ctx context.Context) {
		// ticker := time.NewTicker(5 * time.Second)
		// defer ticker.Stop()
		// buffer := make([]byte, common.T_64KB)
		// for {
		// 	<-ticker.C
		// 	select {
		// 	case <-ctx.Done():
		// 		{
		// 			return
		// 		}
		// 	default:
		// 		{
		// 		}
		// 	}
		// 	n, err := w.Sensor.Session().Transport.Read(buffer)
		// 	if err != nil {
		// 		log.Error(err)
		// 		w.Sensor.OnError(err)
		// 		w.Sensor.OffLine()
		// 		return
		// 	}
		// 	w.Sensor.OnData(buffer[:n])
		// }

	}(w.Ctx)
}
