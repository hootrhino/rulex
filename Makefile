APP=rulex
VERSION=preview
.PHONY: all
all:
	make build

.PHONY: build
build:
	go mod tidy
	go build -ldflags "-s -w" -o ${APP}-${VERSION} main.go

.PHONY: xx
xx:
	make build

.PHONY: windows
windows:
	go mod tidy
	SET GOOS=windows
	go build -ldflags "-s -w" -o ${APP}-${VERSION}.exe main.go

.PHONY: package
package:
	echo ${VERSION} > VERSION
	zip ${APP}-${VERSION}.zip ./${APP}-${VERSION}

.PHONY: run
run:
	go run -race main.go run

.PHONY: docker
docker:
	docker build . -t ${APP}/${APP}:${VERSION} --rm

.PHONY: test
test:
	go test rulex/test -v

.PHONY: cover
cover:
	go test rulex/test -v -cover

.PHONY: clean
clean:
	go clean
	rm *.db
	rm ${APP}-${VERSION}

.PHONY: tag
tag:
	make package
	gh release create $(VERSION) ./${APP}-${VERSION}.zip -F changelog.md

.PHONY: proto
proto:
	go get google.golang.org/protobuf/cmd/protoc-gen-go@latest
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./rulexrpc/grpc_resource.proto
.PHONY: clean-grpc
clean-grpc:
	rm ./rulexrpc/*.pb.go