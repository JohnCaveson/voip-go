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
	ID            string
	Username      string
	Addr          net.IP
	Port          int
	SignalingAddr string
}

type Discoverer struct {
	server  *mdns.Server
	entries chan *mdns.ServiceEntry
}

func NewDiscoverer(username string, port int, signalingAddr string) (*Discoverer, error) {
	info := []string{username}
	if signalingAddr != "" {
		info = append(info, "@"+signalingAddr)
	}
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
	for {
		select {
		case <-ctx.Done():
			return peers, ctx.Err()
		case entry, ok := <-entriesCh:
			if !ok {
				return peers, nil
			}

			username := ""
			signalingAddr := ""
			for _, field := range entry.InfoFields {
				if len(field) > 0 && field[0] == '@' {
					signalingAddr = field[1:]
				} else {
					username = field
				}
			}

			peer := PeerInfo{
				ID:            entry.Host,
				Username:      username,
				Addr:          entry.AddrV4,
				Port:          entry.Port,
				SignalingAddr: signalingAddr,
			}

			if entry.AddrV4 == nil {
				peer.Addr = entry.AddrV6
			}

			if peer.SignalingAddr == "" && peer.Addr != nil {
				peer.SignalingAddr = fmt.Sprintf("ws://%s:%d/signaling", peer.Addr.String(), peer.Port)
			}

			if peer.Addr != nil {
				peers = append(peers, peer)
			}
		}
	}
}

func (d *Discoverer) DiscoverWithTimeout(timeout time.Duration) ([]PeerInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return d.Discover(ctx)
}

func (d *Discoverer) Close() error {
	return d.server.Shutdown()
}
