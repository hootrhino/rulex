# 注册系统进程
目前采用supervisor来管理rulex进程，具体配置如下:
- 首先安装:
```sh
pip install supervisor
```
- rulex 配置文件
```ini
[inet_http_server]
port=127.0.0.1:6001
username=admin
password=PaSs@rd

[supervisord]
logfile=/tmp/supervisord.log
logfile_maxbytes=100MB
logfile_backups=10
loglevel=info
pidfile=/tmp/supervisord.pid
nodaemon=false
minfds=1024
minprocs=200

[supervisorctl]
serverurl=unix:///tmp/supervisor.sock


[program:rulex]
command=~/rulex/rulex run
autostart=true
startsecs=10
autorestart=true
startretries=3
redirect_stderr=true
stdout_logfile_maxbytes=20MB
stdout_logfile_backups = 20
stdout_logfile=/opt/rulex/logs/rulex.log
stopasgroup=false
killasgroup=false

```

- 进程管理
```
supervisorctl
supervisorctl status all
```