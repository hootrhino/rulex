---@diagnostic disable: undefined-global
--
-- App use lua syntax, goto https://hootrhino.github.io for more document
-- APPID: APP6b28330ff4be4b0ba2f3e9317c4e2a47
--
AppNAME = "LED-RGB"
AppVERSION = "1.0.0"
AppDESCRIPTION = ""
--
-- Main
--
function Main(arg)
    while true do
        ws1608:GPIOSet("red", 1)
        time:Sleep(200)
        ws1608:GPIOSet("red", 0)
        time:Sleep(200)
        --
        ws1608:GPIOSet("green", 1)
        time:Sleep(200)
        ws1608:GPIOSet("green", 0)
        time:Sleep(200)
        --
        ws1608:GPIOSet("blue", 1)
        time:Sleep(200)
        ws1608:GPIOSet("blue", 0)
        time:Sleep(200)
    end
end
