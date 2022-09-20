# RULEX Framework

### RULEX 是一个轻量级工业类边缘网关开发框架

> 在这里专门解释一下：什么是框架，框架顾名思义就是可以用来实现别的系统的一套 抽象结构。为什么说这个呢？因为很多朋友上来不看文档就说 RULEX 怎么连接XXX设备。实际上这个框架不是给最终用户来使用的，而是给企业开发人员来用的。你假设这是个WEB框架，就可以轻松实现网站等。同样的，**RULEX是一个集成了流处理和外挂设备驱动支持的物联网网关开发框架**，你可以用这套框架来实现你的产品逻辑。你在下面看到的界面，以及在release界面下载的那个可执行程序，实际上是一个最小Demo，而不是完整的产品。就好比Golang的GIN这个框架，为了教你怎么用，专门写了个Demo网站，而不是说GIN是网站。
> 也许后期会有一些开源作者基于这个框架开发一些硬件产品，我们会收集起来，供给大家把玩参考。

**希望大家第一眼就看到这个说明 更多文档参考(国内可能需要科学上网): https://rulex.pages.dev**

## 架构设计

<div style="text-align:center">
<img src="./README_RES/structure.png"/>
</div>

## 快速开始

### 构建(Linux)

```sh
git clone https://github.com/wwhai/rulex.git
cd rulex
make
```

## 支持的平台

| 平台    | 架构   | 编译测试 |
| ------- | ------ | -------- |
| Windows | X86-64 | 通过     |
| Linux   | X86-64 | 通过     |
| ARM64   | ARM-64 | 通过     |
| ARM32   | ARM-32 | 通过     |
| MacOS   | X86-64 | 通过     |
| 其他    | 未知   | 未知     |

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
function Success()
    -- do some things
end

function Failed(error)
    -- do some things
end

Actions = {
    function(data)
        print("return => ", rulexlib:JqSelect(".[] | select(.hum < 20)", data))
        return true, data
    end
}
```

### 数据中转

```lua
function Success()
    -- do some things
end

function Failed(error)
    -- do some things
end

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
function Success()
    -- do some things
end

function Failed(error)
    -- do some things
end

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

## 详细文档

`<a href="https://rulex.pages.dev">`[点我查看详细文档]`</a>`

## 社区

- QQ群：608382561
- 微信：nilyouth( 加好友后进群, 暗号：RULEX )

  <div style="text-align:center">
    <img src="./README_RES/wx.jpg" width="150px" />
    </div>
- 博客1：https://wwhai.gitee.io
- 博客2：https://wwhai.github.io

## Star

<img src="https://starchart.cc/i4de/rulex.svg">
