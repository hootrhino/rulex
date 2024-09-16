# RULEX Framework

### RULEX 是一个轻量级边缘网关开发框架，助力快速实现边缘数据处理和云边协同方案。

## 架构设计

<div style="text-align:center">
     
![image](https://github.com/user-attachments/assets/f60ee10f-d67d-43d5-a54d-3f2e0b13b26d)

<img src="./README_RES/structure.png"/>
</div>

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

## 社区

- QQ群：608382561
- 文档: <a href="https://hootrhino.github.io">[点我查看详细文档]</a>
- 微信：nilyouth( 加好友后进群, 别忘了来个小星星, 暗号：RULEX )
  <div style="text-align:center">
    <img src="./README_RES/wx.jpg" width="150px" />
  </div>

## Star
<img src="https://starchart.cc/hootrhino/rulex.svg">
