package main

import (
	"context"
	"log"

	"github.com/voip-app/internal/channel"
	"github.com/voip-app/internal/config"
	"github.com/voip-app/internal/signaling"
	"github.com/voip-app/internal/storage"
)

type App struct {
	ctx            context.Context
	cfg            config.Config
	storage        storage.Storage
	channelMgr     *channel.Manager
	signalingClient *signaling.SignalingClient
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.cfg = config.Load()
	a.cfg.AppMode = config.AppModeHosted
	a.cfg.StorageType = config.StorageTypeMongoDB

	if a.cfg.MongoDBURI == "" {
		log.Fatal("VOIP_MONGODB_URI is required for hosted mode")
	}

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

	if a.cfg.ServerAddr != "" {
		go a.connectSignaling()
	}
}

func (a *App) shutdown(ctx context.Context) {
	if a.signalingClient != nil {
		a.signalingClient.Close()
	}
	if a.storage != nil {
		a.storage.Close()
	}
}

func (a *App) initStorage() (storage.Storage, error) {
	return storage.NewMongoDBStorage(a.cfg.MongoDBURI)
}

func (a *App) connectSignaling() {
	client, err := signaling.NewSignalingClient(a.cfg.ServerAddr)
	if err != nil {
		log.Printf("Failed to connect to signaling server: %v", err)
		return
	}
	a.signalingClient = client
	log.Printf("Connected to signaling server at %s", a.cfg.ServerAddr)
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
