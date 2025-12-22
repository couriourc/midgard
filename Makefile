.PHONY: help build run dev test clean docker-build docker-up docker-down

# 默认目标
help:
	@echo "Midgard Gateway - 可用命令:"
	@echo ""
	@echo "  开发命令:"
	@echo "    make dev          - 启动开发环境（包含 Redis）"
	@echo "    make build        - 构建应用"
	@echo "    make run          - 运行应用"
	@echo "    make test         - 运行测试"
	@echo "    make clean        - 清理构建文件"
	@echo ""
	@echo "  Docker 命令:"
	@echo "    make docker-build - 构建 Docker 镜像"
	@echo "    make docker-up    - 启动所有服务"
	@echo "    make docker-down  - 停止所有服务"
	@echo "    make docker-api   - 构建仅 API 模式的镜像"
	@echo ""
	@echo "  前端命令:"
	@echo "    make frontend-dev - 启动前端开发服务器"
	@echo "    make frontend-build - 构建前端"

# 开发环境
dev:
	@echo "启动开发环境..."
	docker-compose -f docker-compose.dev.yml up -d redis
	@echo "等待服务启动..."
	sleep 2
	@echo "启动应用..."
	go run main.go

# 构建
build:
	@echo "构建应用..."
	go build -o midgard main.go
	@echo "构建完成: ./midgard"

# 运行
run: build
	@echo "运行应用..."
	./midgard

# 测试
test:
	@echo "运行测试..."
	go test ./...

# 清理
clean:
	@echo "清理构建文件..."
	rm -f midgard midgard.exe
	rm -rf web/dist
	@echo "清理完成"

# Docker 构建（包含前端）
docker-build:
	@echo "构建 Docker 镜像（包含前端）..."
	docker-compose build

# Docker 构建（仅 API）
docker-api:
	@echo "构建 Docker 镜像（仅 API）..."
	docker-compose build --build-arg ENABLE_FRONTEND=false

# Docker 启动
docker-up:
	@echo "启动所有服务..."
	docker-compose up -d
	@echo "服务已启动"
	@echo "访问: http://midgard.localhost (如果启用了前端)"

# Docker 停止
docker-down:
	@echo "停止所有服务..."
	docker-compose down

# 前端开发
frontend-dev:
	@echo "启动前端开发服务器..."
	cd web && npm run dev

# 前端构建
frontend-build:
	@echo "构建前端..."
	cd web && npm run build

