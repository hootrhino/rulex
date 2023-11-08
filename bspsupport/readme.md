# 跨平台

这里放置一些对特定硬件的支持库，一般指的是定制化网关产品。如果有不同操作系统上的实现库，建议统一放置此处。可参考下面的hello文件里面的程序。

## 当前兼容

### EEKIT 网关

EEKIT 是 RULEX 团队的默认硬件，操作系统为 `64位OpenWrt、Armbian`, CPU 架构为 `64位全志H3`。EEKIT 网关的lua标准库命名空间为 `eekit`。

### 树莓派4B+

除此之外，还对 `树莓派4B`的 GPIO 做了支持。树莓派的lua标准库命名空间为 `raspberry`。

## 环境变量

如果要启用硬件特性，需要在启动的时候加入 `ARCHSUPPORT` 环境变量来指定运行的版本, 例如要在 EEKIT-H3网关上运行：

```sh
ARCHSUPPORT=EEKITH3 rulex run
```

## 支持硬件列表

| 硬件名             | 环境参数  | 示例                              |
| ------------------ | --------- | --------------------------------- |
| EEKIT H3版本网关   | EEKITH3   | `ARCHSUPPORT=EEKITH3 rulex run`   |
| EEKIT T507版本网关 | EEKITT507 | `ARCHSUPPORT=EEKITT507 rulex run` |
| EEKIT T113版本网关 | EEKITT113 | `ARCHSUPPORT=EEKITT113 rulex run` |
| 树莓派4B、4B+      | RPI4      | `ARCHSUPPORT=RPI4B rulex run`     |
| 玩客云S805         | WKYS805   | `ARCHSUPPORT=WKYS805 rulex run`   |

> 警告: 这些属于板级高级功能，和硬件架构以及外设有关，默认关闭。 如果你自己需要定制，最好针对自己的硬件进行跨平台适配, 如果没有指定平台，可能会导致预料之外的结果。

## 常见函数

### EEKIT H3版本网关

1. GPIO 设置

   ```lua
   rhinopi:GPIOSet(Pin, Value)
   ```
   参数表

   | 参数名 | 类型 | 说明           |
   | ------ | ---- | -------------- |
   | Pin    | int  | GPIO引脚       |
   | Value  | int  | 高低电平, 0、1 |
2. GPIO 获取

   ```lua
   rhinopi:GPIOGet(Pin)
   ```
   | 参数名 | 类型 | 说明     |
   | ------ | ---- | -------- |
   | Pin    | int  | GPIO引脚 |

## 示例脚本
1. 玩客云WS1608
```lua
function Main(arg)
    while true do
        ws1608:GPIOSet("red", 1)
        time:Sleep(2000)
        ws1608:GPIOSet("red", 0)
        time:Sleep(2000)
    end
end

```
>必须使用这个系统：Linux aml-s812 5.9.0-rc7-aml-s812 #20.12 SMP Sun Dec 13 22:50:05 CST 2020 armv7l GNU/Linux, Armbian 20.12 Buster \l

1. EEKIT 网关
```lua
function Main(arg)
    while true do
        rhinopi:GPIOSet(6, 1)
        time:Sleep(2000)
        rhinopi:GPIOSet(7, 0)
        time:Sleep(2000)
    end
end

```