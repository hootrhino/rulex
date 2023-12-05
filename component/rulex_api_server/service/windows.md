在 Windows 中，你可以使用 PowerShell 脚本来设置 WiFi 连接。以下是一个示例 PowerShell 脚本，演示了如何连接到一个 WiFi 网络：

```powershell
# WiFi 网络名称和密码
$ssid = "your-ssid"
$passphrase = "your-password"

# 使用 netsh 命令连接到 WiFi 网络
netsh wlan connect name="$ssid" ssid="$ssid" keyMaterial="$passphrase"
```

将上述脚本保存为 `.ps1` 文件，然后在 PowerShell 中运行脚本即可。请将 `your-ssid` 替换为实际的 WiFi 网络名称，将 `your-password` 替换为实际的 WiFi 密码。

请注意，运行 PowerShell 脚本需要管理员权限。你可以右键点击 PowerShell 图标，选择以管理员身份运行，然后运行脚本。

另外，如果你想在批处理脚本中调用 PowerShell 脚本，你可以将以下代码添加到批处理脚本中：

```batch
powershell.exe -ExecutionPolicy Bypass -File "path\to\your-script.ps1"
```

将 `path\to\your-script.ps1` 替换为实际的 PowerShell 脚本文件路径。这将以绕过执行策略的方式运行 PowerShell 脚本。