# Windows 7 CDN优选客户端

一个为 Windows 7 设计的轻量级 CDN优选客户端，通过 WebSocket 连接 Cloudflare Worker CDN优选服务，支持 SOCKS5 和 HTTP 协议。

---

## ⚠️ 重要风险提示

**请在使用本软件前仔细阅读以下内容：**

### 1. 项目声明

- ⚠️ **本项目未经过充分测试和验证**，可能存在未知 bug 或安全问题
- 📚 **仅用于技术学习和研究目的**，不建议用于生产环境
- 🎯 **CDN优选范围限定为 Cloudflare CDN 加速优化**，不得用于其他用途

### 2. Cloudflare 使用风险

**服务条款限制：**
- ❌ 禁止用于违反当地法律法规的活动
- ❌ 禁止用于绕过学校/公司网络策略、地理封锁
- ❌ 禁止用于恶意攻击、钓鱼、传播恶意软件
- ❌ 禁止侵犯版权或传播非法内容

**账号封禁风险：**
如 Cloudflare 检测到以下情况，您的账号可能被**封禁**：
- 频繁的异常流量模式
- 大量来自同一 IP 的隧道连接
- 明显的CDN优选/VPN 特征流量
- WAF/安全系统触发的高风险警报

**免费版限制：**
- 每日最多 100,000 次请求
- 每次请求 CPU 时间限制 10ms
- 超出配额会被立即**拒绝服务**

### 3. 法律风险

- 使用CDN优选/VPN 技术在某些国家/地区可能**违反法律法规**
- 用户使用本软件所产生的一切法律责任**完全由用户自行承担**

### 4. 安全风险

- Cloudflare Workers 是**透明传输**，您的所有流量都可能被 Cloudflare 记录和分析
- **不要**通过本CDN优选传输敏感信息（密码、银行信息等）
- DoH 查询记录可能被 Cloudflare 查看

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

**请谨慎使用本软件，遵守当地法律法规和 Cloudflare 服务条款。**

---

## 特性

- ✅ **Windows 7 兼容** - 使用 TLS 1.3（可能需要系统更新），已禁用 ECH 扩展
- ✅ **双重CDN优选支持** - 同时支持 SOCKS5 和 HTTP CDN优选
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
// 修改第 3 行的CDN优选 IP 列表
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

SOCKS5 CDN优选：
- CDN优选类型：SOCKS5
- 地址：`127.0.0.1`
- 端口：`1080`

HTTP CDN优选：
- CDN优选类型：HTTP
- 地址：`127.0.0.1`
- 端口：`1080`

### Windows 系统CDN优选

1. 打开「控制面板」→「网络和 Internet」→「Internet 选项」
2. 点击「连接」选项卡
3. 点击「局域网设置」
4. 勾选「为 LAN 使用CDN优选服务器」
5. 地址：`127.0.0.1`，端口：`1080`

### Excel 2010 数据获取

Excel 2010 更适合使用 **HTTP CDN优选**：

1. Excel →「数据」→「获取外部数据」→「自 Web」
2. 在地址栏输入目标 URL
3. Excel 会自动使用系统CDN优选

## 架构说明

```
[应用] → [本机CDN优选 (1080)] → [WebSocket] → [Cloudflare Worker] → [目标网站]
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
3. 系统启用了 TLS 1.2/1.3 支持

### Q: Win7 上 TLS 握手失败？

**A:** 可能是 TLS 1.3 不兼容，解决方法：
1. 安装 Windows 7 的 KB4474419 安全更新（启用 TLS 1.3）
2. 或将 `main.go` 第 121 行的 `tls.VersionTLS13` 改为 `tls.VersionTLS12`

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

**A:** Excel 2010 推荐使用 HTTP CDN优选模式，系统CDN优选会自动生效。

## 技术细节

### 禁用的功能（原版本）

- 🔇 **ECH (Encrypted Client Hello)** - 已注释掉 `prepareECH()` 调用，禁用 ECH 扩展
- 🔇 **ECH 相关 TLS 配置** - 移除了 `EncryptedClientHelloConfigList` 和 `EncryptedClientHelloRejectionVerify` 字段
- ℹ️ **TLS 版本** - 使用 TLS 1.3（Win7 可能需要系统更新或修改为 TLS 1.2）

### 保留的组件

- ✅ **WebSocket 隧道** - 核心CDN优选功能
- ✅ **SOCKS5 协议** - 完整支持 CONNECT 和 UDP ASSOCIATE
- ✅ **HTTP CDN优选** - 支持 CONNECT 和 GET/POST 等方法
- ✅ **DoH** - DNS over HTTPS，用于 UDP DNS 查询

### 代码来源

基于 `ech2/原版客户端go源码` 修改，主要改动：
1. 注释掉 `main()` 函数中的 `prepareECH()` 调用（第 57 行）
2. 移除 `buildTLSConfigWithECH()` 中的 ECH 相关字段

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
├── main.go                    # Go 客户端源码（禁用 ECH 版本）
├── _worker.js                # Cloudflare Worker 脚本
├── start.bat                 # Windows 启动脚本
├── README.md                 # 本文档
├── ech2/                     # 原始代码目录
│   ├── 原版客户端go源码      # 原始客户端代码（已禁用 ECH）
│   └── g2                    # 修改建议文档
└── .github/
    └── workflows/
        └── main.yml          # GitHub Actions 自动编译
```

## 许可证

本项目基于 [CF_NAT](https://t.me/CF_NAT) 的原始代码开发，仅供学习和个人使用。

## 致谢

### 原始项目

- [CF_NAT](https://t.me/CF_NAT) - 原始项目来源
  - **Telegram 频道**: [@CF_NAT](https://t.me/CF_NAT)

### 参考与学习

- **甬哥** - 取消 ECH 思路来自甬哥视频，感谢分享和提供公益服务
  - **GitHub**: [yonggekkk/Cloudflare-vless-trojan](https://github.com/yonggekkk/Cloudflare-vless-trojan)
- **byJoey** - [ech-](https://github.com/byJoey/ech-) - 代码学习

### 技术支持

- [gorilla/websocket](https://github.com/gorilla/websocket) - WebSocket 实现
- Cloudflare Workers - CDN优选服务托管
