package config

import (
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.NetworkMode != NetworkModeLAN {
		t.Errorf("expected LAN mode, got %s", cfg.NetworkMode)
	}
	if cfg.Port != 9321 {
		t.Errorf("expected port 9321, got %d", cfg.Port)
	}
	if len(cfg.STUNURLs) == 0 {
		t.Error("expected at least one STUN URL")
	}
	if cfg.Username != "anonymous" {
		t.Errorf("expected username anonymous, got %s", cfg.Username)
	}
}

func TestLoadFromEnv(t *testing.T) {
	os.Setenv("VOIP_NETWORK_MODE", "wan")
	os.Setenv("VOIP_PORT", "9090")
	os.Setenv("VOIP_USERNAME", "testuser")
	os.Setenv("VOIP_DATA_DIR", "/tmp/voip-data")
	os.Setenv("VOIP_SERVER_ADDR", "192.168.1.1:9321")
	defer os.Unsetenv("VOIP_NETWORK_MODE")
	defer os.Unsetenv("VOIP_PORT")
	defer os.Unsetenv("VOIP_USERNAME")
	defer os.Unsetenv("VOIP_DATA_DIR")
	defer os.Unsetenv("VOIP_SERVER_ADDR")

	cfg := Load()

	if cfg.NetworkMode != NetworkModeWAN {
		t.Errorf("expected WAN mode, got %s", cfg.NetworkMode)
	}
	if cfg.Port != 9090 {
		t.Errorf("expected port 9090, got %d", cfg.Port)
	}
	if cfg.Username != "testuser" {
		t.Errorf("expected testuser, got %s", cfg.Username)
	}
	if cfg.DataDir != "/tmp/voip-data" {
		t.Errorf("expected /tmp/voip-data, got %s", cfg.DataDir)
	}
	if cfg.ServerAddr != "192.168.1.1:9321" {
		t.Errorf("expected 192.168.1.1:9321, got %s", cfg.ServerAddr)
	}
}

func TestLoadInvalidPort(t *testing.T) {
	os.Setenv("VOIP_PORT", "invalid")
	defer os.Unsetenv("VOIP_PORT")

	cfg := Load()
	if cfg.Port != 9321 {
		t.Errorf("expected default port 9321 for invalid env, got %d", cfg.Port)
	}
}

func TestLoadTURNConfig(t *testing.T) {
	os.Setenv("VOIP_TURN_URL", "turn:turn.example.com:3478")
	os.Setenv("VOIP_TURN_USERNAME", "user")
	os.Setenv("VOIP_TURN_PASSWORD", "pass")
	defer os.Unsetenv("VOIP_TURN_URL")
	defer os.Unsetenv("VOIP_TURN_USERNAME")
	defer os.Unsetenv("VOIP_TURN_PASSWORD")

	cfg := Load()

	if cfg.TURNConfig.URL != "turn:turn.example.com:3478" {
		t.Errorf("expected turn URL, got %s", cfg.TURNConfig.URL)
	}
	if cfg.TURNConfig.Username != "user" {
		t.Errorf("expected turn user, got %s", cfg.TURNConfig.Username)
	}
	if cfg.TURNConfig.Password != "pass" {
		t.Errorf("expected turn pass, got %s", cfg.TURNConfig.Password)
	}
}
