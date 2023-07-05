package test

// import (
// 	"fmt"
// 	"golang.org/x/sys/unix"
// 	"testing"
// )

// func Test_USB_PLUG(t *testing.T) {

// 	fd, err := unix.Socket(
// 		unix.AF_NETLINK,
// 		unix.SOCK_RAW,
// 		unix.NETLINK_KOBJECT_UEVENT,
// 	)

// 	if err != nil {
// 		fmt.Println(err. common.Error())
// 	}

// 	err = unix.Bind(fd, &unix.SockaddrNetlink{
// 		Family: unix.AF_NETLINK,
// 		Groups: 1,
// 		Pid:    0,
// 	})

// 	if err == nil {
// 		for {
// 			data := make([]byte, 1024)
// 			n, _, _ := unix.Recvfrom(fd, data, 0)
// 			if n != 0 {
// 				fmt.Println(string(data[:n]))
// 			}
// 		}
// 	} else {
// 		fmt.Println(err. common.Error())
// 	}
// }
