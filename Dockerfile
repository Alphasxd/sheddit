FROM golang:1.22-alpine AS builder

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct

# 移动到工作目录：/build
WORKDIR /build

# 将代码复制到容器中
COPY . .

# 在含go环境的镜像中将代码编译成二进制可执行文件 app
RUN go mod tidy
RUN go build -o sheddit .

###################
# 接下来创建一个小镜像 因为此时我们已经得到可执行文件了 不需要有go环境了
###################
FROM debian:stretch-slim

# 从builder镜像中把/dist/app 拷贝到当前目录
COPY ./wait-for.sh /
COPY ./template template/
COPY ./static static/
COPY ./config config/
COPY --from=builder /build/sheddit /

# 需要暴露的端口
EXPOSE 8088

# 需要运行的命令
ENTRYPOINT ["/sheddit", "-f", "config/config.yaml"]