# Windows 7 代理客户端

一个为 Windows 7 设计的轻量级代理客户端，通过 WebSocket 连接 Cloudflare Worker 代理服务，支持 SOCKS5 和 HTTP 协议。

## 特性

- ✅ **Windows 7 完全兼容** - 使用 TLS 1.2，无需 ECH 扩展
- ✅ **双重代理支持** - 同时支持 SOCKS5 和 HTTP 代理
- ✅ **DoH 支持** - 自动通过 HTTPS 转发 DNS 查询到 Cloudflare
- ✅ **轻量级** - 单文件无依赖，编译后仅几 MB
- ✅ **支持 IP 直连** - 可指定服务器 IP 绕过 DNS 解析
- ✅ **Token 认证** - 支持 WebSocket 子协议令牌验证

## 兼容性

| 平台 | 支持 |
|------|------|
| Windows 7 | ✅ |
| Windows 8/8.1 | ✅ |
| Windows 10/11 | ✅ |
| Linux (交叉编译) | ✅ |

## 快速开始

### 1. 准备 Cloudflare Worker

将 `_worker.js` 部署到 Cloudflare Workers：

```javascript
// 修改第 3 行的代理 IP 列表
const CF_FALLBACK_IPS = ['your-proxy-ip.com'];

// 修改第 13 行的认证令牌（可选）
const token = 'your-token-here';
```

### 2. 编译客户端

#### 使用 GitHub Actions（推荐）

**方式一：推送代码自动触发**

```bash
git add .
git commit -m "Initial commit"
git push
```

**方式二：手动触发（推荐）**

1. 将代码推送到 GitHub 仓库
2. 访问仓库的 **Actions** 页面
3. 选择左侧的 **"Auto Build Win7 x64 & Sync Deps"** 工作流
4. 点击 **"Run workflow"** 按钮
5. 选择分支，点击绿色的 **"Run workflow"** 开始编译

编译完成后在 Actions 页面下载 `proxy-win7-x64.exe`。

### 3. 运行客户端

编辑 `start.bat` 配置参数：

```batch
REM 本地监听地址（默认：127.0.0.1:1080）
set LISTEN_ADDR=127.0.0.1:1080

REM Cloudflare Worker 地址（必填！）
set SERVER_ADDR=your-worker.workers.dev:443

REM 可选：指定服务器 IP（绕过 DNS）
set SERVER_IP=1.2.3.4

REM 可选：认证令牌
set TOKEN=your-token-here
```

双击 `start.bat` 运行。

#### 命令行运行（可选）

```batch
# 基础使用
proxy-win7-x64.exe -l 127.0.0.1:1080 -f your-worker.workers.dev:443

# 使用 IP 直连
proxy-win7-x64.exe -l 127.0.0.1:1080 -f your-worker.workers.dev:443 -ip 1.2.3.4

# 使用 Token 认证
proxy-win7-x64.exe -l 127.0.0.1:1080 -f your-worker.workers.dev:443 -token your-token-here
```

## 参数说明

| 参数 | 说明 | 示例 | 必填 |
|------|------|------|------|
| `-l` | 本地监听地址 | `127.0.0.1:1080` | 否（默认 30000） |
| `-f` | 服务器地址 | `worker.workers.dev:443` | **是** |
| `-ip` | 指定服务器 IP | `1.2.3.4` | 否 |
| `-token` | 认证令牌 | `secret-token` | 否 |

## 配置应用程序

### 浏览器（Firefox/Chrome/Edge）

SOCKS5 代理：
- 代理类型：SOCKS5
- 地址：`127.0.0.1`
- 端口：`1080`

HTTP 代理：
- 代理类型：HTTP
- 地址：`127.0.0.1`
- 端口：`1080`

### Windows 系统代理

1. 打开「控制面板」→「网络和 Internet」→「Internet 选项」
2. 点击「连接」选项卡
3. 点击「局域网设置」
4. 勾选「为 LAN 使用代理服务器」
5. 地址：`127.0.0.1`，端口：`1080`

### Excel 2010 数据获取

Excel 2010 更适合使用 **HTTP 代理**：

1. Excel →「数据」→「获取外部数据」→「自 Web」
2. 在地址栏输入目标 URL
3. Excel 会自动使用系统代理

## 架构说明

```
[应用] → [本机代理 (1080)] → [WebSocket] → [Cloudflare Worker] → [目标网站]
                                    ↓
                              [DoH (cloudflare-dns.com)]
```

- SOCKS5 UDP 请求（DNS）通过 DoH 转发到 Cloudflare
- 其他 TCP 流量通过 WebSocket 隧道转发

## 常见问题

### Q: 编译后无法在 Win7 运行？

**A:** 确保：
1. 使用 Go 1.20.x 版本编译
2. 目标系统安装了最新的根证书更新
3. 系统启用了 TLS 1.2 支持

### Q: 连接失败？

**A:** 检查：
1. Worker 脚本中的 `CF_FALLBACK_IPS` 是否正确
2. 防火墙是否允许 443 端口
3. 尝试使用 `-ip` 参数指定服务器 IP

### Q: DNS 解析失败？

**A:** DoH 可能被阻断，尝试：
1. 修改 `queryDoHForProxy` 函数中的 DoH 服务器
2. 使用 `-ip` 参数指定服务器 IP

### Q: Excel 无法获取数据？

**A:** Excel 2010 推荐使用 HTTP 代理模式，系统代理会自动生效。

## 技术细节

### 移除的组件（原版本）

- ❌ ECH (Encrypted Client Hello) - Win7 不支持
- ❌ TLS 1.3 强制要求 - 改为 TLS 1.2
- ❌ HTTPS DNS 记录解析 - 不再需要

### 保留的组件

- ✅ WebSocket 隧道 - 核心代理功能
- ✅ SOCKS5 协议 - 完整支持 CONNECT 和 UDP ASSOCIATE
- ✅ HTTP 代理 - 支持 CONNECT 和 GET/POST 等方法
- ✅ DoH - DNS over HTTPS，用于 UDP DNS 查询

### 协议兼容性

| 消息类型 | Go 客户端 | JS 服务端 |
|----------|-----------|-----------|
| `CONNECT:target\|data` | ✓ | ✓ |
| 二进制数据 | ArrayBuffer | ✓ |
| `CLOSE` | ✓ | ✓ |
| `ERROR:` | ✓ | ✓ |

## 项目结构

```
ech-win7-repo/
├── main.go              # Go 客户端源码
├── _worker.js          # Cloudflare Worker 脚本
├── start.bat           # Windows 启动脚本
├── README.md           # 本文档
├── g1                  # Gemini 建议文档（参考）
└── .github/
    └── workflows/
        └── main.yml    # GitHub Actions 自动编译
```

## 许可证

本项目基于 [CF_NAT](https://t.me/CF_NAT) 的原始代码开发，仅供学习和个人使用。

## 致谢

### 原始项目

- [CF_NAT](https://t.me/CF_NAT) - 原始项目来源
  - **Telegram 频道**: [@CF_NAT](https://t.me/CF_NAT)

### 参考与学习

- **CM** - 服务支持
- **YGkkk** - 取消ECH思路来自甬哥视频
- [ech-](https://github.com/byJoey/ech-) - JS 代码参考

### 技术支持

- [gorilla/websocket](https://github.com/gorilla/websocket) - WebSocket 实现
- Cloudflare Workers - 代理服务托管
