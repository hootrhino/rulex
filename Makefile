APP=rulenginex
VERSION=0.0.1

.PHONY: all
all:
	make build

.PHONY: build
build:
	echo "{" > metainfo.json
	echo '"version":"${VERSION}",' >> metainfo.json
	echo '"app":"${APP}",' >> metainfo.json
	echo '"build_ts":' '"$(shell date '+%Y-%m-%d %H:%M:%S')"' >> metainfo.json
	echo "}" >> metainfo.json

	go build -o ${APP}-${VERSION} main.go

.PHONY: run
run:
	go run -race main.go

.PHONY: docker
docker:
	docker build . -t ${APP}/${APP}:${VERSION}

.PHONY: clean
clean:
	go clean
	rm ${APP}-${VERSION}
	rm ${APP}-${VERSION}.zip
	rm metainfo.json

.PHONY: clocs
clocs:
	echo "# Lines" > clocs.md
	echo '```sh' >> clocs.md
	cloc ./ >> clocs.md
	echo '```'  >> clocs.md

.PHONY: package
package:
	zip ${APP}-${VERSION}.zip ./${APP}-${VERSION} ./metainfo.json