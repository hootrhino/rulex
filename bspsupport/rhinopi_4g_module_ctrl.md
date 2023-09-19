# RhinoPi 的4G模组操作接口
## 简介
RhinoPi自带一个4G上网模组，可以实现接入移动网络远程上网。该控制接口只可兼容上海移远公司的下列模组：
- EC200T 系列
- EC200S 系列
- EC200A 系列
- EC200N-CN
- EC600S-CN
- EC600N-CN
- EC800N-CN
- EG912Y-EU
- EG915N-EU

RhinoPi使用的是 EC200A 系列，更多模块资料请参考此处：https://www.quectel.com/cn/product/ec200a-series 。
## 操作
在shell输入ifconfig即可看到 usb0 网卡，该网卡就是对应的4G模组：
```sh
usb0: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500
        inet 10.175.11.79  netmask 255.0.0.0  broadcast 10.255.255.255
        inet6 fe80::d6eb:f5eb:da12:d4ab  prefixlen 64  scopeid 0x20<link>
        ether 02:0c:29:a3:9b:6d  txqueuelen 1000  (Ethernet)
        RX packets 319  bytes 20479 (20.4 KB)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 489  bytes 52138 (52.1 KB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0
```
## 扩展知识
### 4G模块的信号强度原理
4G模块的CSQ（Cellular Signal Quality）信号通常是用来表示移动网络信号强度的一种方式。CSQ值通常以分数或分贝（dBm）为单位。它是一个0到31之间的整数，其中0表示没有信号，31表示最强的信号。

具体的CSQ值与不同的硬件、模块和供应商有关，因此CSQ值的解释可能因设备而异。一般来说，以下是一些通用的CSQ值的解释：

- 0：没有信号。
- 1-9：非常弱的信号，可能无法建立连接。
- 10-14：较弱的信号，但可能可以建立连接。
- 15-19：中等强度的信号。
- 20-31：非常强的信号，信号质量非常好。

有些设备还可以提供更详细的信号质量信息，例如dBm值。dBm（分贝毫瓦）是一种衡量信号强度的标准单位，通常负数值表示较弱的信号，例如：

- -50 dBm：非常强的信号。
- -70 dBm：较强的信号。
- -90 dBm：较弱的信号。
- -110 dBm：非常弱的信号。

要获取CSQ或dBm值，你通常需要查询4G模块的AT命令或使用相关API。不同的模块和供应商提供的AT命令和API可能不同，因此你需要查阅模块的文档以了解如何获取信号质量信息。

一般来说，使用AT命令可能是获取CSQ或dBm值的常见方法。例如，你可以发送类似于以下的AT命令来获取CSQ值：

```sh
AT+CSQ\r\n
```

模块将返回CSQ值的响应。再次强调，确保查阅模块的文档以了解如何正确地获取信号质量信息，并将其解释为信号质量的度量标准。