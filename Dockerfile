FROM golang:alpine AS builder
# FROM seekwe/go-builder:latest AS builder
LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /zlsgo-app

ADD go.mod .
ADD go.sum .
RUN go mod download

COPY . .

# 编译指令
RUN go build -ldflags="-s -w" -o /app/zls .

FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata

ENV TZ Asia/Shanghai

WORKDIR /app

COPY --from=builder /app/zls /app/zls

# 暴露一个端口
EXPOSE 8181

# 执行程序
CMD ["./zls"]

# 根据 Dockerfile 生成 Docker 镜像
# docker build -t zlsapp:v1 -f ./Dockerfile  .

# 启动容器
# docker run --rm -it -p 8181:8181 zlsapp:v1

# 进入容器
# docker run --rm -it zlsapp:v1 sh