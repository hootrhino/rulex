package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/bluele/gcache"
)

func TestGoCache(t *testing.T) {

	gc := gcache.New(20).
		LRU().
		Build()
	gc.SetWithExpire("key", "ok", time.Second*10)
	gc.Set("h", "w")
	value, _ := gc.Get("key")
	fmt.Println("Get:", value)

	// Wait for value to expire
	time.Sleep(time.Second * 10)

	value, err := gc.Get("key")
	if err != nil {
		panic(err)
	}
	fmt.Println("Get:", value)
}
