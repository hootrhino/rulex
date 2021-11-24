# 项目构建
## 后端
后端是golang项目,首先装环境
```sh
sudo apt update -y
sudo apt install jq cloc protobuf-compiler gcc-mingw-w64-x86-64 gcc-arm-linux-gnueabi gcc-aarch64-linux-gnu -y
```
编译：
```
linux: make xx
windows: make windows
```
或者直接 `go build`.

## 前端
前端是NPM项目，clone代码以后执行：
```sh
npm install
npm build:prod
```
如果需要源码构建 dashboard，请运行 `release_pkg.sh` 脚本，或者单独运行前端项目的打包命令,将最新构建的前端项目下的dist目录下的文件全部复制到 `plugin\http_server\www` 目录下即可.
