package archsupport

//#include <stdio.h>
import "C"

/*
*
* 这是Linux库函数调用，本案例没有意义，主要用来展示规范化
*
 */
func HelloWorld() {
	C.puts(C.CString("Hello, World\n"))
}
