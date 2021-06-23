FROM golang:alpine

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY="https://goproxy.cn,direct"
RUN apk add build-base
MAINTAINER wwhai "cnwwhai@gmail.com"
ADD . /rulenginex
WORKDIR /rulenginex

CMD make run