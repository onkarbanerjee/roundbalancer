package backend

import (
	"net/http/httputil"
	"sync"
)

type Backend struct {
	ID        string
	URL       string
	isHealthy bool
	Service   *httputil.ReverseProxy
	mu        *sync.RWMutex
}

func (b *Backend) IsHealthy() bool {
	// get the read lock before reading its health status value
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.isHealthy
}

func (b *Backend) UpdateHealthStatus(isHealthy bool) {
	// lock it before updating the value
	b.mu.Lock()
	defer b.mu.Unlock()

	b.isHealthy = isHealthy
}

type Pool interface {
	GetHealthyBackends() []Backend
	GetAllBackends() []Backend
}
type pool struct {
	Backends []*Backend
}

func NewPool(backends []*Backend) Pool {
	return &pool{Backends: backends}
}

func (p *pool) GetHealthyBackends() []Backend {
	var healthyBackends []Backend
	for _, b := range p.Backends {
		if !b.IsHealthy() {
			continue
		}
		healthyBackends = append(healthyBackends, *b)
	}

	return healthyBackends
}

func (p *pool) GetAllBackends() []Backend {
	var allBackends []Backend
	for _, b := range p.Backends {
		allBackends = append(allBackends, *b)
	}

	return allBackends
}
