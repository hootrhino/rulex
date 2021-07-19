FROM golang:alpine
LABEL author="wwhai"
LABEL email="cnwwhai@gmail.com"
LABEL homepage="rulex.ezlinker.cn"
ENV APP_NAME="rulex"
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY="https://goproxy.cn,direct"
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add build-base zip

ADD . /rulex/
WORKDIR /rulex
RUN make build
RUN cp ./rulex-* ./rulex
EXPOSE 2580
CMD ./$APP_NAME run

