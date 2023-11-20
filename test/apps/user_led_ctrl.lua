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
        rhinopi:Led1On()
        time:Sleep(200)
        rhinopi:Led1Off()
        time:Sleep(200)
    end
end
