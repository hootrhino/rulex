# 安装脚本
此处介绍两类应用安装脚本。
## Systemctl
Systemctl是Ubuntu系自带的进程管理器，这也是大部分人能接触到的最简单的一种，下面给出个示例。
1. 配置脚本
    ```ini
    [Unit]
    Description=rulex
    After=network-online.target rc-local.service nss-user-lookup.target
    Wants=network-online.target

    [Service]
    User=root
    Type=simple
    WorkingDirectory=/usr/local/rulexapp
    ExecStart=/usr/local/rulexapp/rulex
    Restart=on-failure
    RestartSec=5s
    [Install]
    WantedBy=multi-user.target

    ```
2. 操作指令
    ```sh
    sudo systemctl start rulex
    sudo systemctl enable rulex
    sudo systemctl status rulex
    ```

## Linux 原生脚本
Linux 原生脚本一半放在 `/etc/network/interfaces.d`目录下。
```sh
#! /bin/sh
APP_NAME="/root/rulex"
while true; do
    APP_PROCESS_COUNT=`ps aux | grep ${APP_NAME} | grep -v grep |wc -l`
    if [ "${APP_PROCESS_COUNT}" -lt "1" ];then
        ${APP_NAME} -c 1 &
        elif [ "${APP_PROCESS_COUNT}" -gt "1" ];then
        killall -9 $APP_NAME
    fi
    sleep 5
done
```