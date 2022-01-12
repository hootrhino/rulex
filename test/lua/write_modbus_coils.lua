---@diagnostic disable: undefined-global
-- Success
function Success()
    --
end
-- Failed
function Failed(error)
    print(error)
end
-- Actions
Actions = {function(data)
    local t = {
        ["type"] = 15,
        ["params"] = {
            ["address"] = 3,
            ["quantity"] = 4,
            -- 写多个线圈，因为modbus每个寄存器的大小是2字节，因此下面尝试写了2个寄存器是4字节
            ["values"] = {0xFF, 0x00, 0xFF, 0x00}
        }
    }
    rulexlib:WriteInStream('INEND_77c16142-f849-48c3-b150-34aed2d0d9ae', rulexlib:JsonEncode(t))
    return false, data
end}
