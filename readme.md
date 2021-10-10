# RuleX

RuleX 是一个轻量级网关，支持多种数据接入以及数据流筛选，可以理解为一个数据路由器。

> 当前处于极其不稳定阶段,请勿尝试.
## 预览
### 登录
![res](README_RES/1.png)
### 首页
![res](README_RES/2.png)
### 服务
![res](README_RES/3.png)
### 资源
![res](README_RES/4.png)

## 快速开始
### 构建
```sh
git clone https://github.com/wwhai/rulex.git
cd rulex
make # on windows: make windows
```
> ProtoFile需要在Linux下编译, 需要安装: `sudo apt install protobuf-compiler -y`
### 启动
```sh
./rulex run ./conf/default.data
2021/09/20 17:09:05 cfg.go:24: [info] Init rulex config 
2021/09/20 17:09:05 cfg.go:34: [info] Rulex config init success. 
2021/09/20 17:09:05 utils.go:71: [info] 
 -----------------------------------------------------------     
~~~/=====\       ██████╗ ██╗   ██╗██╗     ███████╗██╗  ██╗       
~~~||\\\||--->o  ██╔══██╗██║   ██║██║     ██╔════╝╚██╗██╔╝       
~~~||///||--->o  ██████╔╝██║   ██║██║     █████╗   ╚███╔╝        
~~~||///||--->o  ██╔══██╗██║   ██║██║     ██╔══╝   ██╔██╗        
~~~||\\\||--->o  ██║  ██║╚██████╔╝███████╗███████╗██╔╝ ██╗       
~~~\=====/       ╚═╝  ╚═╝ ╚═════╝ ╚══════╝╚══════╝╚═╝  ╚═╝
-----------------------------------------------------------
2021/09/20 17:09:05 utils.go:74: [info] rulex start successfully
2021/09/20 17:09:05 http_api_server.go:139: [info] Http server started on http://127.0.0.1:2580
2021/09/20 17:09:05 grpc_resource.go:92: [info] RulexRpc resource started on [::]:2581
2021/09/20 17:09:05 coap_resource.go:71: [info] Coap resource started on [udp]:2582
2021/09/20 17:09:05 http_resource.go:47: [info] HTTP resource started on [0.0.0.0]:2583
2021/09/20 17:09:05 udp_resource.go:50: [info] UDP resource started on [0.0.0.0]:2584
```
> `./conf/default.data` 是已经设置好的测试数据,方便大家调试体验。
## Dashboard
```
浏览器输入：http://127.0.0.1:2580
```

## HTTP API
```
源码根目录下：`plugin\http_server\openapi.yml`
```

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
        print("return => ", stdlib:JqSelect(".[] | select(.hum < 20)", data))
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
        stdlib:DataToMongo("OUTEND_83775a94-9f64-4d37-be17-45dd0c90f56d", data)
        -- 持久化到 Mysql:
        stdlib:DataToMysql("OUTEND_83775a94-9f64-4d37-be17-45dd0c90f56d", data)
        -- 推送化到 Kafka:
        stdlib:DataToKafka("OUTEND_83775a94-9f64-4d37-be17-45dd0c90f56d", data)
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
<div style="text-align:center;">
    <a href="https://wwhai.github.io/rulex_doc_html">[点我查看详细文档]</a>
<div>
