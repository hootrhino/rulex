package test

import (
	"fmt"
	"net"
	"testing"
)

func Test_get_hostname(t *testing.T) {
	//ss, err := GetLocalIP()
	// if err != nil {
	// 	t. common.Error(err)
	// }
	//t.Log(ss)
	localAddresses()
}
func localAddresses() {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Print(fmt.Errorf("localAddresses: %+v\n", err.Error()))
		return
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Print(fmt.Errorf("localAddresses: %+v\n", err.Error()))
			continue
		}
		for _, a := range addrs {

			fmt.Println(a.String())
		}
	}
}

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() ([]string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return []string{}, err
	}
	ss := make([]string, len(addrs)-1)
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				// ip+net 172.30.128.1/20 172.30.128.1
				ss = append(ss, ipnet.IP.String())
			}
		}
	}
	return ss, nil
}
