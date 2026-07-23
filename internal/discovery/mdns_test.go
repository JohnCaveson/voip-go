package discovery

import (
	"testing"
	"time"
)

func TestNewDiscoverer(t *testing.T) {
	d, err := NewDiscoverer("testuser", 9321, "ws://localhost:9321/signaling")
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
}

func TestDiscoverWithTimeout(t *testing.T) {
	d1, err := NewDiscoverer("user1", 9321, "ws://localhost:9321/signaling")
	if err != nil {
		t.Fatal(err)
	}
	defer d1.Close()

	time.Sleep(500 * time.Millisecond)

	peers, err := d1.DiscoverWithTimeout(2 * time.Second)
	if err != nil {
		t.Fatal(err)
	}

	_ = peers
}

func TestClose(t *testing.T) {
	d, err := NewDiscoverer("testuser", 9321, "ws://localhost:9321/signaling")
	if err != nil {
		t.Fatal(err)
	}

	if err := d.Close(); err != nil {
		t.Fatal(err)
	}
}
