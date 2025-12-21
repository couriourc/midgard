# 前端构建阶段（可选）
FROM node:18-alpine AS frontend-builder

WORKDIR /app/web

# 复制前端依赖文件
COPY web/package*.json ./

# 安装依赖
RUN npm install

# 复制前端源代码
COPY web .

# 构建前端应用
RUN npm run build

# 创建空目录作为占位符（用于禁用前端时）
FROM alpine:latest AS frontend-empty
RUN mkdir -p /app/web/dist

# Go 应用构建阶段
FROM golang:1.22.4-alpine AS backend-builder

WORKDIR /app

# 复制 Go 依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制 Go 源代码
COPY . .

# 构建参数：是否包含前端（默认：true）
ARG ENABLE_FRONTEND=true

# 根据构建参数选择复制前端文件或空目录
COPY --from=frontend-builder /app/web/dist ./web/dist

# 构建 Go 应用
RUN go build -o midgard main.go

# 最终镜像
FROM alpine:latest

WORKDIR /app

# 安装必要的运行时依赖
RUN apk add --no-cache ca-certificates tzdata

# 复制 Go 应用
COPY --from=backend-builder /app/midgard .

# 复制配置文件
COPY --from=backend-builder /app/config ./config

# 复制前端文件（如果存在）
ARG ENABLE_FRONTEND=true
COPY --from=backend-builder /app/web/dist ./web/dist

# 创建数据目录
RUN mkdir -p /app/data

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./midgard"]
