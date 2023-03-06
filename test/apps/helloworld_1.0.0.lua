AppNAME = "test_demo"
AppVERSION = "1.0.0"
AppDESCRIPTION = "A demo app"

function Main(arg)
	print(arg)
	for key, value in pairs(_G) do
		print("_G:", key, value)
	end
	print("Message From Main:", AppNAME, AppVERSION, AppDESCRIPTION, rulexlib:Time())
	return 0
end
