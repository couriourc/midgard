# Midgard Gateway

Midgard Gateway 是一个基于 Go、Traefik、Memcached 和 Vue 3 的网关代理管理工具。

## 项目简介
技术栈使用 go1.22.4(gin+gorm+ sqlite/postgres) +traefik +memcached + shadcn dashboard (vue3 vite + shadcn ) 的网关代理管理工具，具有以下功能：

1. 支持通过导入 openapi.json 或者提供 openapi.json 进行解析，生成一个 collection，一个collection相当于一组，可以设置代理的 endpoint。在 collection 中可以设置对外的网关前缀，这样就可以外部调用我的端口，我进行完整的代理转发，记录拦截外部的请求日志，并记录请求耗时。

2. 支持多个 collection ，支持停用 collection 以及启用，从而控制是否能请求我对应的 collection 进行访问。

3. 在 collection 中能配置，/health 的地址，来判断服务是否正常运行，并支持设置 collection 记录的方式为滚动记录，以及限制条目数目，并支持能够设置是否通过 memcached 来进行缓存，缓存的依据包括，json 化参数体等方式


整体最后采用 docker 部署
## 功能特性

1. **OpenAPI 导入**：支持通过 URL 或上传 JSON 文件导入 OpenAPI 规范，自动生成代理端点
2. **集合管理**：
   - 创建、编辑、删除集合
   - 启用/停用集合控制访问权限
   - 设置对外网关前缀
3. **健康检查**：配置健康检查路径和间隔，自动监控后端服务状态
4. **日志记录**：
   - 记录请求详细信息（路径、方法、状态码、耗时等）
   - 支持滚动日志和条目限制
5. **缓存支持**：
   - 通过 Memcached 缓存请求响应
   - 可配置缓存策略（参数、请求体或全部）
6. **Dashboard**：直观的 Web 界面管理所有功能

## 技术栈

- **后端**：Go 1.22.4
- **API 网关**：Traefik
- **缓存**：Memcached
- **前端**：Vue 3 + Vite + Shadcn UI
- **部署**：Docker

## 项目结构

```
midgard/
├── config/           # 配置管理
├── internal/
│   ├── api/          # API 路由和处理
│   ├── collection/   # 集合管理
│   └── proxy/        # 代理功能
├── web/              # 前端代码
├── Dockerfile        # 应用 Dockerfile
├── docker-compose.yml # Docker 部署配置
├── go.mod            # Go 依赖
└── main.go           # 主程序入口
```

## 快速开始

### 使用 Docker 部署

1. 克隆项目
   ```bash
git clone <repository-url>
cd midgard
```

2. 启动服务
   ```bash
docker-compose up -d
```

3. 访问 Dashboard
   ```
http://midgard.localhost
```

### 手动构建

#### 后端

1. 安装依赖
   ```bash
go mod download
```

2. 构建
   ```bash
go build -o midgard
```

3. 运行
   ```bash
./midgard
```

#### 前端

1. 进入前端目录
   ```bash
cd web
```

2. 安装依赖
   ```bash
npm install
```

3. 开发模式
   ```bash
npm run dev
```

4. 构建生产版本
   ```bash
npm run build
```

## API 文档

### 集合管理

- `GET /api/collections` - 获取所有集合
- `POST /api/collections` - 创建新集合
- `GET /api/collections/{id}` - 获取单个集合
- `PUT /api/collections/{id}` - 更新集合
- `DELETE /api/collections/{id}` - 删除集合
- `POST /api/collections/{id}/toggle` - 启用/停用集合
- `POST /api/collections/{id}/import-openapi` - 导入 OpenAPI 规范

### 代理请求

```
/proxy/{prefix}/{path}
```

其中 `prefix` 是 collection 配置的对外网关前缀。

## 配置

配置文件位于 `config/config.yaml`

```yaml
server:
  port: 8080

database:
  type: sqlite
  dsn: midgard.db
  # 或使用 PostgreSQL:
  # type: postgres
  # host: localhost
  # port: 5432
  # user: postgres
  # password: postgres
  # dbname: midgard

memcached:
  host: localhost
  port: 11211

log:
  level: info
  max_entries: 1000
  rolling: true
```

## 许可证

MIT
