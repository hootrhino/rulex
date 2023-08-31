## Linux script
下面是一个示例的 Linux 守护进程脚本，用于管理名为 "rulex" 的程序。这个脚本可以启动、停止、重启和检查程序的状态。

```bash
#!/bin/bash

# 守护进程名称
DAEMON_NAME="rulex"

# 程序执行路径和参数
DAEMON_CMD="/path/to/rulex"
DAEMON_ARGS=""

# 守护进程启动函数
start() {
    echo "Starting $DAEMON_NAME..."
    $DAEMON_CMD $DAEMON_ARGS &
}

# 守护进程停止函数
stop() {
    echo "Stopping $DAEMON_NAME..."
    pkill -f "$DAEMON_CMD"
}

# 守护进程重启函数
restart() {
    stop
    start
}

# 检查守护进程状态
status() {
    if pgrep -f "$DAEMON_CMD" >/dev/null; then
        echo "$DAEMON_NAME is running."
    else
        echo "$DAEMON_NAME is not running."
    fi
}

# 主要逻辑
case "$1" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        restart
        ;;
    status)
        status
        ;;
    *)
        echo "Usage: $0 {start|stop|restart|status}"
        exit 1
        ;;
esac

exit 0
```

将上面的脚本保存为名为 "rulex_daemon.sh" 的文件，并将其中的 `/path/to/rulex` 替换为实际的程序执行路径。然后为脚本添加执行权限：

```bash
chmod +x rulex_daemon.sh
```

使用以下命令来管理守护进程：

- 启动守护进程：
  ```bash
  ./rulex_daemon.sh start
  ```

- 停止守护进程：
  ```bash
  ./rulex_daemon.sh stop
  ```

- 重启守护进程：
  ```bash
  ./rulex_daemon.sh restart
  ```

- 检查守护进程状态：
  ```bash
  ./rulex_daemon.sh status
  ```

请根据实际情况进行适当的调整，并确保你了解每个函数的作用。

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