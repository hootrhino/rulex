---@diagnostic disable: undefined-global
-- {
--     "d1":{
--         "tag":"d1",
--         "function":3,
--         "slaverId":1,
--         "address":1,
--         "quantity":3,
--         "value":"000100010000"
--     }
-- }
Actions = {
    function(args)
        local dataT, err0 = rulexlib:J2T(data)
        if err0 ~= nil then
            print("ERROR:", err0)
            goto END
        end
        for _, entity in ipairs(dataT) do
            print("tag", entity.tag)
            print("value", entity.value)
        end
        ::END::
        return true, args
    end
}