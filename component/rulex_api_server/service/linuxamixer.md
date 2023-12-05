<!--
 Copyright (C) 2023 wangwenhai

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

`amixer` 是 Linux 系统中用于控制音频设置的命令行工具。它可以用来调整音量、静音状态、选择音频输入/输出设备等。以下是一些常用的 `amixer` 命令示例：

1. **查看当前音量设置**：
   ```
   amixer get Master
   ```

2. **将音量设置为特定值**（以百分比为单位）：
   ```
   amixer set Master 50%
   ```

3. **增加音量**：
   ```
   amixer set Master 10%+
   ```

4. **减少音量**：
   ```
   amixer set Master 10%-
   ```

5. **静音/取消静音**：
   ```
   amixer set Master mute
   amixer set Master unmute
   ```

6. **选择音频设备**（如果系统支持多个音频设备）：
   ```
   amixer -c 1 set Master 50%
   ```

7. **查看详细音频控制选项**：
   ```
   amixer controls
   ```

8. **调整 PCM 音量**：
   ```
   amixer set PCM 50%
   ```

这些只是一些常见的 `amixer` 命令示例，实际上 `amixer` 提供了丰富的音频控制选项，你可以根据你的需求进行相应的调整。请注意，`amixer` 命令可能在不同的 Linux 发行版中略有不同，你可以根据你的系统进行适当的调整。如果你想要通过 Golang 来执行这些命令，可以使用 `os/exec` 包。