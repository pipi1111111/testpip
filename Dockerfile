# 使用官方Go镜像作为构建环境
FROM golang:1.19 as builder

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum文件
COPY go.mod ./

# 下载依赖
RUN go mod download

# 复制源代码到容器中
COPY . .

# 编译应用程序
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o webhook .

# 使用scratch作为运行环境
FROM scratch

# 从构建阶段复制编译好的程序
COPY --from=builder /app/webhook /webhook

# 设置运行时的工作目录
WORKDIR /

# 暴露端口
EXPOSE 8088

# 运行webhook程序
ENTRYPOINT ["/webhook"]
