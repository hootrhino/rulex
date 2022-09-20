APP=rulex
.PHONY: all
all:
	make build

.PHONY: build
build:
	go mod tidy
	chmod +x gen_info.sh
	go generate
	CGO_ENABLED=1 GOOS=linux
	go build -v -ldflags "-s -w" -o ${APP}

.PHONY: xx
xx:
	make build

.PHONY: windows
windows:
	go mod tidy
	SET GOOS=windows
	SET CGO_ENABLED=1
	go build -ldflags "-s -w" -o ${APP}.exe

.PHONY: arm32
arm32:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm CC=arm-linux-gnueabi-gcc\
	    go build -ldflags "-s -w -linkmode external -extldflags -static"

.PHONY: arm64
arm64:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc\
	    go build -ldflags "-s -w -linkmode external -extldflags -static"

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
