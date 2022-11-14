package glogger

import (
	"net"

	"github.com/sirupsen/logrus"
)

var private_remote_logger *UdpLogger
var RemoteLogger *logrus.Logger = logrus.New()

/*
*
* 日志记录器, 使用UDP协议将日志打到云端
*
 */
type UdpLogger struct {
	udpServer string
	udpPort   int
	Sn        string `json:"sn"`
	Uid       string `json:"uid"`
}

func NewUdpLogger(sn, uid, udpServer string, udpPort int) *UdpLogger {
	return &UdpLogger{Sn: sn, Uid: uid, udpServer: udpServer, udpPort: udpPort}
}

//
func (lw *UdpLogger) Write(b []byte) (n int, err error) {
	socket, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.ParseIP(lw.udpServer),
		Port: lw.udpPort,
	})
	if err != nil {
		return 0, nil
	}
	defer socket.Close()
	socket.Write(b)
	return 0, nil
}

func (lw *UdpLogger) Close() error {
	return nil
}

func StartRemoteLogger(sn, uid, udpServer string, udpPort int) {
	private_remote_logger = NewUdpLogger(sn, uid, udpServer, udpPort)
	RemoteLogger.Formatter = new(logrus.JSONFormatter)
	RemoteLogger.SetReportCaller(true)
	RemoteLogger.Formatter.(*logrus.JSONFormatter).PrettyPrint = false
	RemoteLogger.SetOutput(private_remote_logger)
}
