AppNAME = "test_demo"
AppVERSION = "1.0.0"
AppDESCRIPTION = "A demo app"

function Main(arg)
	print("Message From Main:", AppNAME, AppVERSION, AppDESCRIPTION, applib:Time())
	return 0
end
