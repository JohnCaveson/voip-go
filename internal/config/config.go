package config

import (
	"os"
	"strconv"
)

type NetworkMode string

const (
	NetworkModeLAN NetworkMode = "lan"
	NetworkModeWAN NetworkMode = "wan"
)

type AppMode string

const (
	AppModeP2P    AppMode = "p2p"
	AppModeHosted AppMode = "hosted"
)

type StorageType string

const (
	StorageTypeSQLite  StorageType = "sqlite"
	StorageTypeMongoDB StorageType = "mongodb"
)

type Config struct {
	NetworkMode NetworkMode
	AppMode     AppMode
	Port        int
	DataDir     string
	STUNURLs    []string
	TURNConfig  TURNConfig
	ServerAddr  string
	Username    string
	StorageType StorageType
	MongoDBURI  string
}

type TURNConfig struct {
	URL      string
	Username string
	Password string
}

func DefaultConfig() Config {
	return Config{
		NetworkMode: NetworkModeLAN,
		AppMode:     AppModeP2P,
		Port:        9321,
		DataDir:     "./data",
		STUNURLs: []string{
			"stun:stun.l.google.com:19302",
			"stun:stun1.l.google.com:19302",
		},
		TURNConfig: TURNConfig{},
		Username:   "anonymous",
		StorageType: StorageTypeSQLite,
		MongoDBURI:  "mongodb://localhost:27017",
	}
}

func Load() Config {
	cfg := DefaultConfig()

	if v := os.Getenv("VOIP_APP_MODE"); v == string(AppModeHosted) {
		cfg.AppMode = AppModeHosted
	}
	if v := os.Getenv("VOIP_NETWORK_MODE"); v == string(NetworkModeWAN) {
		cfg.NetworkMode = NetworkModeWAN
	}
	if v := os.Getenv("VOIP_STORAGE_TYPE"); v == string(StorageTypeMongoDB) {
		cfg.StorageType = StorageTypeMongoDB
	}
	if v := os.Getenv("VOIP_MONGODB_URI"); v != "" {
		cfg.MongoDBURI = v
	}
	if v := os.Getenv("VOIP_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p > 0 && p < 65536 {
			cfg.Port = p
		}
	}
	if v := os.Getenv("VOIP_DATA_DIR"); v != "" {
		cfg.DataDir = v
	}
	if v := os.Getenv("VOIP_SERVER_ADDR"); v != "" {
		cfg.ServerAddr = v
	}
	if v := os.Getenv("VOIP_USERNAME"); v != "" {
		cfg.Username = v
	}
	if v := os.Getenv("VOIP_TURN_URL"); v != "" {
		cfg.TURNConfig.URL = v
	}
	if v := os.Getenv("VOIP_TURN_USERNAME"); v != "" {
		cfg.TURNConfig.Username = v
	}
	if v := os.Getenv("VOIP_TURN_PASSWORD"); v != "" {
		cfg.TURNConfig.Password = v
	}

	return cfg
}
