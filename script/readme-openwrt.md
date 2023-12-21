# Openwrt daemon script
该脚本是RULEX的`Openwrt系统`操作脚本，处理RULEX的安装、启动、停止、卸载等。
## 基础使用
### 下载
将安装包解压:
```sh
unzip rulex-arm32linux-$VERSION.zip -d rulex
```
### 安装
```sh
./rulex-openwrt.sh install
```

### 卸载
```sh
./rulex-openwrt.sh uninstall
```

### 使用
下面的脚本一定要在root权限下执行,或者使用sudo。
```bash
# 启动
./rulex-openwrt.sh start
# 停止
./rulex-openwrt.sh stop
# 重启
./rulex-openwrt.sh restart
# 状态
./rulex-openwrt.sh status
```

## 守护进程
```sh
# 打开crontab
sudo crontab -e
# 输入
@reboot (export ARCHSUPPORT=EEKITH3 && /etc/init.d/rulex.service start > /var/log/rulex.log 2>&1)
```