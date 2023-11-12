Actions = {
    function(args)
        local Json = rulexlib:T2J(
            {
                method = 'userDefineTopic',
                userDefineTopic = '/userDefineTopic',
                params = {
                    tag = "device",
                    temp = 12.34,
                    hum = 56.78,
                }
            }
        )
        print("Json ->:", Json)
        print(data:ToMqtt('UUID', Json))
    end
}
