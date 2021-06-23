APP=rulenginex
VERSION=0.0.1
.PHONY: build
build:
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