---@diagnostic disable: undefined-global

-- Actions
Actions = { function(args)
    applib:WriteDevice('DEVICE71770e5db1f84bdfa6099cb3c7f6c48e',
        "cmd",
        rulexlib:T2J({ {
            SlaverId = 1,
            Function = 16,
            Address  = 0,
            Value    = { 11, 22, 32 },
        } }))
    return false, data
end }
