# RULEX Framework

**希望大家第一眼就看到这个说明 更多文档参考(国内可能需要科学上网): https://rulex.pages.dev**
> 如果您阅读过 RULEX 的源码，你会发现里面有很多很愚蠢的设计（比如对资源的状态管理、类型硬编码设计等），因为特殊历史原因导致了其设计上有一些很糟粕的地方，如有建议请不吝赐教，一起让这个框架更加优秀！同时未来随着版本的迭代，很多低级问题会逐步被重构完善。

### RULEX 是一个轻量级工业类边缘网关开发框架

## 架构设计

<div style="text-align:center">
<img src="./README_RES/structure.png"/>
</div>

## 快速开始

### HelloWorld
下面展示一个最简单的设备数据转发案例：
```go
AppNAME = 'UdpServerTest'
AppVERSION = '0.0.1'

function Main(arg)
    for i = 1, 10, 1 do
        local data = { temp = 20.15 , humi = 34}
        local err = applib:DataToUdp('UdpServer', applib:T2J(data))
        applib:Sleep(100)
    end
    return 0
end


```
这个 DEMO 展示了如何把一个简单的数据推到 UDP 服务器端


## 支持的平台

| 架构   | 操作系统             | 测试 |
| ------ | -------------------- | ---- |
| X86-64 | X86-64-Linux\Windows | 通过 |
| ARM64  | ARM-64-Linux         | 通过 |
| ARM32  | ARM-32-Linux         | 通过 |
| Mips   | Arm-Mips-Linux       | 通过 |

除此之外，还可以在 Armbian、OpenWrt 等小众平台上编译成功。

## 规则引擎

### 规则定义

```lua

function Success()
    -- do some things
end

function Failed(error)
    -- do some things
end

Actions = {
    function(data)
        return true, data
    end
}

```

### 数据筛选

```lua
Actions = {
    function(data)
        print("return => ", rulexlib:JqSelect(".[] | select(.hum < 20)", data))
        return true, data
    end
}
```

### 数据中转

```lua
Actions = {
    function(data)
        -- 持久化到 MongoDb:
        rulexlib:DataToMongo("45dd0c90f56d", data)
        -- 持久化到 Mysql:
        rulexlib:DataToMysql("45dd0c90f56d", data)
        -- 推送化到 Kafka:
        rulexlib:DataToKafka("45dd0c90f56d", data)
        return true, data
    end
}
```

### 云端计算

```lua
Actions = {
    function(data)
        -- PyTorch 训练数据:
        cloud:PyTorchTrainCNN(data)
        -- PyTorch 识别:
        local V = cloud:PyTorchCNN(data)
        print(V)
        return true, data
    end
}
```

## 社区

- QQ群：608382561
- 文档: <a href="https://rulex.pages.dev">[点我查看详细文档]</a>
- 微信：nilyouth( 加好友后进群, 别忘了来个小星星, 暗号：RULEX )
  <div style="text-align:center">
    <img src="./README_RES/wx.jpg" width="150px" />
  </div>
- 博客1：https://wwhai.gitee.io
- 博客2：https://wwhai.github.io

## 贡献者
鸣谢各位给RULEX贡献代码的大佬。

### RULEX
<a href="https://github.com/i4de/rulex/graphs/contributors">
  <img src="https://contributors-img.web.app/image?repo=i4de/rulex" />
</a>

### RULEX Other
<a href="https://github.com/i4de/rulex-dashboard/graphs/contributors">
  <img src="https://contributors-img.web.app/image?repo=i4de/rulex-dashboard" />
</a>

## Star
<img src="https://starchart.cc/i4de/rulex.svg">
