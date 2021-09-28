APP=rulex
VERSION=V0.0.1
.PHONY: all
all:
	make build

.PHONY: build
build:
	go mod tidy
	chmod 755 ./gen_version.sh
	chmod +x ./gen_version.sh
	chmod 755 ./gen_proto.sh
	chmod +x ./gen_proto.sh
	sed -i "s/\r//" ./gen_proto.sh
	sed -i "s/\r//" ./gen_version.sh
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
	git tag -a "${VERSION}" -m "$(<./notes/${VERSION}.txt)"
	git commit -m "release: ${VERSION}"
	git push origin ${VERSION} -f

.PHONY: clean-grpc
clean-grpc:
	rm ./rulexrpc/*.pb.go
	rm ./xstream/*.pb.go