---@diagnostic disable: undefined-global

-- Success
function Success()

end

-- Failed
function Failed(error)
    print("failed:", error)
end

-- Actions

Actions = {
    function(data)
        local dataTable, err1 = rulexlib:J2T(data)
        if err1 ~= nil then
            return true, data
        end
        for _k, entity in pairs(dataTable) do
            rulexlib:DataToMongo("OUT8404b5afb7fe4baea335ebcb0d821491", rulexlib:T2J(entity["value"]))
        end
        return true, data
    end
}
