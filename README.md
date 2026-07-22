# VoIP Desktop Application

A peer-to-peer and hosted VoIP application built with Wails, React, and Go.

## Two App Modes

| | **P2P** | **Hosted** |
|---|---|---|
| Storage | SQLite (local file) | MongoDB (centralized) |
| Signaling | Embedded WebSocket | Standalone server |
| Discovery | mDNS (LAN peers) | Server-managed |
| Privacy | All data stays on device | Data stored centrally |
| Deployment | Desktop binary | Docker Compose |

## Prerequisites

- Go 1.26+
- Node.js 18+
- Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- Docker & Docker Compose (for hosted mode)

## Quick Start

```bash
# Install dependencies
cd frontend && npm install && cd ..
go mod tidy

# Development (P2P mode, default)
wails dev -tags webkit2_41
```

## Building

### P2P App (Privacy-focused)

```bash
# Development
wails dev -tags webkit2_41

# Production build
wails build -tags webkit2_41 -o voip-p2p

# Run production binary
./build/bin/voip-p2p
```

### Hosted App (Centralized)

```bash
# Start MongoDB + signaling server
docker compose up -d

# Development (connects to Docker services)
VOIP_APP_MODE=hosted \
VOIP_MONGODB_URI=mongodb://localhost:27017 \
VOIP_SERVER_ADDR=ws://localhost:9321/signaling \
wails dev -tags webkit2_41

# Production build
wails build -tags webkit2_41 -o voip-hosted

# Run production binary
VOIP_APP_MODE=hosted \
VOIP_MONGODB_URI=mongodb://localhost:27017 \
VOIP_SERVER_ADDR=ws://localhost:9321/signaling \
./build/bin/voip-hosted
```

### Standalone Signaling Server

```bash
# Run directly
go run ./server/cmd/server --port 9321

# Or build and run
go build -o signaling-server ./server/cmd/server
./signaling-server --port 9321
```

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `VOIP_APP_MODE` | `p2p` | App mode: `p2p` or `hosted` |
| `VOIP_NETWORK_MODE` | `lan` | Network mode: `lan` or `wan` |
| `VOIP_PORT` | `9321` | Signaling server port |
| `VOIP_DATA_DIR` | `./data` | Local data directory (P2P) |
| `VOIP_SERVER_ADDR` | | Signaling server address (e.g., `ws://host:9321/signaling`) |
| `VOIP_USERNAME` | `anonymous` | Display username |
| `VOIP_MONGODB_URI` | `mongodb://localhost:27017` | MongoDB connection URI (hosted) |
| `VOIP_STUN_URLS` | Google STUN | Comma-separated STUN server URLs |
| `VOIP_TURN_URL` | | TURN server URL |
| `VOIP_TURN_USERNAME` | | TURN username |
| `VOIP_TURN_PASSWORD` | | TURN password |

## Running Tests

### Run All Tests

```bash
go test ./... -v
```

### Run Specific Test Suites

```bash
# Config tests
go test ./internal/config/... -v

# Storage tests (SQLite)
go test ./internal/storage/... -v -run "Test(Create|Get|List|Update|Delete|Send|User)"

# Channel manager tests
go test ./internal/channel/... -v

# Signaling server tests
go test ./internal/signaling/... -v

# API type tests
go test ./pkg/api/... -v

# Discovery tests
go test ./internal/discovery/... -v
```

### MongoDB Integration Tests

MongoDB integration tests require a running MongoDB instance:

```bash
# Start MongoDB
docker compose up -d mongodb

# Run MongoDB tests
MONGODB_URI=mongodb://localhost:27017 go test ./internal/storage/... -v -run "TestMongoDB"
```

### Run Tests with Coverage

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Project Structure

```
voip-go/
├── cmd/
│   ├── p2p/              # P2P desktop app entry point
│   │   ├── main.go
│   │   └── app.go
│   └── hosted/           # Hosted desktop app entry point
│       ├── main.go
│       └── app.go
├── internal/
│   ├── config/           # App configuration
│   ├── models/           # Data models
│   ├── storage/          # Storage backends (SQLite, MongoDB)
│   ├── channel/          # Channel management
│   ├── signaling/        # WebSocket signaling
│   └── discovery/        # mDNS peer discovery
├── pkg/api/              # Shared API types
├── server/               # Standalone signaling server
├── frontend/             # React + TypeScript UI
├── docker-compose.yml    # MongoDB + signaling server
├── wails-p2p.json        # Wails config for P2P build
└── wails-hosted.json     # Wails config for hosted build
```

## Docker Compose (Hosted Mode)

```bash
# Start all services
docker compose up -d

# View logs
docker compose logs -f

# Stop all services
docker compose down

# Stop and remove data
docker compose down -v
```

Services:
- **mongodb**: MongoDB 7 on port 27017
- **signaling**: WebSocket signaling server on port 9321
