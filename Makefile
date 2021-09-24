APP=rulex
VERSION=preview
.PHONY: all
all:
	make build

.PHONY: build
build:
	go mod tidy
	go generate
	go build -ldflags "-s -w" -o ${APP} main.go

.PHONY: xx
xx:
	make build

.PHONY: windows
windows:
	go mod tidy
	SET GOOS=windows
	go build -ldflags "-s -w" -o ${APP}.exe main.go

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

.PHONY: tag
tag:
	make package
	gh release create $(VERSION) ./${APP}-${VERSION}.zip -F changelog.md

.PHONY: clean-grpc
clean-grpc:
	rm ./rulexrpc/*.pb.go
	rm ./xstream/*.pb.go