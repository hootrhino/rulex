# Linux systemctl
该脚本是RULEX的`systemctl`操作脚本，处理RULEX的安装、启动、停止、卸载等。
## 基础使用
将安装包解压:
```sh
unzip rulex-arm32linux-v0.6.2.zip -d rulex
```

下面的脚本一定要在root权限下执行,或者使用sudo。
- 安装
    ```sh
    ./rulex_systemctl.sh install
    ```
- 启动
    ```sh
    ./rulex_systemctl.sh start
    ```
- 状态
    ```sh
    ./rulex_systemctl.sh status
    ```
- 停止
    ```sh
    ./rulex_systemctl.sh stop
    ```
- 卸载
    ```sh
    ./rulex_systemctl.sh uninstall
    ```