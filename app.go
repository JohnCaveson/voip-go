package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/voip-app/internal/channel"
	"github.com/voip-app/internal/config"
	"github.com/voip-app/internal/discovery"
	"github.com/voip-app/internal/storage"

	// Register MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

type App struct {
	ctx        context.Context
	cfg        config.Config
	storage    storage.Storage
	channelMgr *channel.Manager
	discoverer *discovery.Discoverer
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.cfg = config.Load()

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

	if a.cfg.NetworkMode == config.NetworkModeLAN {
		d, err := discovery.NewDiscoverer(a.cfg.Username, a.cfg.Port)
		if err != nil {
			log.Printf("Failed to start discovery: %v", err)
		} else {
			a.discoverer = d
			go a.discoverPeers()
		}
	}
}

func (a *App) initStorage() (storage.Storage, error) {
	switch a.cfg.StorageType {
	case config.StorageTypeMySQL:
		return storage.NewMySQLStorage(a.cfg.MySQLDSN)
	default:
		dbDir := a.cfg.DataDir
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return nil, err
		}
		return storage.NewSQLiteStorage(filepath.Join(dbDir, "voip.db"))
	}
}

func (a *App) discoverPeers() {
	for {
		peers, err := a.discoverer.DiscoverWithTimeout(10 * time.Second)
		if err != nil {
			log.Printf("Discovery error: %v", err)
			continue
		}
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

func (a *App) Shutdown() {
	if a.discoverer != nil {
		a.discoverer.Close()
	}
	if a.storage != nil {
		a.storage.Close()
	}
}
