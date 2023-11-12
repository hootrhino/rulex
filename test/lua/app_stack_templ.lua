-- 应用名称
AppNAME = "test_demo"
-- 版本信息,必须符合 x.y.z 格式
AppVERSION = "1.0.0"
-- 关于这个应用的一些描述信息
AppDESCTIPTION = "A demo app"
-- 必须包含 APPMain(arg) 函数作为 app 启动点
function Main(arg)
	while true do
		local value, err = rhinopi:GPIOGet(6)
		if err ~= nil
		then
			print(err)
		else
			print("value ok:", value)
		end
	end
end
