FROM registry.cn-chengdu.aliyuncs.com/zhzy0518/golang:1.23 AS builder

ENV CGO_ENABLED=0
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /build

# 复制项目中的 go.mod 和 go.sum文件并下载依赖信息（这里只针对go语言而已）
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o app .

FROM registry.cn-chengdu.aliyuncs.com/zhzy0518/alpine:3.20.2

ENV TZ Asia/Shanghai

WORKDIR /app

COPY --from=builder /build/config/config.yaml .
COPY --from=builder /build/app .

CMD ["./app","-c","config.yaml"]