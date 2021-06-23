## Bench test
Copy to EMQX console.
### 100000 Message
- 100000 message
- peer message 1 byte

```erlang
    Count = 100000,
    {ok, C} = emqtt:start_link([{host, "localhost"}, {clientid, <<"BenchX">>}]),
    {ok, _} = emqtt:connect(C),
    lists:foreach(fun(I) ->
        emqtt:publish(C, <<"test">>, erlang:integer_to_binary(I), qos0)
    end, lists:seq(1, Count)).
```