# 第一阶段：构建阶段
FROM golang:latest AS builder

# 设置工作目录
WORKDIR /app

# 复制本地的 Go 代码到容器中
COPY . .

# 构建 Go 应用
RUN CGO_ENABLED=0 go build -o main .

# 第二阶段：运行阶段
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 从构建阶段拷贝编译好的可执行文件到当前容器中
COPY --from=builder /app/main .

# 暴露应用运行时需要的端口
EXPOSE 8080

# 运行应用
CMD ["./main"]