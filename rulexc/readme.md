# RULEXC 客户端工具使用文档
## RULEXC 简介
RULEXC 是 RULEX Client 的意思，是命令行下的客户端工具，帮助我们调试和管理RULEX。
## 使用文档

- 查看系统参数
```sh
go run ./rulexc.go system
```
- 查看入口列表
```sh
go run ./rulexc.go inends
```
- 查看单个入口
```sh
go run ./rulexc.go inends  '<id>'
```
- 查看出口列表
```sh
go run ./rulexc.go outends
```
- 查看单个出口
```sh
go run ./rulexc.go outends '<id>'
```
- 查看规则列表
```sh
go run ./rulexc.go rules
```
- 查看单个规则
```sh
go run ./rulexc.go rules '<id>'
```
- 查看插件列表
```sh
go run ./rulexc.go plugins
```
- 查看单个插件
```sh
go run ./rulexc.go plugins '<id>'
```
- 创建入口
```sh
go run ./rulexc.go inend-create --config  '<your json format config>'
```
- 创建出口
```sh
go run ./rulexc.go out-create --config  '<your json format config>'
```