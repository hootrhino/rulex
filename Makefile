APP=rulex
# 获取操作系统信息
distro=$(shell cat /etc/os-release | grep PRETTY_NAME | cut -d '"' -f 2)
kernel=$(shell uname -r)
cpu=$(shell uname -m)
host=$(shell hostname)
ip=$(shell hostname -I)
memory=$(shell free -m | awk 'NR==2{printf "%.2fGB\n", $$2/1000}')
disk=$(shell df -h | awk '$$NF=="/"{printf "%s\n", $$2}')


.PHONY: all
all:
	@echo "\e[41m[*] Distro \e[0m: \e[36m ${distro} \e[0m"
	@echo "\e[41m[*] Kernel \e[0m: \e[36m ${kernel} \e[0m"
	@echo "\e[41m[*] Cpu    \e[0m: \e[36m ${cpu} \e[0m"
	@echo "\e[41m[*] Memory \e[0m: \e[36m ${memory} \e[0m"
	@echo "\e[41m[*] Host   \e[0m: \e[36m ${host} \e[0m"
	@echo "\e[41m[*] IP     \e[0m: \e[36m ${ip} \e[0m"
	@echo "\e[41m[*] Disk   \e[0m: \e[36m ${disk} \e[0m"
	make build

.PHONY: build
build:
	chmod +x gen_info.sh
	go generate
	CGO_ENABLED=1 GOOS=linux
	go build -v -ldflags "-s -w" -o ${APP}

.PHONY: x64linux
x64linux:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=gcc\
	    go build -ldflags "-s -w -linkmode external -extldflags -static" -o ${APP}-x64linux

.PHONY: windows
windows:
	GOOS=windows go build -ldflags "-s -w" -o ${APP}.exe

.PHONY: arm32
arm32:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm CC=arm-linux-gnueabi-gcc\
	    go build -ldflags "-s -w -linkmode external -extldflags -static" -o ${APP}-arm32linux

.PHONY: arm64
arm64:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc\
	    go build -ldflags "-s -w -linkmode external -extldflags -static" -o ${APP}-arm64linux

.PHONY: mips32
mips32:
	# sudo apt-get install gcc-mips-linux-gnu
	GOOS=linux GOARCH=mips CGO_ENABLED=1 CC=mips-linux-gnu-gcc\
	    go build -ldflags "-s -w -linkmode external -extldflags -static" -o ${APP}-mips32linux

.PHONY: mips64
mips64:
	# sudo apt-get install gcc-mips-linux-gnu
	GOOS=linux GOARCH=mips64 CGO_ENABLED=1 CC=mips-linux-gnu-gcc\
	    go build -ldflags "-s -w -linkmode external -extldflags -static" -o ${APP}-mips64linux

.PHONY: mipsel
mipsle:
	# sudo apt-get install gcc-mipsel-linux-gnu
	GOOS=linux GOARCH=mipsle CGO_ENABLED=1 GOMIPS=softfloat CC=mipsel-linux-gnu-gcc\
	    go build -ldflags "-s -w -linkmode external -extldflags -static" -o ${APP}-mipslelinux

.PHONY: run
run:
	go run -race run

.PHONY: test
test:
	go test rulex/test -v

.PHONY: cover
cover:
	go test rulex/test -v -cover

.PHONY: clean
clean:
	go clean
	rm _release -rf
	rm *.db *log.txt -rf
