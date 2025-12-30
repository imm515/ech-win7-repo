# ClawCloud Go ech服务端

Go 语言版本的甬哥 JS 服务端，用于 ClawCloud Run 部署。

### 5. 使用声明

**本软件仅用于：**
- ✅ 合法的网络加速
- ✅ 合理的网络访问需求
- ✅ 技术学习和研究

**本软件禁止用于：**
- ❌ 任何违法活动
- ❌ 绕过网络访问限制（学校、公司等）
- ❌ 访问非法内容
- ❌ 侵犯他人权益

### ⏱️ 删除提醒

**请在下载后 24 小时内删除本软件及相关文件。**

本软件仅用于技术学习和研究目的，请勿长期保存或滥用。

### 免责声明

开发者不对以下内容承担责任：
- 软件使用导致的任何法律问题
- 数据泄露或隐私泄露
- 账号被封禁或服务中断
- 网络访问问题或连接故障
- 使用软件造成的任何直接或间接损失

## 特性

- ✅ **100% 协议兼容** - 完全对齐甬哥 JS 版的 WebSocket 协议
- ✅ **Token 认证** - 通过 Sec-WebSocket-Protocol 头部验证
- ✅ **并发控制** - 限制最大 100 连接，防止资源滥用
- ✅ **超时统一** - WebSocket 和 TCP 超时统一为 5 分钟
- ✅ **心跳保活** - 响应客户端 Ping/Pong
- ✅ **健康检查** - `/health` 端点用于监控
- ✅ **优雅关闭** - 支持 SIGINT/SIGTERM 信号
- ✅ **完整日志** - 记录所有关键操作和错误

## 文件说明

```
ech3/
├── main.go          # 服务端源码
├── go.mod           # Go 模块定义
├── go.sum           # 依赖校验和
├── Dockerfile        # 容器镜像构建
└── README.md        # 本文档

.github/workflows/
└── deploy-server.yml # GitHub Actions 自动构建
```

## 环境变量

| 变量 | 说明 | 必填 | 默认值 |
|------|------|------|---------|
| `token` | 认证令牌 | 否 | 空（不验证） |
| `PORT` | 监听端口 | 否 | 3000 |

## 快速开始

### 1. 配置 GitHub Secrets

在 GitHub 仓库设置中添加以下 Secrets：

- `DOCKER_USERNAME`: Docker Hub 用户名
- `DOCKER_PASSWORD`: Docker Hub Access Token 或密码

### 2. 推送代码触发构建

```bash
git add ech3/
git commit -m "Add Go server for ClawCloud"
git push
```

代码推送后会自动触发 GitHub Actions 构建 Docker 镜像。

### 3. 在 ClawCloud Run 部署

1. **镜像地址**: `YOUR_USERNAME/echclawcloudrun:latest`
2. **环境变量**:
   - `token`: 您的认证令牌（必须与客户端一致）
   - `PORT`: `3000`
3. **容器端口**: `3000`
4. **资源配置**:
   - CPU: `0.1` vCPU
   - 内存: `128MB`

### 4. 本地客户端连接

使用 Win7 客户端连接：

```cmd
proxy-win7-x64.exe -f your-domain.clawcloudrun.com:443 -token your-secret-token
```

## 协议说明

### WebSocket 握手

客户端在握手时通过 `Sec-WebSocket-Protocol` 头部传递 token：

```
Sec-WebSocket-Protocol: your-secret-token
```

### 消息格式

| 指令 | 方向 | 格式 | 说明 |
|-------|------|------|------|
| CONNECT | 客户端→服务端 | `CONNECT:target\|firstFrame` | 连接目标地址 |
| CONNECTED | 服务端→客户端 | `CONNECTED` | 连接成功 |
| ERROR | 服务端→客户端 | `ERROR:message` | 错误通知 |
| CLOSE | 双向 | `CLOSE` | 关闭连接 |
| Binary | 双向 | ArrayBuffer | 二进制数据 |

### CONNECT 消息详解

```
CONNECT:host:port|firstFrameData
```

- `host:port`: 目标地址（支持 IPv6 格式 `[ipv6]:port`）
- `|`: 分隔符
- `firstFrameData`: 首帧数据（可为空）

## 开发

### 本地测试（需要 Go 环境）

```bash
cd ech3
go mod download
go run main.go
```

### 本地 Docker 测试

```bash
cd ech3
docker build -t echclawcloudrun .
docker run -p 3000:3000 -e token=your-secret-token echclawcloudrun
```

## 性能指标

- **最大并发**: 100 连接
- **缓冲区大小**: 32KB
- **超时时间**: 5 分钟
- **日志间隔**: 5 分钟

## 安全说明

- 仅允许无 Origin 的请求（防止浏览器直接连接）
- 支持 Token 认证
- 连接数限制防止滥用
- 日志记录所有连接事件

## 技术栈

- **Go**: 1.21
- **WebSocket**: gorilla/websocket v1.5.3
- **容器**: Alpine Linux
- **镜像大小**: ~10MB

## 许可证

本项目基于甬哥 JS 版本开发，仅供学习和个人使用。
