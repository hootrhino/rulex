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

| 硬件名          | 环境参数 | 示例                              |
| --------------- | -------- | --------------------------------- |
| EEKITH3版本网关 | EEKITH3  | `ARCHSUPPORT=EEKITH3 rulex run` |
| 树树莓派4B、4B+ | RPI4     | `ARCHSUPPORT=RPI4B rulex run`   |

> 警告: 这些属于板级高级功能，和硬件架构以及外设有关，默认关闭。 如果你自己需要定制，最好针对自己的硬件进行跨平台适配, 如果没有指定平台，可能会导致预料之外的结果。

## 常见函数

### EEKIT H3版本网关

1. GPIO 设置

   ```lua
   eekit:GPIOSet(Pin, Value)
   ```
   参数表

   | 参数名 | 类型 | 说明           |
   | ------ | ---- | -------------- |
   | Pin    | int  | GPIO引脚       |
   | Value  | int  | 高低电平, 0、1 |
2. GPIO 获取

   ```lua
   eekit:GPIOGet(Pin)
   ```
   | 参数名 | 类型 | 说明     |
   | ------ | ---- | -------- |
   | Pin    | int  | GPIO引脚 |
