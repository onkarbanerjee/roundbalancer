package backend

import (
	"net/http/httputil"
	"sync"
)

type Backend struct {
	ID        string
	URL       string
	isHealthy bool
	*httputil.ReverseProxy
}

func NewBackend(id string, url string, proxy *httputil.ReverseProxy) *Backend {
	return &Backend{ID: id, URL: url, isHealthy: true, ReverseProxy: proxy}
}

func (b *Backend) IsHealthy() bool {
	return b.isHealthy
}

func (b *Backend) UpdateHealthStatus(isHealthy bool) {
	b.isHealthy = isHealthy
}

type Pool struct {
	Backends []*Backend
	sync.RWMutex
}

func NewPool(backends []*Backend) *Pool {
	return &Pool{Backends: backends}
}
func (p *Pool) GetHealthyBackends() []Backend {
	p.RLock()
	defer p.RUnlock()

	var healthyBackends []Backend
	for _, b := range p.Backends {
		if !b.IsHealthy() {
			continue
		}
		healthyBackends = append(healthyBackends, *b)
	}

	return healthyBackends
}
