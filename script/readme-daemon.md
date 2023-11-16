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
# 通用Linux
这个脚本是针对无systemctl支持的系统使用。

- 安装服务：`bash ./rulex_daemon.sh install`
- 启动服务：`bash ./rulex_daemon.sh start`
- 停止服务：`bash ./rulex_daemon.sh stop`
- 重启服务：`bash ./rulex_daemon.sh restart`
- 禁用服务：`bash ./rulex_daemon.sh disable`
- 启用服务：`bash ./rulex_daemon.sh enable`
- 卸载服务：`bash ./rulex_daemon.sh uninstall`