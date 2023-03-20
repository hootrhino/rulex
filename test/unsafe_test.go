package test

import (
	"fmt"
	"testing"
	"unsafe"
)

type M struct {
	a1 uint
	a2 uint
	a3 uint
	a4 uint
	a5 uint
}

func (m *M) String() string {
	return fmt.Sprintf("a1=%v; a2=%v; a3=%v; a4=%v; a5=%v;", m.a1, m.a2, m.a3, m.a4, m.a5)
}
func Test_unsafe_pointer(t *testing.T) {

	var a1 uint = 1
	var a2 uint = 2
	var a3 uint = 3
	var a4 uint = 4
	var a5 uint = 5
	t.Log("real pointer of a1: ", unsafe.Pointer(&a1))
	t.Log("real pointer of a2: ", unsafe.Pointer(&a2))
	t.Log("real pointer of a3: ", unsafe.Pointer(&a3))
	t.Log("real pointer of a4: ", unsafe.Pointer(&a4))
	t.Log("real pointer of a5: ", unsafe.Pointer(&a5))
	t.Log("array: ", (*((*[5]int)(unsafe.Pointer(&a1)))))
	t.Log("m: ", (*((*M)(unsafe.Pointer(&a1)))).String())
	t.Logf("up1+0: 0x%x  %v", uintptr(unsafe.Pointer(&a1))+8*0, (*(*uint)(unsafe.Pointer(uintptr(unsafe.Pointer(&a1)) + unsafe.Sizeof(uint(1))*0))))
	t.Logf("up1+1: 0x%x  %v", uintptr(unsafe.Pointer(&a1))+8*1, (*(*uint)(unsafe.Pointer(uintptr(unsafe.Pointer(&a1)) + unsafe.Sizeof(uint(1))*1))))
	t.Logf("up1+2: 0x%x  %v", uintptr(unsafe.Pointer(&a1))+8*2, (*(*uint)(unsafe.Pointer(uintptr(unsafe.Pointer(&a1)) + unsafe.Sizeof(uint(1))*2))))
	t.Logf("up1+3: 0x%x  %v", uintptr(unsafe.Pointer(&a1))+8*3, (*(*uint)(unsafe.Pointer(uintptr(unsafe.Pointer(&a1)) + unsafe.Sizeof(uint(1))*3))))
	t.Logf("up1+4: 0x%x  %v", uintptr(unsafe.Pointer(&a1))+8*4, (*(*uint)(unsafe.Pointer(uintptr(unsafe.Pointer(&a1)) + unsafe.Sizeof(uint(1))*4))))

}
