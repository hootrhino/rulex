<!--
 Copyright (C) 2023 wwhai

 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU Affero General Public License as
 published by the Free Software Foundation, either version 3 of the
 License, or (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU Affero General Public License for more details.

 You should have received a copy of the GNU Affero General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.
-->

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