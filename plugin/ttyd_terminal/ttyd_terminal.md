# WEB 终端
该插件主要是用来在浏览器打开终端，方便debug、配置等，默认端口是：7681，最好不要对外开放这个端口，仅仅内网使用即可。

## 安装依赖
首先需要了解一点：RULEX 本身并没有 Terminal 功能，是借助了一个开源项目“ttyd”，因此需要提前安装好这个包，否则没有效果。下面是通过源码安装指令：
```bash
sudo apt-get install build-essential cmake git libjson-c-dev libwebsockets-dev
git clone https://github.com/tsl0922/ttyd.git
cd ttyd && mkdir build && cd build
cmake ..
make && sudo make install
```
## 配置参数
```ini
[plugin.ttyd]
#
# Enable
#
enable = true
#
# Server port
#
listen_port = 7681
```
### 参数说明
- enable: 是否开启
- listen_port：监听端口