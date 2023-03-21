package test

import (
	"testing"

	"github.com/emirpasic/gods/maps/linkedhashmap"
)

func Test_Name(t *testing.T) {
	m := linkedhashmap.New()
	m.Put(1, []int{100, 200, 300})
	m.Put(2, "2")
	m.Put(3, "3")
	m.Put(4, "4")
	m.Put(5, "5")
	m.Put(6, "6")
	b, err := m.ToJSON()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))
	_, _ = m.Get(2) // b, true
	_, _ = m.Get(3) // nil, false
	_ = m.Values()  // []interface {}{"b", "a"} (insertion-order)
	_ = m.Keys()    // []interface {}{2, 1} (insertion-order)
	m.Remove(1)     // 2->b
	m.Clear()       // empty
	m.Empty()       // true
	m.Size()        // 0
}
