# Midgard Gateway 快速启动指南

## 方式一：使用启动脚本（推荐）

### Linux/macOS
```bash
chmod +x start.sh
./start.sh
```

### Windows
```cmd
start.bat
```

## 方式二：使用 Makefile

### 开发环境（包含 Memcached）
```bash
make dev
```

### 构建并运行
```bash
make build
make run
```

### 查看所有可用命令
```bash
make help
```

## 方式三：手动启动

### 1. 安装依赖
```bash
# Go 依赖
go mod download

# 前端依赖（如果需要）
cd web && npm install
```

### 2. 配置环境（可选）
```bash
# 复制配置文件示例
cp config/config.example.yaml config/config.yaml
cp env.example .env

# 根据需要修改配置文件
```

### 3. 启动 Memcached（可选，如果使用缓存）
```bash
# 使用 Docker
docker-compose -f docker-compose.dev.yml up -d memcached

# 或使用本地安装的 Memcached
memcached -d -m 64 -p 11211
```

### 4. 运行应用
```bash
# 开发模式
go run main.go

# 或构建后运行
go build -o midgard main.go
./midgard
```

## 方式四：使用 Docker Compose

### 完整部署（包含前端）
```bash
docker-compose up -d
```

### 仅启动依赖服务（Memcached、PostgreSQL）
```bash
docker-compose -f docker-compose.dev.yml up -d
```

## 默认配置

- **端口**: 8080
- **数据库**: SQLite (midgard.db)
- **Memcached**: localhost:11211
- **前端**: 启用

## 环境变量

可以通过环境变量覆盖配置：

```bash
export PORT=8080
export DATABASE_TYPE=sqlite
export DATABASE_DSN=midgard.db
export MEMCACHED_HOST=localhost
export MEMCACHED_PORT=11211
export ENABLE_FRONTEND=true
```

## 访问应用

- **前端界面**: http://localhost:8080
- **API 端点**: http://localhost:8080/api
- **健康检查**: http://localhost:8080/health

## 常见问题

### 端口被占用
修改 `config/config.yaml` 中的 `server.port` 或设置环境变量 `PORT`

### Memcached 连接失败
确保 Memcached 服务正在运行，或设置 `MEMCACHED_HOST=""` 禁用缓存

### 数据库锁定错误
SQLite 在高并发下可能出现锁定，建议：
1. 使用 PostgreSQL
2. 或减少并发请求

