package sensor_server

import (
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
	Sn() string                 // 获取编号
	Session() Session           // 会话
	Ping() []byte               // PING包
	OnRegister(sn string) error // 注册成功
	OnLine()                    // 上线
	OffLine()                   // 掉线
	OnError(error)              // 出错
	OnData([]byte)              // 来数据
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
