## 第一阶段：构建Go应用
## 使用官方的Go基础镜像
#FROM golang:1.22-alpine AS builder
#
## 设置工作目录
#WORKDIR /app
#
## 复制go.mod和go.sum文件到工作目录
#COPY go.mod go.sum ./
#
## 下载项目依赖
#RUN go mod download
#
## 复制项目源代码到工作目录
#COPY . .
#
## 构建Go应用，生成可执行文件
#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 第二阶段：创建轻量级运行时镜像
FROM alpine:latest

# 安装必要的依赖（可选）
#RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /root/

# 从第一阶段的镜像中复制生成的可执行文件到当前工作目录
COPY --from=builder /app/main .

# 暴露应用监听的端口（根据实际情况修改）
EXPOSE 8081

# 定义容器启动时执行的命令
CMD ["./main"]
