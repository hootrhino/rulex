## 注意

## SNMP
github.com/gosnmp/gosnmp：这个库有点问题，效率不高，因为大量用了循环导致监控点数量多了以后，就会大量吃CPU,因此这个SNMP的监控[不推荐正式场景]使用，后期可能会找个好点的库重写这个功能.
## 资源状态同步
下面接口很重要，需要非常准确的返回资源的状态，才能触发重启。
```
Status() typex.ResourceState {}
```

## 测试MQTT
```
#!/bin/bash
echo "Publish test."
for i in {1..1000}; do
    mosquitto_pub -h  127.0.0.1 -p 1883 -t '$X_IN_END' -q 2 -m "{\"temp\": $RANDOM,\"hum\":$RANDOM}"
    echo "Publish ", $i, " Ok."
done

```