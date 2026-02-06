FROM golang:1.23.2 AS builder
MAINTAINER suqing <zz13168@hotmail.com>

WORKDIR /app
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

COPY . .
RUN make build

FROM centos:centos7.9.2009
WORKDIR /app

ENV TZ Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN mkdir -p /app/config && mkdir -p /qae/log

COPY --from=builder /app/bin/go_starter /app
COPY --from=builder /app/config/config_docker.yml /app/config/config.yml

RUN chmod +x /app/go_starter
CMD ["/app/go_starter", "http", "-c", "/app/config/config.yml"]