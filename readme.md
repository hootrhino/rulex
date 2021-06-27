# RulEngine X
RulEngine X 是一个轻量级规则引擎(名字看起来有点怪实际上是 rule + Engine 的组合词，中文发音为"若金克斯")。主要用来中转上游数据和吐出数据到目标点，可以理解为一个数据路由器。主要被设计用来做物联网网关或者服务端存在数据筛选的场景。
## 功能
- MQTT 数据输入
- HTTP 数据输入
- COAP 数据输入
- UDP 数据输入
- LUA 自定义业务逻辑支持
- SQL 字段筛选支持

## API接口
- HTTP RestFul API
## 管理界面
- Web dashboard

## 测试

```sh
make run
```
> 测试依赖于 main.go, 需要MQTT环境，本地装一个测试。

## 编译
```sh
make build
```
## Docker打包
```sh
make docker
```
## 制作压缩包
```sh
make package
```
## 统计代码
```sh
make clocs
```