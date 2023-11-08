AppNAME = "Test1"
AppVERSION = "1.0.0"
AppDESCRIPTION = ""
--
-- Main
--

function Main(arg)
    while true do
        for i = 1, 5, 1 do
            local result, err = applib:CtrlDevice('DEVICEda7ea0bdcf364ca7b7dda5e0cca647d7', "0" .. i)
            print("|*** CtrlDevice [0x01] result=>", result, err)
            time:Sleep(50)
        end
    end
    return 0
end
