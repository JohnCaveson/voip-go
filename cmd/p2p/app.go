package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/voip-app/internal/channel"
	"github.com/voip-app/internal/config"
	"github.com/voip-app/internal/discovery"
	"github.com/voip-app/internal/signaling"
	"github.com/voip-app/internal/storage"
)

type Peer struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Addr          string `json:"addr"`
	SignalingAddr string `json:"signaling_addr"`
}

type App struct {
	ctx          context.Context
	cfg          config.Config
	storage      storage.Storage
	channelMgr   *channel.Manager
	discoverer   *discovery.Discoverer
	hub          *signaling.Hub
	httpServer   *http.Server
	signalingURL string
	mu           sync.RWMutex
	peers        map[string]Peer
}

func NewApp() *App {
	return &App{
		peers: make(map[string]Peer),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.cfg = config.Load()
	a.cfg.AppMode = config.AppModeP2P
	a.cfg.StorageType = config.StorageTypeSQLite

	s, err := a.initStorage()
	if err != nil {
		log.Printf("Failed to initialize storage: %v", err)
		return
	}
	a.storage = s

	a.channelMgr = channel.NewManager(s)
	if err := a.channelMgr.Init(ctx); err != nil {
		log.Printf("Failed to init channels: %v", err)
	}

	a.startSignalingServer()

	if a.cfg.NetworkMode == config.NetworkModeLAN {
		d, err := discovery.NewDiscoverer(a.cfg.Username, a.cfg.Port, a.signalingURL)
		if err != nil {
			log.Printf("Failed to start discovery: %v", err)
		} else {
			a.discoverer = d
			go a.discoverPeers()
		}
	}
}

func (a *App) shutdown(ctx context.Context) {
	if a.httpServer != nil {
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("Signaling server shutdown error: %v", err)
		}
	}
	if a.discoverer != nil {
		a.discoverer.Close()
	}
	if a.storage != nil {
		a.storage.Close()
	}
}

func (a *App) startSignalingServer() {
	a.hub = signaling.NewHub()
	mux := http.NewServeMux()
	mux.Handle("/signaling", a.hub)

	addr := fmt.Sprintf(":%d", a.cfg.Port)
	a.httpServer = &http.Server{Addr: addr, Handler: mux}
	a.signalingURL = fmt.Sprintf("ws://localhost:%d/signaling", a.cfg.Port)

	go func() {
		log.Printf("Signaling server starting on %s", addr)
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Signaling server error: %v", err)
		}
	}()
}

func (a *App) initStorage() (storage.Storage, error) {
	dbDir := a.cfg.DataDir
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, err
	}
	return storage.NewSQLiteStorage(filepath.Join(dbDir, "voip.db"))
}

func (a *App) discoverPeers() {
	for {
		peers, err := a.discoverer.DiscoverWithTimeout(10 * time.Second)
		if err != nil {
			log.Printf("Discovery error: %v", err)
			continue
		}

		a.mu.Lock()
		for _, p := range peers {
			if p.Username == a.cfg.Username {
				continue
			}
			a.peers[p.ID] = Peer{
				ID:            p.ID,
				Username:      p.Username,
				Addr:          p.Addr.String(),
				SignalingAddr: p.SignalingAddr,
			}
		}
		a.mu.Unlock()

		if len(peers) > 0 {
			log.Printf("Discovered %d peer(s)", len(peers))
		}
	}
}

func (a *App) GetChannels() []channel.ChannelInfo {
	channels, err := a.channelMgr.List(a.ctx)
	if err != nil {
		return nil
	}

	var result []channel.ChannelInfo
	for _, ch := range channels {
		result = append(result, channel.ChannelInfo{
			ID:        ch.ID,
			Name:      ch.Name,
			Type:      string(ch.Type),
			IsDefault: ch.IsDefault,
		})
	}
	return result
}

func (a *App) CreateChannel(name, chType string) error {
	_, err := a.channelMgr.Create(a.ctx, name, channel.ParseType(chType))
	return err
}

func (a *App) DeleteChannel(id string) error {
	return a.channelMgr.Delete(a.ctx, id)
}

func (a *App) RenameChannel(id, newName string) error {
	return a.channelMgr.Rename(a.ctx, id, newName)
}

func (a *App) GetConfig() config.Config {
	return a.cfg
}

func (a *App) GetDiscoveredPeers() []Peer {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var result []Peer
	for _, p := range a.peers {
		result = append(result, p)
	}
	return result
}

func (a *App) GetSignalingURL() string {
	return a.signalingURL
}

func (a *App) SaveLayout(layoutData string) error {
	return a.storage.SetSetting(a.ctx, "layout", layoutData)
}

func (a *App) GetLayout() string {
	data, err := a.storage.GetSetting(a.ctx, "layout")
	if err != nil {
		return ""
	}
	return data
}
