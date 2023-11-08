AppNAME = "RMD_ISTMP10"
AppVERSION = "1.0.0"
AppDESCRIPTION = "RMD_ISTMP10 温度传感器"
--
-- Main
--

function Main(arg)
    while true do
        local result, err = applib:CtrlDevice('DEVICEb2a06c0130fd443e97d3980611aa3064', "010300010001D5CA")
        if err ~= nil then
            stdlib:Debug("!!!*** CtrlDevice Error=>" .. err)
        else
            stdlib:Debug("√√√*** CtrlDevice result=>" .. result)
        end
        time:Sleep(100)
    end
end
