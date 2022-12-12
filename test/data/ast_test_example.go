package test

// This is doc1
func T_F1() {

}

/*
 @overview: This is doc2
 @arg1: 1
*/
func T_F2() {

}

// @overview: 尝试把Go翻译成LUA
// function add(a,b)
//     return a+b
// end
//
//
func GoToLuaDemo(a, b int) int {
	c := a + b
	return c
}
