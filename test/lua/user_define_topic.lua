Actions = {
    function(data)
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
        print(rulexlib:DataToMqtt('UUID', Json))
    end
}
