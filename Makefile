APP=rulex
VERSION=preview
.PHONY: all
all:
	make build

.PHONY: build
build:
	go mod tidy
	go build -ldflags "-s -w" -o ${APP}-${VERSION} main.go

.PHONY: package
package:
	echo "{" > metainfo.json
	echo '"version":"${VERSION}",' >> metainfo.json
	echo '"app":"${APP}",' >> metainfo.json
	echo '"build_ts":' '"$(shell date '+%Y-%m-%d %H:%M:%S')"' >> metainfo.json
	echo "}" >> metainfo.json
	zip ${APP}-${VERSION}.zip ./${APP}-${VERSION} ./metainfo.json

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
	rm ${APP}-${VERSION}
	rm ${APP}-${VERSION}.zip
	rm metainfo.json
	rm coverage.out

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