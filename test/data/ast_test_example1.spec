package rulexlib

//@desc:数据转发到HTTP服务器
func __RULEX_DataToHttp(
	uuid string, //@arg: HTTP UUID
	data string, //@arg: 数据
) error //@arg: 错误信息

//@desc:数据转发到TdEngine服务器
func __RULEX_DataToTdEngine(
	uuid string, //@arg: Tdengine UUID
	data string, //@arg: 数据
) error //@arg: 错误信息


//@desc:JSON解析
func __RULEX_parseJson(
	json string, //@arg: 数据
) (
	string, //@arg: 解码后的JSON
	error, //@arg: 错误信息
)
