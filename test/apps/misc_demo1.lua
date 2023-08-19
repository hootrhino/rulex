AppNAME = 'Data'
AppVERSION = '0.0.1'
function Main(arg)
    print("XOR: 0101=>", misc:XOR("0101", 0))
    print("XOR: 01010001=>", misc:XOR("01010001", 1))
    return 0
end
