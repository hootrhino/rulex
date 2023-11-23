# Linux daemon script
该脚本是RULEX的`通用Linux系统`操作脚本，处理RULEX的安装、启动、停止、卸载等。
## 基础使用
### 下载
将安装包解压:
```sh
unzip rulex-arm32linux-$VERSION.zip -d rulex
```
### 安装
```sh
./rulex-daemon.sh install
```

### 卸载
```sh
./rulex-daemon.sh uninstall
```

### 使用
下面的脚本一定要在root权限下执行,或者使用sudo。
```bash
# 启动
./rulex-daemon.sh start
# 停止
./rulex-daemon.sh stop
# 重启
./rulex-daemon.sh restart
# 状态
./rulex-daemon.sh status
```
