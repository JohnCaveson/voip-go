# Gather

A peer-to-peer and hosted communication app built with Wails, React, and Go. Friends can gather, chat, and talk — or communities can spin up their own hosted instance.

## Two App Modes

| | **P2P** | **Hosted** |
|---|---|---|
| Storage | SQLite (local file) | MongoDB (centralized) |
| Signaling | Embedded WebSocket | Standalone server |
| Discovery | mDNS (LAN peers) | Server-managed |
| Privacy | All data stays on device | Data stored centrally |
| Deployment | Desktop binary | Docker Compose |

---

## Getting Started

### 1. Install System Dependencies

**Linux (Debian/Ubuntu):**
```bash
sudo apt update
sudo apt install -y \
  build-essential \
  libgtk-3-dev \
  libwebkit2gtk-4.1-dev \
  libayatana-appindicator3-dev \
  librsvg2-dev
```

**Linux (Fedora):**
```bash
sudo dnf install -y \
  gcc \
  gtk3-devel \
  webkit2gtk4.1-devel \
  libappindicator-gtk3-devel \
  librsvg2-devel
```

**macOS:**
```bash
xcode-select --install
```

**Windows:**
- Install [WebView2](https://developer.microsoft.com/en-us/microsoft-edge/webview2/) (usually pre-installed on Windows 10/11)
- Install [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) or MinGW for CGO

### 2. Install Go

Download and install Go 1.26+ from [go.dev](https://go.dev/dl/).

Verify:
```bash
go version
```

### 3. Install Node.js

Download and install Node.js 18+ from [nodejs.org](https://nodejs.org/).

Verify:
```bash
node --version
npm --version
```

### 4. Install Wails CLI

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

Verify:
```bash
wails doctor
```

`wails doctor` will check your system has all required dependencies. Fix any issues it reports before continuing.

### 5. Clone and Setup

```bash
git clone https://github.com/your-username/voip-go.git
cd voip-go

# Install Go dependencies
go mod tidy

# Install frontend dependencies
cd frontend
npm install
cd ..
```

### 6. Run in Development

**P2P mode (default — everything runs locally):**
```bash
wails dev -tags webkit2_41
```

This starts the app with hot-reload. Changes to Go or React code auto-refresh.

**Hosted mode (requires MongoDB + signaling server):**
```bash
# Start backend services
docker compose up -d

# Run the app connected to them
VOIP_APP_MODE=hosted \
VOIP_MONGODB_URI=mongodb://localhost:27017 \
VOIP_SERVER_ADDR=ws://localhost:9321/signaling \
wails dev -tags webkit2_41
```

---

## Building for Production

### P2P App

```bash
wails build -tags webkit2_41 -o gather-p2p
./build/bin/gather-p2p
```

### Hosted App

```bash
wails build -tags webkit2_41 -o gather-hosted

# Start backend services first
docker compose up -d

# Run the app
VOIP_APP_MODE=hosted \
VOIP_MONGODB_URI=mongodb://localhost:27017 \
VOIP_SERVER_ADDR=ws://localhost:9321/signaling \
./build/bin/gather-hosted
```

### Standalone Signaling Server

```bash
go build -o signaling-server ./server/cmd/server
./signaling-server --port 9321
```

---

## Frontend Development (without Wails)

If you want to work on the frontend UI only without the Go backend:

```bash
cd frontend
npm run dev
```

This starts Vite on `http://localhost:5173`. The UI will run with placeholder data since Wails bindings aren't available outside the app.

---

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `VOIP_APP_MODE` | `p2p` | App mode: `p2p` or `hosted` |
| `VOIP_NETWORK_MODE` | `lan` | Network mode: `lan` or `wan` |
| `VOIP_PORT` | `9321` | Signaling server port |
| `VOIP_DATA_DIR` | `./data` | Local data directory (P2P) |
| `VOIP_SERVER_ADDR` | | Signaling server WebSocket URL |
| `VOIP_USERNAME` | `anonymous` | Display username |
| `VOIP_MONGODB_URI` | `mongodb://localhost:27017` | MongoDB connection URI |
| `VOIP_TURN_URL` | | TURN server URL |
| `VOIP_TURN_USERNAME` | | TURN username |
| `VOIP_TURN_PASSWORD` | | TURN password |

---

## Running Tests

### All Tests

```bash
go test ./... -v
```

### By Package

```bash
go test ./internal/config/... -v        # Config loading
go test ./internal/channel/... -v       # Room management
go test ./internal/signaling/... -v     # WebSocket signaling
go test ./internal/storage/... -v       # SQLite storage
go test ./pkg/api/... -v               # API types
go test ./internal/discovery/... -v     # mDNS discovery
```

### MongoDB Integration Tests

Requires a running MongoDB instance:

```bash
docker compose up -d mongodb
MONGODB_URI=mongodb://localhost:27017 go test ./internal/storage/... -v -run "TestMongoDB"
```

### Test Coverage

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## Docker Compose (Hosted Mode)

```bash
docker compose up -d          # Start MongoDB + signaling server
docker compose logs -f        # View logs
docker compose down           # Stop services
docker compose down -v        # Stop and delete data
```

Services:
- **mongodb** — MongoDB 7 on port 27017 (data persisted in `mongo_data` volume)
- **signaling** — WebSocket signaling server on port 9321

---

## Project Structure

```
voip-go/
├── cmd/
│   ├── p2p/                  # P2P desktop app
│   │   ├── main.go           # Wails entry point
│   │   └── app.go            # App logic (SQLite + mDNS)
│   └── hosted/               # Hosted desktop app
│       ├── main.go           # Wails entry point
│       └── app.go            # App logic (MongoDB + signaling)
├── internal/
│   ├── config/               # Environment-based configuration
│   ├── models/               # Data models (Room, Message, User)
│   ├── storage/              # Storage backends
│   │   ├── interface.go      # Storage interface
│   │   ├── sqlite.go         # SQLite implementation
│   │   └── mongodb.go        # MongoDB implementation
│   ├── channel/              # Room management
│   ├── signaling/            # WebSocket signaling (server + client)
│   └── discovery/            # mDNS LAN peer discovery
├── pkg/api/                  # Shared API types
├── server/                   # Standalone signaling server
│   └── cmd/server/main.go
├── frontend/                 # React + TypeScript UI
│   ├── src/
│   │   ├── App.tsx           # Main app component
│   │   ├── App.css           # Warm earth-tone theme
│   │   └── components/       # UI components
│   └── wailsjs/              # Auto-generated Wails bindings
├── docker-compose.yml        # MongoDB + signaling server
├── wails-p2p.json            # Wails config for P2P build
└── wails-hosted.json         # Wails config for hosted build
```

---

## Troubleshooting

**`wails dev` fails with "webview not found":**
Run `wails doctor` and install any missing system dependencies.

**Frontend shows blank screen:**
```bash
cd frontend && npm install && npm run build
```

**Go build fails with "cannot find module":**
```bash
go mod tidy
```

**MongoDB connection refused (hosted mode):**
Make sure Docker is running and the services are up:
```bash
docker compose ps
docker compose up -d
```

**Port 9321 already in use:**
Either stop the other process or use a different port:
```bash
VOIP_PORT=9322 wails dev -tags webkit2_41
```
