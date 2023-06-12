function Main(Arg)
    -- 读取10组刀具的平均厚度
    -- {1,2,3,4,5,6,7,8,9,10}
    local CutterData, err = applib:ReadDevice('Dev1', 'D1', "count=10")
    if err ~= nil then
        error(err)
        return
    end
    -- 交给 ID为'AI-001'的AI模型去计算结果
    -- 输出结果是一个数组，维度取决于模型输出参数
    -- Result: {1}
    local Result, err = aibase:Infer('AI-001', CutterData)
    if err ~= nil then
        error(err)
        return
    end
    print('Result =>', Result)
    return 0
end
