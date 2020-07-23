FROM golang:1.14-alpine AS build

LABEL maintainer=""

WORKDIR /build

ADD . /build

ENV GO111MODULE=on \
    GOPROXY="https://goproxy.cn" \
    GOARCH=amd64 \
    CGO_ENABLED=1 \
    GOOS=linux 

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories \
&& apk add --no-cache git gcc musl-dev protoc libprotoc protobuf libprotobuf protobuf-dev \
&& go get -u github.com/golang/protobuf/proto \
&& go get -u github.com/golang/protobuf/protoc-gen-go \
&& export PATH="/root/go/bin:$PATH"

RUN go build -tags musl -o wechat_server_linux_amd64 /build

FROM alpine:3.12 AS prod

COPY --from=build /build/wechat_server_linux_amd64 /usr/local/bin/
COPY --from=build /build/config.toml /zhimiao/
COPY --from=build /build/docker-entrypoint.sh /usr/local/bin/

RUN set -eux \
    && chmod a+x /usr/local/bin/docker-entrypoint.sh \
    && addgroup -S -g 1000 zhimiao \
    && adduser -S -G zhimiao -u 1000 zhimiao \
    && sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories \
    && apk upgrade -U -a --no-cache \
    && apk add --no-cache \
    'su-exec>=0.2'

EXPOSE 1316

WORKDIR /zhimiao

ENTRYPOINT ["docker-entrypoint.sh"]

CMD ["wechat_server_linux_amd64"]
