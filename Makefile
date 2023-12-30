APP=rulex
# 获取操作系统信息
distro=$(shell cat /etc/os-release | grep PRETTY_NAME | cut -d '"' -f 2)
kernel=$(shell uname -r)
cpu=$(shell uname -m)
host=$(shell hostname)
ip=$(shell hostname -I)
memory=$(shell free -m | awk 'NR==2{printf "%.2fGB\n", $$2/1000}')
disk=$(shell df -h | awk '$$NF=="/"{printf "%s\n", $$2}')
arch=$(uname -m)
version=$(shell git describe --tags $(git rev-list --tags --max-count=1))
shortVersion=-X 'github.com/hootrhino/rulex/typex.MainVersion=$(version)'
.PHONY: all
all:
	@echo "\e[41m[*] Distro \e[0m: \e[36m ${distro} \e[0m"
	@echo "\e[41m[*] Arch   \e[0m: \e[36m ${arch} \e[0m"
	@echo "\e[41m[*] Kernel \e[0m: \e[36m ${kernel} \e[0m"
	@echo "\e[41m[*] Cpu    \e[0m: \e[36m ${cpu} \e[0m"
	@echo "\e[41m[*] Memory \e[0m: \e[36m ${memory} \e[0m"
	@echo "\e[41m[*] Host   \e[0m: \e[36m ${host} \e[0m"
	@echo "\e[41m[*] IP     \e[0m: \e[36m ${ip} \e[0m"
	@echo "\e[41m[*] Disk   \e[0m: \e[36m ${disk} \e[0m"
	go generate
	make build

.PHONY: build
build:
	CGO_ENABLED=1 GOOS=linux
	go generate
	go build -ldflags "$(shortVersion) -s -w" -v -o ${APP}

.PHONY: x64linux
x64linux:
	go generate
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=gcc\
	    go build -ldflags "$(shortVersion) -s -w" -v -o ${APP}-x64linux

.PHONY: windows
windows:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc\
	    go build -ldflags "$(shortVersion) -s -w" -o ${APP}-windows.exe

.PHONY: arm32
arm32:
	go generate
	CGO_ENABLED=1 GOOS=linux GOARCH=arm CC=arm-linux-gnueabi-gcc\
	    go build -ldflags "$(shortVersion) -s -w -linkmode external -extldflags -static" -o ${APP}-arm32linux

.PHONY: arm64
arm64:
	go generate
	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc\
	    go build -ldflags "$(shortVersion) -s -w -linkmode external -extldflags -static" -o ${APP}-arm64linux

.PHONY: mips32
mips32:
	go generate
	# sudo apt-get install gcc-mips-linux-gnu
	GOOS=linux GOARCH=mips CGO_ENABLED=1 CC=mips-linux-gnu-gcc\
	    go build -ldflags "$(shortVersion) -s -w -linkmode external -extldflags -static" -o ${APP}-mips32linux

.PHONY: mips64
mips64:
	go generate
	# sudo apt-get install gcc-mips-linux-gnu
	GOOS=linux GOARCH=mips64 CGO_ENABLED=1 CC=mips-linux-gnu-gcc\
	    go build -ldflags "$(shortVersion) -s -w -linkmode external -extldflags -static" -o ${APP}-mips64linux

.PHONY: mipsel
mipsle:
	go generate
	# sudo apt-get install gcc-mipsel-linux-gnu
	GOOS=linux GOARCH=mipsle CGO_ENABLED=1 GOMIPS=softfloat CC=mipsel-linux-gnu-gcc\
	    go build -ldflags "$(shortVersion) -s -w -linkmode external -extldflags -static" -o ${APP}-mipslelinux

.PHONY: release
release:
	bash ./release_pkg.sh

.PHONY: run
run:
	go run -race run

.PHONY: test
test:
	go test ${APP}/test -v

.PHONY: cover
cover:
	go test ${APP}/test -v -cover

.PHONY: clean
clean:
	go clean
	rm _release -rf
	rm ${APP}-arm32linux ${APP}-arm64linux *.db *.txt  *.txt.gz -rf upload/*
