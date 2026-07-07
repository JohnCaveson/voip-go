package discovery

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/hashicorp/mdns"
)

const serviceName = "_voip-app._tcp"

type PeerInfo struct {
	ID       string
	Username string
	Addr     net.IP
	Port     int
}

type Discoverer struct {
	server  *mdns.Server
	entries chan *mdns.ServiceEntry
}

func NewDiscoverer(username string, port int) (*Discoverer, error) {
	info := []string{username}
	service, err := mdns.NewMDNSService(
		username,
		serviceName,
		"",
		"",
		port,
		nil,
		info,
	)
	if err != nil {
		return nil, fmt.Errorf("new mdns service: %w", err)
	}

	server, err := mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		return nil, fmt.Errorf("mdns server: %w", err)
	}

	return &Discoverer{
		server:  server,
		entries: make(chan *mdns.ServiceEntry, 10),
	}, nil
}

func (d *Discoverer) Discover(ctx context.Context) ([]PeerInfo, error) {
	entriesCh := make(chan *mdns.ServiceEntry, 10)

	go func() {
		mdns.Lookup(serviceName, entriesCh)
		close(entriesCh)
	}()

	var peers []PeerInfo
	for entry := range entriesCh {
		username := ""
		if len(entry.InfoFields) > 0 {
			username = entry.InfoFields[0]
		}

		peer := PeerInfo{
			ID:       entry.Host,
			Username: username,
			Addr:     entry.AddrV4,
			Port:     entry.Port,
		}

		if entry.AddrV4 == nil {
			peer.Addr = entry.AddrV6
		}

		if peer.Addr != nil {
			peers = append(peers, peer)
		}
	}

	return peers, nil
}

func (d *Discoverer) DiscoverWithTimeout(timeout time.Duration) ([]PeerInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return d.Discover(ctx)
}

func (d *Discoverer) Close() error {
	return d.server.Shutdown()
}
