package archsupport

import (
	"syscall"
	"unsafe"
)

var (
	user32, _          = syscall.LoadLibrary("user32.dll")
	procMessageBoxW, _ = syscall.GetProcAddress(user32, "MessageBoxW")
)

func HelloWorld() {
	defer syscall.FreeLibrary(user32)
	win32HelloWorld(0, "Hello, Win32 API(Go) World!", "Hello, World!", 0)
}

/*
*
* 这是win32系统调用，会弹出一个框显示一些文本，本案例没有意义，主要用来展示规范化
*
 */
func win32HelloWorld(hwnd uintptr, text string, caption string, style uintptr) int32 {
	ret, _, _ := syscall.Syscall6(
		uintptr(procMessageBoxW),
		4,
		hwnd,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(text))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(caption))),
		style,
		0,
		0)
	return int32(ret)
}
