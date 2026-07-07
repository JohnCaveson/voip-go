# VoIP Desktop Application - Build Plan

## Tech Stack

| Area             | Choice                          |
| ---------------- | ------------------------------- |
| GUI Framework    | Wails v2 (stable)               |
| Frontend         | React + TypeScript + Tailwind   |
| Frontend WebRTC  | Browser WebRTC APIs (WebView)   |
| Server WebRTC    | Pion WebRTC v4 (SFU/TURN)      |
| Signaling        | WebSocket (`gorilla/websocket`) |
| Storage          | SQLite (`modernc.org/sqlite`)   |
| LAN Discovery    | mDNS (`hashicorp/mdns`)         |
| Project Structure| Multi-module Go workspace       |

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
