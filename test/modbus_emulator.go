package test

import (
	"log"
	"time"

	"github.com/tbrandon/mbserver"
)

func _T() {
	server := mbserver.NewServer()
	err := server.ListenTCP("127.0.0.1:502")
	if err != nil {
		log.Printf("%v\n", err)
		return
	}
	defer server.Close()

	for {
		time.Sleep(1 * time.Second)
		server.Coils = []byte{1, 2, 3, 4, 5, 6, 7, 8}
	}
}
