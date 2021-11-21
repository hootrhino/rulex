# 测试应用
## MQTT测试
```sh
#!/bin/bash
echo "Publish test."
for i in {1..1000}; do
    mosquitto_pub -h  127.0.0.1 -p 1883 -t '$X_IN_END' -q 2 -m "{\"temp\": $RANDOM,\"hum\":$RANDOM}"
    echo "Publish ", $i, " Ok."
done

```

## 串口模拟
### 安装 socat
```
sudo apt install socat
```
### 运行
```sh
socat -d -d -d  pty,raw,echo=0 pty,raw,echo=1
```
### 输出
```
2021/09/11 19:18:03 socat[6074] I socat by Gerhard Rieger and contributors - see www.dest-unreach.org
2021/09/11 19:18:03 socat[6074] I This product includes software developed by the OpenSSL Project for use in the OpenSSL Toolkit. (http://www.openssl.org/)
2021/09/11 19:18:03 socat[6074] I This product includes software written by Tim Hudson (tjh@cryptsoft.com)
2021/09/11 19:18:03 socat[6074] I setting option "raw"
2021/09/11 19:18:03 socat[6074] I setting option "echo" to 0
2021/09/11 19:18:03 socat[6074] I openpty({5}, {6}, {"/dev/pts/0"},,) -> 0
2021/09/11 19:18:03 socat[6074] N PTY is /dev/pts/0
2021/09/11 19:18:03 socat[6074] I setting option "raw"
2021/09/11 19:18:03 socat[6074] I setting option "echo" to 1
2021/09/11 19:18:03 socat[6074] I openpty({7}, {8}, {"/dev/pts/1"},,) -> 0
2021/09/11 19:18:03 socat[6074] N PTY is /dev/pts/1

```

## 关键输出
```
2021/09/11 19:18:03 socat[6074] N PTY is /dev/pts/0
2021/09/11 19:18:03 socat[6074] N PTY is /dev/pts/1
```
标识开启了两个模拟串口:

- `/dev/pts/0`
- `/dev/pts/1`

一个用来发送，一个用来接收。
