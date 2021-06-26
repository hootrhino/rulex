# RulEngine X
RulEngine X 是一个轻量级规则引擎(名字看起来有点怪实际上是 rule + Engine 的组合词，中文发音为"若金克斯")。主要用来中转上游数据和吐出数据到目标点，可以理解为一个数据路由器。主要被设计用来做物联网网关或者服务端存在数据筛选的场景。
## 功能
- MQTT 数据输入
- HTTP 数据输入
- COAP 数据输入
- UDP 数据输入
- LUA 自定义业务逻辑支持
- SQL 字段筛选支持

## API接口
- HTTP RestFul API
## 用户管理
- Web dashboard

## Bench test
Copy to EMQX console.
- 100000 message
- peer message 1 byte

```erlang
    Count = 100000,
    {ok, C} = emqtt:start_link([{host, "localhost"}, {clientid, <<"BenchX">>}]),
    {ok, _} = emqtt:connect(C),
    lists:foreach(fun(I) ->
        io:format("Send:~p\n",[I]),
        emqtt:publish(C, <<"$X_IN_END">>, erlang:integer_to_binary(I), qos0)
    end, lists:seq(1, Count)).
```