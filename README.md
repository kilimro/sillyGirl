# MianPro

[![Version](https://img.shields.io/badge/version-v1.1.9-1677ff)](VERSION)
[![Go](https://img.shields.io/badge/Go-1.18+-00ADD8?logo=go)](go.mod)
[![Docker](https://img.shields.io/badge/Docker-pull-2496ED?logo=docker)](https://hub.docker.com/r/kilimro/mianpro)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

English | [中文](README_ZH.md)

A multi-platform chat bot framework built on Go. Connect WeChat, Feishu, QQ, Telegram, DingTalk and more through a single unified admin panel.

## Features

- **Multi-platform adapters** — WeChat (YYW / ClawBot / FlowBot), Feishu, QQ, Telegram, DingTalk, QQ Guild, Pagermaid, Web Bot
- **Plugin system** — Extend with JavaScript or Python plugins, hot-reloadable, no restart needed
- **Web admin panel** — Beautiful Ant Design Vue dashboard for bots, plugins, rules, storage and more
- **Message rules** — Pattern-based auto-reply, group forwarding, command permissions
- **Cross-platform** — Runs on Windows, macOS and Linux
- **Embedded frontend** — Frontend assets bundled into the Go binary, no separate web server needed

## Quick Start

### Prerequisites

- Go 1.18+ (for building from source)
- Node.js 18+ with corepack enabled (for JS plugins, optional)
- Python 3 (for Python plugins, optional)

### Build from Source

```bash
# Clone
git clone https://github.com/kilimro/sillyGirl.git
cd sillyGirl

# Build frontend
cd frontend
npm install
npx vite build
cd ..

# Build binary
go build -o mianpro .
```

### Run

```bash
./mianpro
# on Windows
MianPro.exe
```

Open **http://localhost:8080/admin** in your browser.

## Deployment

### Docker (Recommended)

```bash
# Pull image
docker pull kilimro/mianpro:latest

# Run
docker run -d \
  --name mianpro \
  -p 8080:8080 \
  -v mianpro-data:/data \
  --restart unless-stopped \
  kilimro/mianpro:latest
```

### Linux Server

```bash
# Build for Linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o mianpro .

# Run in background
nohup ./mianpro > mianpro.log 2>&1 &
```

### macOS

```bash
go build -o mianpro .
./mianpro
```

### Windows

1. Download `MianPro.exe` from [Releases](https://github.com/kilimro/sillyGirl/releases)
2. Install [Node.js LTS](https://nodejs.org) (required for JS plugins)
3. Run `corepack enable` in Command Prompt
4. Double-click `MianPro.exe`
5. Open http://localhost:8080/admin

### Cross-compile for Windows from macOS/Linux

```bash
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o MianPro.exe .
```

## Adapters

| Platform | Connection | Config |
|---|---|---|
| YYW WeChat | Webhook callback | API base URL + auth key |
| Feishu (Lark) | WebSocket SDK | App ID + App Secret |
| ClawBot | WebSocket | Token |
| FlowBot | HTTP | Token + host + port |
| QQ | WebSocket | App ID + App Secret |
| Telegram | Long polling | Bot token |
| DingTalk | Stream | Client ID + Client Secret |
| QQ Guild | Webhook | App ID + App Secret |
| Pagermaid | Webhook | Token |
| Web Bot | Built-in | Always on |

## Project Structure

```
├── adapters/          # Platform adapters
│   ├── feishu/        # Feishu (Lark) bot
│   ├── yyw/           # YYW WeChat plugin
│   ├── clawbot/       # ClawBot
│   └── ...
├── core/              # Core framework (storage, routing, plugins)
├── frontend/          # Vue 3 + Ant Design Vue admin panel
├── main.go            # Entry point
└── go.mod
```

## Documentation

- [Adapter Guide](docs/adapters.md)
- [API & Storage](docs/api-storage.md)
- [Plugin Development](docs/plugin-development.md)
- [Screenshots](docs/screenshots.md)

## License

MIT
