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

//@namespace: rulexlib
//@desc: 数据转发到HTTP服务器
//@example: 示例代码
func __RULEXLIB_DataToHttp(
	// 1
	uuid string, //HTTP UUID
	// 2
	data string, //数据
) error

//@namespace: rulexlib
//@desc: 数据转发到TdEngine服务器
//@example: 示例代码
func __RULEXLIB_DataToTdEngine(
	// 1
	uuid string, /*Tdengine UUID*/
	// 2
	data string, /*数据*/
) error
