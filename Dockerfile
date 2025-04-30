FROM crpi-u4133maz63mkibc0.cn-chengdu.personal.cr.aliyuncs.com/shaco-go/golang:1.23.6 AS builder

ENV CGO_ENABLED=0
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /build

# 复制项目中的 go.mod 和 go.sum文件并下载依赖信息（这里只针对go语言而已）
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o app .

FROM crpi-u4133maz63mkibc0.cn-chengdu.personal.cr.aliyuncs.com/shaco-go/alpine:3.21.3

ARG ENV_VAR

ENV TZ Asia/Shanghai
ENV MY_ENV_VAR=${ENV_VAR}

WORKDIR /app

COPY --from=builder /build/config/development.yaml .
COPY --from=builder /build/config/production.yaml .
COPY --from=builder /build/app .

CMD ./app -c $MY_ENV_VAR
