Actions = {
    function(data)
        print('Data From Ithings:', data)
        local Json = rulexlib:T2J(
            {
                method = 'report',
                params = {
                    tag = 'key',
                    temp = 0.1,
                    hum = 0.1,
                }
            }
        )
        print('Data to Ithings:', Json)
        rulexlib:DataToIthings('INe1e769cff3b9467394564ca78c7bc93b', Json)
        return true, data
    end
}
