APP=rulex
VERSION=0.0.1

.PHONY: all
all:
	make build

.PHONY: build
build:
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
