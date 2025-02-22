package backends

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Backend struct {
	ID          string
	Service     *httputil.ReverseProxy
	LivenessURL *url.URL
	isHealthy   bool
	mu          *sync.RWMutex
}

func NewBackend(id string, service *httputil.ReverseProxy, url2 *url.URL) *Backend {
	return &Backend{
		ID:          id,
		Service:     service,
		LivenessURL: url2,
		isHealthy:   false,
		mu:          &sync.RWMutex{},
	}
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
	fmt.Println("seting health status", isHealthy, b.ID)
	b.isHealthy = isHealthy
}

type GroupOfBackends interface {
	GetAllBackends() []*Backend
}

type group struct {
	group []*Backend
}

func (g *group) GetAllBackends() []*Backend {
	return g.group
}

func NewGroup(backendServers []*Backend) GroupOfBackends {
	return &group{group: backendServers}
}
