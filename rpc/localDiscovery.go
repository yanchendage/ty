package rpc

import (
	"errors"
	"math"
	"math/rand"
	"sync"
	"time"
)

type LocalDiscovery struct {
	r *rand.Rand
	mu sync.RWMutex
	servers []string
	index int
}

func NewLocalDiscovery(servers []string) *LocalDiscovery {
	d := &LocalDiscovery{
		servers: servers,
		r:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	d.index = d.r.Intn(math.MaxInt32 - 1)
	return d
}

var _ IDiscovery = (*LocalDiscovery)(nil)

func (d *LocalDiscovery) Refresh() error {
	return nil
}

func (d *LocalDiscovery) Update(servers []string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.servers = servers
	return nil
}

func (d *LocalDiscovery) Get(mode LoadMode) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	n := len(d.servers)
	if n == 0 {
		return "", errors.New("rpc discovery: no available servers")
	}
	switch mode {
	case Random:
		return d.servers[d.r.Intn(n)], nil
	case RoundRobin:
		s := d.servers[d.index%n] // servers could be updated, so mode n to ensure safety
		d.index = (d.index + 1) % n
		return s, nil
	default:
		return "", errors.New("rpc discovery: not supported select mode")
	}
}

func (d *LocalDiscovery) GetAll() ([]string, error) {
	d.mu.RLocker()
	defer d.mu.RUnlock()
	// return a copy of d.servers
	servers := make([]string, len(d.servers), len(d.servers))
	copy(servers, d.servers)
	return servers, nil
}

