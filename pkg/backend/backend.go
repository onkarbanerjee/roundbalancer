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
	mu *sync.RWMutex
}

func NewBackend(id string, url string, proxy *httputil.ReverseProxy) *Backend {
	return &Backend{ID: id, URL: url, isHealthy: true, ReverseProxy: proxy, mu: &sync.RWMutex{}}
}

func (b *Backend) IsHealthy() bool {
	// get the read lock before reading it's health status value
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

type Pool struct {
	Backends []*Backend
	sync.RWMutex
}

func NewPool(backends []*Backend) *Pool {
	return &Pool{Backends: backends}
}

func (p *Pool) Add(backend *Backend) {
	p.Lock()
	defer p.Unlock()
	
	p.Backends = append(p.Backends, backend)
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
