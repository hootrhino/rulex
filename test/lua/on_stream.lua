---@diagnostic disable: undefined-global

-- Success
function Success()
    print("======> success")
end
-- Failed
function Failed(error)
    print("======> failed:", error)
end

-- Actions
Actions = {
    function(data)
        stdlib:WriteInStream('INEND_e1a3f6df-1048-4582-aeac-9e7c45352db1', data)
        stdlib:WriteOutStream('OUTEND_e1a3f6df-1048-4582-aeac-9e7c45352db1', data)
        return true, data
    end
}
