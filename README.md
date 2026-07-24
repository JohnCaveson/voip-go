# Gather

A peer-to-peer and hosted communication app built with Wails, React, and Go. Friends can gather, chat, and talk — or communities can spin up their own hosted instance.

## Quick Start

```bash
# 1. Install system deps (see below)
# 2. Install Go 1.26+, Node 18+, Wails CLI
# 3. Clone and run
git clone https://github.com/your-username/voip-go.git
cd voip-go
go mod tidy && cd frontend && npm install && cd ..
cd cmd/p2p
wails dev -tags webkit2_41   # <-- the build tag is required on most Linux distros
```

The app window opens with hot-reload. The signaling server starts on `:9321` and mDNS discovery begins immediately. You should see `"Discovered N peer(s)"` in the terminal if other instances are running on your LAN.

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

**Linux (Arch):**
```bash
sudo pacman -S --needed \
  webkit2gtk-4.1 \
  gtk3 \
  pkgconf
```

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

Verify — this is the most important pre-flight check:
```bash
wails doctor
```

`wails doctor` reports missing system libraries. The only critical one is **libwebkit**. Everything else is optional. If `wails doctor` says `libwebkit` is "Not Found" but you know it's installed (e.g. on Arch), that's fine — it sometimes can't detect it. The build tag `-tags webkit2_41` (below) forces the correct linking.

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
cd cmd/p2p
wails dev -tags webkit2_41
```

> **Important:** You must pass `-tags webkit2_41` on most Linux distros. Without it the build will fail with `Package webkit2gtk-4.0 was not found`. This tag tells Wails to link against `webkit2gtk-4.1` instead of the older `4.0`.

This starts the app with hot-reload. Changes to Go or React code auto-refresh. You'll see three things in the terminal:
1. **Vite** dev server on `http://localhost:5173` (frontend assets)
2. **Wails dev server** on a random port (browser access to Go bindings)
3. **The native app window** itself

If you edit Go code, the entire app recompiles (takes a few seconds). If you edit frontend code, Vite hot-reloads instantly.

**Hosted mode (requires MongoDB + signaling server):**
```bash
# Start backend services
docker compose up -d

# Run the app connected to them
cd cmd/hosted
VOIP_APP_MODE=hosted \
VOIP_MONGODB_URI=mongodb://localhost:27017 \
VOIP_SERVER_ADDR=ws://localhost:9321/signaling \
wails dev -tags webkit2_41
```

---

## Running Over a Local Network

Once built, anyone on your LAN can connect to the app. Here's how:

### 1. Build the App

```bash
cd cmd/p2p && wails build -tags webkit2_41 -o gather
```

### 2. Find Your Local IP

```bash
# Linux
ip addr show | grep "inet " | grep -v 127.0.0.1

# macOS
ipconfig getifaddr en0

# Windows
ipconfig
```

Note your IP (e.g., `192.168.1.50`).

### 3. Run the App

```bash
./build/bin/gather
```

The app starts a signaling server on port 9321 and advertises via mDNS.

### 4. Connect from Another Machine

On the other machine, either:
- **Run the same binary** — it will discover peers automatically via mDNS and show them in the sidebar under "Peers Nearby"
- **Or distribute the binary** — copy `build/bin/gather` to the other machine (same OS/architecture), run it, and it will find the first peer via mDNS

Both peers need to be on the same subnet for mDNS discovery to work.

### 5. Start Talking

1. Both users enter a username when prompted
2. Click on an audio room (e.g., "Lounge")
3. Click the join button ("Start Yapping", "Hop In", etc.)
4. Your microphone activates and WebRTC establishes a peer-to-peer audio connection
5. Text chat also works — messages are relayed through the signaling server

### Firewall Notes

If peers can't discover each other, allow these ports:

```bash
# Linux (ufw)
sudo ufw allow 9321/tcp    # Signaling server
sudo ufw allow 5353/udp    # mDNS discovery

# macOS — System Preferences > Security > Firewall
# Allow incoming connections for the Gather app
```

### Custom Username

Set your display name via environment variable instead of using the modal:

```bash
VOIP_USERNAME=YourName ./build/bin/gather
```

---

## Running Over the Internet (WAN)

> **Coming soon.** The hosted mode with Docker Compose is the foundation for this.

To let friends connect over the internet:

1. **Deploy the signaling server** to a VPS (DigitalOcean, Hetzner, etc.)
2. **Run MongoDB + signaling** via Docker Compose on the server
3. **Open port 9321** in the server's firewall
4. **Run the app** in hosted mode on each client:

```bash
VOIP_APP_MODE=hosted \
VOIP_SERVER_ADDR=ws://your-server-ip:9321/signaling \
VOIP_USERNAME=YourName \
./build/bin/gather-hosted
```

This requires the `gather-hosted` build (see [Building for Production](#building-for-production)).

For NAT traversal on peer-to-peer audio (when both clients are behind routers), a TURN server is needed. See the `VOIP_TURN_*` environment variables.

---

## Building for Production

### Build Commands

```bash
# P2P mode (from repo root or cmd/p2p)
cd cmd/p2p && wails build -tags webkit2_41 -o voip-p2p

# Hosted mode
cd cmd/hosted && wails build -tags webkit2_41 -o voip-hosted
```

The output is a single binary in `build/bin/`. No installer needed — just copy and run.

### Supported Platforms

| Target | Command | Notes |
|---|---|---|
| Linux (amd64) | `-platform linux/amd64` | Requires webkit2gtk on target machine |
| Linux (arm64) | `-platform linux/arm64` | For Raspberry Pi, ARM servers |
| macOS (amd64) | `-platform darwin/amd64` | Intel Macs |
| macOS (arm64) | `-platform darwin/arm64` | Apple Silicon (M1/M2/M3) |
| Windows (amd64) | `-platform windows/amd64` | Requires WebView2 (usually pre-installed) |

### Distributing the App

**Wails produces a single executable** — no installer, no framework dependencies on the target (except system libraries):

```bash
# Build for your platform
cd cmd/p2p && wails build -tags webkit2_41 -o gather

# The binary is here:
ls build/bin/gather

# Copy it anywhere — USB, shared folder, SCP, etc.
scp build/bin/gather friend@192.168.1.60:~/
```

**What your friend needs:**
- **Linux:** `webkit2gtk-4.1` installed (`sudo apt install libwebkit2gtk-4.1-0`)
- **macOS:** Nothing extra — the binary is self-contained
- **Windows:** WebView2 (comes with Windows 10/11)

That's it. They run the binary, enter a username, and they're in.

---

### Standalone Signaling Server

```bash
go build -tags webkit2_41 -o signaling-server ./server/cmd/server
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
│   │   ├── main.go           # Wails entry point — creates app, configures window
│   │   ├── app.go            # App logic — SQLite storage, mDNS discovery, WebRTC
│   │   ├── wails.json        # Wails config — points to ../../frontend
│   │   ├── go.mod
│   │   └── go.sum
│   └── hosted/               # Hosted desktop app
│       ├── main.go           # Wails entry point
│       ├── app.go            # App logic — MongoDB storage, remote signaling
│       ├── wails.json        # Wails config
│       ├── go.mod
│       └── go.sum
├── internal/
│   ├── config/               # Environment-based configuration (VOIP_* vars)
│   ├── storage/              # Storage backends
│   │   ├── interface.go      # Storage interface (Channels, Messages, Users, Settings)
│   │   ├── sqlite.go         # SQLite implementation (P2P mode)
│   │   └── mongodb.go        # MongoDB implementation (hosted mode)
│   ├── channel/              # Room management — create, join, list channels
│   ├── signaling/            # WebSocket signaling (server + client)
│   └── discovery/            # mDNS LAN peer discovery
├── pkg/
│   ├── api/                  # Shared signaling message types
│   └── models/               # Domain models (User, Channel, Message)
├── server/                   # Standalone signaling server (hosted mode)
│   └── cmd/server/main.go
├── frontend/                 # React + TypeScript UI
│   ├── src/
│   │   ├── App.tsx           # Main app component — mounts layout, modals
│   │   ├── App.css           # Warm earth-tone theme
│   │   ├── store/
│   │   │   └── layoutStore.ts    # Zustand layout state — panel positions/sizes
│   │   ├── utils/
│   │   │   └── layoutCodec.ts    # Layout encode/decode — base64 persistence
│   │   ├── components/       # UI components
│   │   │   ├── Panel.tsx         # Draggable/resizable panel (react-rnd)
│   │   │   ├── Layout.tsx        # Panel layout container
│   │   │   ├── SnapToggle.tsx    # Grid/free-form toggle
│   │   │   ├── TextChannel.tsx   # Text chat panel
│   │   │   ├── VoiceChannel.tsx  # Voice chat panel
│   │   │   ├── Sidebar.tsx       # Channel list, peers, user info panels
│   │   │   ├── ConnectionStatus.tsx
│   │   │   └── ... (modals)
│   │   ├── hooks/            # useSignaling, useWebRTC
│   │   └── services/         # SignalingClient, StorageService
│   ├── wailsjs/              # Auto-generated Wails Go<->JS bindings
│   ├── package.json
│   ├── vite.config.ts
│   └── tsconfig.json
├── go.work                   # Go workspace — links cmd/, internal/, pkg/, server/
├── docker-compose.yml        # MongoDB + signaling server
├── Dockerfile                # Signaling server container
├── README.md
└── PLAN.md
```

---

## Troubleshooting

### Build Errors

**`Package webkit2gtk-4.0 was not found` / `libwebkit not found`:**
Your system has webkit2gtk-4.1, not the older 4.0. Always pass the build tag:
```bash
wails dev -tags webkit2_41
wails build -tags webkit2_41 -o gather
```
The wails.json files already include `"buildTags": "webkit2_41"` but the CLI flag takes precedence and is more reliable.

**`wails dev` fails with "cannot find wails.json":**
You must run from `cmd/p2p/` or `cmd/hosted/` — those directories have their own `wails.json`:
```bash
cd cmd/p2p && wails dev -tags webkit2_41
```

**Go build fails with "cannot find module":**
This is a Go workspace issue. Run from the repo root:
```bash
go mod tidy          # root module
cd cmd/p2p && go mod tidy   # p2p module
cd ../../frontend && npm install
```

**`wails dev` fails with "webview not found":**
Run `wails doctor` and install any missing system dependencies. See the [Install System Dependencies](#1-install-system-dependencies) section.

### Runtime Errors

**Port 9321 already in use:**
Either stop the other process or use a different port:
```bash
VOIP_PORT=9322 wails dev -tags webkit2_41
```

**MongoDB connection refused (hosted mode):**
Make sure Docker is running and the services are up:
```bash
docker compose ps
docker compose up -d
```

**Frontend shows blank screen:**
```bash
cd frontend && npm install && npm run build
```

**App starts but no peers are discovered:**
mDNS only works on the same subnet. Check that:
1. Both machines are on the same network
2. UDP port 5353 is not blocked by your firewall
3. Your router/access point allows multicast traffic

### Where Things Live

| What | Location |
|---|---|
| Local data (P2P mode) | `./data/` directory (SQLite database, settings) |
| Build output | `cmd/p2p/build/bin/` or `cmd/hosted/build/bin/` |
| Wails dev server | Random port shown in terminal (e.g. `http://localhost:34115`) |
| Vite dev server | `http://localhost:5173` |
| Signaling server | Port `9321` (both dev and production) |
| mDNS broadcast | UDP port `5353` |

### Debugging Tips

**View Go logs in dev mode:**
All `println()` and log output from the Go backend appears in the terminal where you ran `wails dev`.

**Frontend dev tools:**
In dev mode, right-click the app window → "Inspect Element" to open WebKit Inspector. This works exactly like Chrome DevTools.

**Inspect Wails bindings:**
In the browser at the Wails dev server URL (e.g. `http://localhost:34115`), you can call Go methods from the console:
```js
// Example: list channels
window.go.main.App.GetChannels()
```

**Reset local data:**
Delete the data directory to start fresh:
```bash
rm -rf cmd/p2p/data/
```

**Test the signaling server standalone:**
```bash
cd server && go build -o signaling-server ./cmd/server
./signaling-server --port 9321
# Then point a client at it:
VOIP_APP_MODE=hosted VOIP_SERVER_ADDR=ws://localhost:9321/signaling wails dev -tags webkit2_41
```
