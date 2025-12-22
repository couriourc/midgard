#!/bin/bash

# Midgard Gateway 启动脚本

set -e

echo "=========================================="
echo "  Midgard Gateway 启动脚本"
echo "=========================================="

# 检查 Go 是否安装
if ! command -v go &> /dev/null; then
    echo "错误: 未找到 Go，请先安装 Go 1.22+"
    exit 1
fi

# 检查 Node.js 是否安装（如果需要构建前端）
if [ "$ENABLE_FRONTEND" != "false" ]; then
    if ! command -v node &> /dev/null; then
        echo "警告: 未找到 Node.js，前端可能无法正常工作"
    fi
fi

# 创建必要的目录
mkdir -p data
mkdir -p config

# 检查配置文件
if [ ! -f "config/config.yaml" ]; then
    echo "警告: 未找到 config/config.yaml，将使用默认配置"
fi

# 设置默认环境变量
export PORT=${PORT:-8080}
export DATABASE_TYPE=${DATABASE_TYPE:-sqlite}
export DATABASE_DSN=${DATABASE_DSN:-midgard.db}
export REDIS_HOST=${REDIS_HOST:-localhost}
export REDIS_PORT=${REDIS_PORT:-6379}
export REDIS_PASSWORD=${REDIS_PASSWORD:-}
export REDIS_DB=${REDIS_DB:-0}
export ENABLE_FRONTEND=${ENABLE_FRONTEND:-true}

echo ""
echo "配置信息:"
echo "  端口: $PORT"
echo "  数据库类型: $DATABASE_TYPE"
echo "  数据库 DSN: $DATABASE_DSN"
echo "  Redis: $REDIS_HOST:$REDIS_PORT"
echo "  前端启用: $ENABLE_FRONTEND"
echo ""

# 检查依赖
echo "检查 Go 依赖..."
if [ ! -f "go.sum" ]; then
    echo "下载 Go 依赖..."
    go mod download
fi

# 构建应用
echo "构建应用..."
go build -o midgard main.go

# 启动应用
echo ""
echo "启动 Midgard Gateway..."
echo "=========================================="
./midgard

