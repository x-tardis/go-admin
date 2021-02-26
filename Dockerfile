FROM golang:1.15.8-alpine3.12 AS builder

ENV GOPROXY https://goproxy.io,https://goproxy.cn,direct
ARG ALPINE_VERSION=v3.12
WORKDIR /release

# build depends and then building
RUN echo https://mirrors.aliyun.com/alpine/${ALPINE_VERSION}/main > /etc/apk/repositories \
     	&& echo https://mirrors.aliyun.com/alpine/${ALPINE_VERSION}/community >> /etc/apk/repositories
RUN apk add git \
            make \
            bash \
            tzdata

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY . .
RUN pwd && ls && make

# 运行环境
FROM frolvlad/alpine-glibc:alpine-3.12_glibc-2.31

MAINTAINER thinkgos "thinkgo@aliyun.com"

ENV TZ=Asia/Shanghai
WORKDIR /root

# 设置时区 Ubuntu和Debian系统, CentOS只需第一句
RUN apk --no-cache add ca-certificates bash
#    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
#    && echo "Asia/Shanghai" > /etc/timezone

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /release/config config
COPY --from=builder /release/static static
COPY --from=builder /release/go-admin go-admin

EXPOSE 8000
CMD ["./go-admin","server","-c", "./config/config.yaml"]

