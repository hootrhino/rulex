package rulexlib

/*
*
* 这里规定一些LUA和golang的类型映射
*  local t = {
*      ["type"] = 5,
*      ["params"] = {
*          ["address"] = 1,
*          ["quantity"] = 1,
*          ["value"] = 0xFF00
*      }
*  }
*
 */
type ModbusW struct {
	SlaverId byte   // 从机ID
	Function int    // 功能码
	Address  uint16 // 地址
	Quantity uint16 // 读写数量
	Value    []byte // 值
}
