package glogger

import (
	"net"
)

/*
*
* 日志记录器, 使用UDP协议将日志打到云端
*
 */
type UdpLogger struct {
}

func NewUdpLogger(filepath string, maxSlotCount int) *UdpLogger {

	return &UdpLogger{}
}

//
func (lw *UdpLogger) Write(b []byte) (n int, err error) {
	socket, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 30000,
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
