# VoIP Desktop Application - Build Plan

## Tech Stack

| Area             | Choice                                    |
| ---------------- | ----------------------------------------- |
| GUI Framework    | Wails v2 (stable)                         |
| Frontend         | React 19 + TypeScript                     |
| Frontend WebRTC  | Browser WebRTC APIs (WebView)             |
| Panel Layout     | react-rnd (drag/resize)                   |
| State Management | Zustand                                   |
| Server WebRTC    | Pion WebRTC v4 (SFU/TURN)                |
| Signaling        | WebSocket (`gorilla/websocket`)           |
| Storage (P2P)    | SQLite (`modernc.org/sqlite`)             |
| Storage (Hosted) | MongoDB (`go.mongodb.org/mongo-driver`)   |
| LAN Discovery    | mDNS (`hashicorp/mdns`)                   |
| Deployment       | Docker Compose (hosted mode)              |
| Project Structure| Multi-module Go workspace                 |

## App Modes

| | P2P | Hosted |
|---|---|---|
| Storage | SQLite (local) | MongoDB (centralized) |
| Signaling | Embedded | Standalone server |
| Discovery | mDNS | Server-managed |
| Privacy | Device-only | Central server |

## Default Channels

| Name         | Type  | Deletable |
| ------------ | ----- | --------- |
| `#general`   | Text  | No        |
| `🔊 General` | Voice | No        |

## Phase 1: Foundation
- Go workspace & modules with deps
- Models: User, Channel, Message
- Storage interface + SQLite implementation + tests
- Channel manager with defaults + tests
- Config + tests
- Wails frontend scaffold

## Phase 2: Signaling & Discovery
- WebSocket signaling server (hub, relay)
- Signaling client wrapper
- mDNS LAN discovery

## Phase 3: VoIP & Screen Sharing
- Frontend WebRTC hooks
- VoiceChannel component
- ScreenShare component

## Phase 4: Text Channels & UI
- TextChannel component
- ChannelList sidebar
- App layout + settings modal

## Phase 5: Server Binary & Polish
- WAN server binary
- Client binary
- Full test suite

## Phase 6: Dual-App Architecture (P2P + Hosted)
- Config: AppMode + MongoDBURI fields
- MongoDB storage backend + tests
- Two build targets: `cmd/p2p/` and `cmd/hosted/`
- Docker Compose with MongoDB + signaling server
- Frontend mode awareness (badges, connection status)
- Removed MySQL storage (replaced by MongoDB)
- Shared frontend embed package
- Updated README with build/run/test procedures

## Phase 7: Customizable Panel UI
- Panel-based layout system (every UI element is a draggable, resizable panel)
- `react-rnd` for drag/resize, `zustand` for state management
- Snap-to-grid mode with toggle (grid/free-form)
- Panel chrome: header, minimize/maximize/close controls
- Default layout mimics original sidebar+content app
- Non-closable panels: header, channel list, first channel
- Layout persistence: encoded base64 string via Go backend (primary), localStorage (fallback)
- Backend: `settings` table (SQLite) / `settings` collection (MongoDB) with `GetSetting`/`SetSetting`
- Wails bindings: `SaveLayout`/`GetLayout` methods on both P2P and Hosted app structs
- Panel splitting: duplicate channel panels
