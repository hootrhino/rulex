FROM golang:alpine3.9
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOPROXY="https://goproxy.cn,direct"
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk add build-base jq
RUN apk --no-cache add ca-certificates
ADD . /rulex/
WORKDIR /rulex/
RUN make

FROM golang:alpine3.9
LABEL name="RULEX"
LABEL author="wwhai"
LABEL email="cnwwhai@gmail.com"
LABEL homepage="https://github.com/hootrhino/rulex"
COPY --from=0 /rulex/ /rulex/
WORKDIR /rulex/

EXPOSE 2580
EXPOSE 2581
EXPOSE 2582
EXPOSE 2583
EXPOSE 2584
EXPOSE 2585
EXPOSE 2586
EXPOSE 2587
EXPOSE 2588
EXPOSE 2589

ENTRYPOINT ["./rulex", "run"]

