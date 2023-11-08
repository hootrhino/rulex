AppNAME = "test_demo"
AppVERSION = "1.0.0"
AppDESCRIPTION = "A demo app"

function Main(arg)
	for i = 1, 10, 1 do
		print("Hello App:", AppNAME, AppVERSION, AppDESCRIPTION, applib:Time())
		time:Sleep(1000)
	end
	return 0
end
