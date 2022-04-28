-- Success
function Success()
end
-- Failed
function Failed(error)
    print("Error:", error)
end
-- {
--     "cmd":[
--         0,
--         1,
--         0,
--         1,
--         0,
--         0,
--         1,
--         1
--     ]
-- }
Actions = {function(data)
    local tb = rulexlib:J2T(data)
    local cmd = tb['cmd']
    for _, value in ipairs(cmd) do
        rulexlib:WriteOutStream('UUID', rulexlib:T2J({
            type = 5,
            params = {
                address = 1,
                quantity = 1,
                value = value
            }
        }))
    end
    return true, data
end}
