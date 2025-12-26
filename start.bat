@echo off
chcp 65001 >nul
title Windows 7 CDN优选客户端
color 0A

echo.
echo ================================
echo    Windows 7 CDN优选客户端
echo ================================
echo.

REM ================= 配置区域 =================

REM 本地监听地址（默认：127.0.0.1:1080）
REM 可选端口：1080, 7890, 10808 等
set LISTEN_ADDR=127.0.0.1:1080

REM Cloudflare Worker 地址（必填！）
REM 格式：your-worker.workers.dev:443
REM 请修改为你的 Worker 域名
set SERVER_ADDR=your-worker.workers.dev:443

REM 可选：指定服务器 IP（绕过 DNS 解析）
REM 使用场景：域名解析失败时，填入 Worker 的优选 IP
REM 格式：1.2.3.4
REM set SERVER_IP=1.2.3.4

REM 可选：认证令牌（与 Worker 脚本中的 token 一致）
REM 如果 Worker 启用了 token 验证，请填写
REM 格式：your-token-here
REM set TOKEN=your-token-here

REM ========================================

echo [配置] 本地监听: %LISTEN_ADDR%
echo [配置] 服务器地址: %SERVER_ADDR%
if defined SERVER_IP echo [配置] 服务器 IP: %SERVER_IP%
if defined TOKEN echo [配置] 认证令牌: 已设置
echo.

REM 检查可执行文件是否存在
if not exist "proxy-win7-x64.exe" (
    echo [错误] 找不到 proxy-win7-x64.exe
    echo.
    echo 请先从 GitHub Actions 下载编译好的程序：
    echo   1. 访问仓库的 Actions 页面
    echo   2. 找到最新的构建任务
    echo   3. 下载 proxy-win7-x64.exe
    echo.
    pause
    exit /b 1
)

REM 检查服务器地址是否为默认值
if "%SERVER_ADDR%"=="your-worker.workers.dev:443" (
    echo [警告] 服务器地址仍为默认值！
    echo [警告] 请先编辑本文件，修改 SERVER_ADDR 为你的 Worker 地址
    echo [警告] 按 Enter 继续启动（或 Ctrl+C 取消）...
    pause
)

echo [启动] 正在启动CDN优选服务...
echo [提示] 按 Ctrl+C 停止服务
echo.
echo ================================
echo.

REM 构建命令行参数
set ARGS=-l %LISTEN_ADDR% -f %SERVER_ADDR%

if defined SERVER_IP (
    set ARGS=%ARGS% -ip %SERVER_IP%
)

if defined TOKEN (
    set ARGS=%ARGS% -token %TOKEN%
)

proxy-win7-x64.exe %ARGS%

if %errorlevel% neq 0 (
    echo.
    echo [错误] 服务异常退出，错误代码: %errorlevel%
    echo.
    pause
)
