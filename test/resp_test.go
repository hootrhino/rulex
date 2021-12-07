package test

import (
	"fmt"
	"net"
	"os"
	"testing"
)

func Test_set(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer conn.Close()
	set(conn)
	get(conn)
}
func set(conn net.Conn) {
	_, e := conn.Write([]byte("*3\r\n$3\r\nSET\r\n$5\r\nkkkkk\r\n$7\r\nvvvvvvv\r\n"))
	if e != nil {
		fmt.Println("Error to send message because of ", e.Error())
	}
	buf := make([]byte, 5)
	reqLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error to read message because of ", err)
		return
	}
	fmt.Println("set response:", string(buf[:reqLen-1]))
}
func get(conn net.Conn) {
	_, e := conn.Write([]byte("GET kkkkk\r\n"))
	if e != nil {
		fmt.Println("Error to send message because of ", e.Error())
	}

	buf := make([]byte, 512)
	reqLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error to read message because of ", err)
		return
	}
	fmt.Println("get response:", string(buf[:reqLen-1]))
}
