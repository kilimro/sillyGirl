# MianPro

[![Version](https://img.shields.io/badge/version-v1.1.9-1677ff)](VERSION)
[![Go](https://img.shields.io/badge/Go-1.18+-00ADD8?logo=go)](go.mod)
[![Docker](https://img.shields.io/badge/Docker-pull-2496ED?logo=docker)](https://hub.docker.com/r/kilimro/mianpro)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

[English](README.md) | 中文

基于 Go 的多平台聊天机器人框架。通过统一的管理后台连接微信、飞书、QQ、Telegram、钉钉等多个平台。

## 功能特性

- **多平台适配器** — 微信（YYW / ClawBot / FlowBot）、飞书、QQ、Telegram、钉钉、QQ 频道、Pagermaid、Web Bot
- **插件系统** — 支持 JavaScript / Python 插件，热重载，无需重启
- **Web 管理后台** — 基于 Ant Design Vue 的美观控制台，管理 Bot、插件、规则、存储等
- **消息规则** — 正则自动回复、群聊搬运、命令权限控制
- **跨平台** — 支持 Windows、macOS、Linux
- **前端内嵌** — 前端资源打包进 Go 二进制，无需单独部署 Web 服务

## 快速开始

### 环境要求

- Go 1.18+（从源码编译需要）
- Node.js 18+ 并启用 corepack（JS 插件需要，可选）
- Python 3（Python 插件需要，可选）

### 从源码编译

```bash
# 克隆
git clone https://github.com/kilimro/sillyGirl.git
cd sillyGirl

# 编译前端
cd frontend
npm install
npx vite build
cd ..

# 编译后端
go build -o mianpro .
```

### 运行

```bash
./mianpro
# Windows 下
MianPro.exe
```

浏览器打开 **http://localhost:8080/admin**。

## 部署方式

### Docker（推荐）

```bash
# 拉取镜像
docker pull kilimro/mianpro:latest

# 运行
docker run -d \
  --name mianpro \
  -p 8080:8080 \
  -v mianpro-data:/data \
  --restart unless-stopped \
  kilimro/mianpro:latest
```

### Linux 服务器

```bash
# 编译 Linux 版
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o mianpro .

# 后台运行
nohup ./mianpro > mianpro.log 2>&1 &
```

### macOS

```bash
go build -o mianpro .
./mianpro
```

### Windows

1. 从 [Releases](https://github.com/kilimro/sillyGirl/releases) 下载 `MianPro.exe`
2. 安装 [Node.js LTS](https://nodejs.org)（JS 插件必须）
3. 在命令提示符运行 `corepack enable`
4. 双击 `MianPro.exe`
5. 浏览器打开 http://localhost:8080/admin

### 从 macOS/Linux 交叉编译 Windows 版

```bash
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o MianPro.exe .
```

## 适配器列表

| 平台 | 连接方式 | 配置项 |
|---|---|---|
| YYW 微信 | Webhook 回调 | API 地址 + 鉴权密钥 |
| 飞书 | WebSocket 长连接 | App ID + App Secret |
| ClawBot | WebSocket | Token |
| FlowBot | HTTP | Token + 地址 + 端口 |
| QQ | WebSocket | App ID + App Secret |
| Telegram | 长轮询 | Bot Token |
| 钉钉 | Stream | Client ID + Client Secret |
| QQ 频道 | Webhook | App ID + App Secret |
| Pagermaid | Webhook | Token |
| Web Bot | 内置 | 默认开启 |

## 项目结构

```
├── adapters/          # 平台适配器
│   ├── feishu/        # 飞书机器人
│   ├── yyw/           # YYW 微信插件
│   ├── clawbot/       # ClawBot
│   └── ...
├── core/              # 核心框架（存储、路由、插件）
├── frontend/          # Vue 3 + Ant Design Vue 管理后台
├── main.go            # 入口文件
└── go.mod
```

## 相关文档

- [适配器指南](docs/adapters.md)
- [API 与存储](docs/api-storage.md)
- [插件开发](docs/plugin-development.md)
- [界面截图](docs/screenshots.md)

## 开源协议

MIT
