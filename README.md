# RULEX Framework

### RULEX 是一个轻量级边缘网关开发框架，助力快速实现边缘数据处理和云边协同方案。

## 架构设计

<div style="text-align:center">
<img src="./README_RES/structure.png"/>
</div>

## 预览

![image](https://user-images.githubusercontent.com/20577297/249867828-afb6c81f-288e-47f9-b7d2-73330896ac30.png)
![image](https://user-images.githubusercontent.com/20577297/249867911-907827d1-5f1d-4ddb-bab7-3fc792f28c41.png)
![image](https://user-images.githubusercontent.com/20577297/249867961-9ca5c333-28d0-4154-9758-297e0bac3ca3.png)
![image](https://user-images.githubusercontent.com/20577297/249868010-8a5f1ca7-0203-4754-a206-cda48d75e331.png)
![image](https://user-images.githubusercontent.com/20577297/249868079-50ec6002-7447-4eca-9ebd-e32cd4d6caff.png)
![image](https://user-images.githubusercontent.com/20577297/249868117-288fffa0-7b96-4f82-85e1-97470f4dce35.png)
![image](https://user-images.githubusercontent.com/20577297/249868160-9662b07c-d189-4cea-be63-94b210abf908.png)


## 快速开始
### 源码编译
#### 环境安装
下面是Ubuntu上搭建环境的指令：
```bash
sudo apt install jq cloc protobuf-compiler \
     gcc-mingw-w64-x86-64 \
     gcc-arm-linux-gnueabi \
     gcc-mips-linux-gnu \
     gcc-mingw-w64 \
     gcc-aarch64-linux-gnu -y
```
> [!TIP]
> 推荐使用 ubuntu18.04 开发。

#### 编译
```sh
make
```

### HelloWorld
下面展示一个最简单的设备数据转发案例,这个 DEMO 展示了如何把一个简单的数据推到 UDP 服务器端
：
```go
AppNAME = 'UdpServerTest'
AppVERSION = '0.0.1'

function Main(arg)
    for i = 1, 10, 1 do
        local data = { temp = 20.15 , humi = 34}
        local err = applib:DataToUdp('UdpServer', applib:T2J(data))
        time:Sleep(100)
    end
    return 0
end


```

## 支持的平台
在下列系统上已经通过全面测试：

| 架构   | 操作系统             | 测试 |
| ------ | -------------------- | ---- |
| X86-64 | X86-64-Linux\Windows | 通过 |
| ARM64  | ARM-64-Linux         | 通过 |
| ARM32  | ARM-32-Linux         | 通过 |
| Mips   | Arm-Mips-Linux       | 通过 |


> [!WARNING]
> 除此之外，还可以在 Armbian、OpenWrt 等小众平台上编译成功。现阶段只针对**Ubuntu16.04**和**Ubuntu18.04**做了大量支持，其他的系统也许能编译成功但是没测试功能是否可用**

## 规则引擎
### 数据筛选

```lua
Actions = {
    function(args)
        print("return => ", rulexlib:JqSelect(".[] | select(.hum < 20)", data))
        return true, args
    end
}
```

### 数据中转

```lua
Actions = {
    function(args)
        -- 持久化到 MongoDb:
        data:ToMongo("45dd0c90f56d", data)
        -- 持久化到 Mysql:
        data:ToMysql("45dd0c90f56d", data)
        -- 推送化到 Kafka:
        data:ToKafka("45dd0c90f56d", data)
        return true, args
    end
}
```


## 社区

- QQ群：608382561
- 文档: <a href="https://hootrhino.github.io">[点我查看详细文档]</a>
- 微信：nilyouth( 加好友后进群, 别忘了来个小星星, 暗号：RULEX )
  <div style="text-align:center">
    <img src="./README_RES/wx.jpg" width="150px" />
  </div>

## 贡献者
鸣谢各位给RULEX贡献代码的社区开发者。

<a href="#">
  <img src="https://contributors-img.web.app/image?repo=hootrhino/rulex" />

</a>

<a href="#">
  <img src="https://contributors-img.web.app/image?repo=hootrhino/rulex-dashboard-vue-old" />
</a>

> [!NOTE]
> 如果您阅读过 RULEX 的源码，你会发现里面有很多很愚蠢的设计（比如对资源的状态管理、类型硬编码设计等），因为特殊历史原因导致了其设计上有一些很糟粕的地方，如有建议请不吝赐教，一起让这个框架更加优秀！同时未来随着版本的迭代，很多低级问题会逐步被重构完善。

## Star
<img src="https://starchart.cc/hootrhino/rulex.svg">
