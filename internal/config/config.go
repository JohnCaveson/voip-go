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

type StorageType string

const (
	StorageTypeSQLite StorageType = "sqlite"
	StorageTypeMySQL  StorageType = "mysql"
)

type Config struct {
	NetworkMode NetworkMode
	Port        int
	DataDir     string
	STUNURLs    []string
	TURNConfig  TURNConfig
	ServerAddr  string
	Username    string
	StorageType StorageType
	MySQLDSN    string
}

type TURNConfig struct {
	URL      string
	Username string
	Password string
}

func DefaultConfig() Config {
	return Config{
		NetworkMode: NetworkModeLAN,
		Port:        9321,
		DataDir:     "./data",
		STUNURLs: []string{
			"stun:stun.l.google.com:19302",
			"stun:stun1.l.google.com:19302",
		},
		TURNConfig: TURNConfig{},
		Username:   "anonymous",
		StorageType: StorageTypeSQLite,
		MySQLDSN:    "root:password@tcp(127.0.0.1:3306)/voip?parseTime=true",
	}
}

func Load() Config {
	cfg := DefaultConfig()

	if v := os.Getenv("VOIP_NETWORK_MODE"); v == string(NetworkModeWAN) {
		cfg.NetworkMode = NetworkModeWAN
	}
	if v := os.Getenv("VOIP_STORAGE_TYPE"); v == string(StorageTypeMySQL) {
		cfg.StorageType = StorageTypeMySQL
	}
	if v := os.Getenv("VOIP_MYSQL_DSN"); v != "" {
		cfg.MySQLDSN = v
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
