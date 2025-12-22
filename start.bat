@echo off
REM Midgard Gateway 启动脚本 (Windows)

echo ==========================================
echo   Midgard Gateway 启动脚本
echo ==========================================

REM 检查 Go 是否安装
where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo 错误: 未找到 Go，请先安装 Go 1.22+
    pause
    exit /b 1
)

REM 创建必要的目录
if not exist "data" mkdir data
if not exist "config" mkdir config

REM 检查配置文件
if not exist "config\config.yaml" (
    echo 警告: 未找到 config\config.yaml，将使用默认配置
)

REM 设置默认环境变量
if "%PORT%"=="" set PORT=8080
if "%DATABASE_TYPE%"=="" set DATABASE_TYPE=sqlite
if "%DATABASE_DSN%"=="" set DATABASE_DSN=midgard.db
if "%REDIS_HOST%"=="" set REDIS_HOST=localhost
if "%REDIS_PORT%"=="" set REDIS_PORT=6379
if "%REDIS_PASSWORD%"=="" set REDIS_PASSWORD=
if "%REDIS_DB%"=="" set REDIS_DB=0
if "%ENABLE_FRONTEND%"=="" set ENABLE_FRONTEND=true

echo.
echo 配置信息:
echo   端口: %PORT%
echo   数据库类型: %DATABASE_TYPE%
echo   数据库 DSN: %DATABASE_DSN%
echo   Redis: %REDIS_HOST%:%REDIS_PORT%
echo   前端启用: %ENABLE_FRONTEND%
echo.

REM 检查依赖
echo 检查 Go 依赖...
if not exist "go.sum" (
    echo 下载 Go 依赖...
    go mod download
)

REM 构建应用
echo 构建应用...
go build -o midgard.exe main.go

REM 启动应用
echo.
echo 启动 Midgard Gateway...
echo ==========================================
midgard.exe

