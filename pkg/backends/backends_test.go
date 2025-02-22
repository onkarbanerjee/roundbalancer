package backends_test

import (
	"net/http/httputil"
	"net/url"
	"testing"

	"github.com/onkarbanerjee/roundbalancer/pkg/backends"
)

// Test creating a new Backend
func TestNewBackend(t *testing.T) {
	targetURL, _ := url.Parse("http://example.com")
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	backend := backends.NewBackend("backend-1", proxy, targetURL)

	if backend.ID != "backend-1" {
		t.Errorf("Expected ID 'backend-1', got %s", backend.ID)
	}
	if backend.Service != proxy {
		t.Errorf("Expected proxy service to be set")
	}
	if backend.LivenessURL.String() != "http://example.com" {
		t.Errorf("Expected LivenessURL 'http://example.com', got %s", backend.LivenessURL.String())
	}
	if backend.IsHealthy() {
		t.Errorf("Expected initial health status to be false")
	}
}

func TestBackendHealthStatus(t *testing.T) {
	targetURL, _ := url.Parse("http://example.com")
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	backend := backends.NewBackend("backend-1", proxy, targetURL)

	if backend.IsHealthy() {
		t.Errorf("Expected initial health status to be false")
	}

	backend.UpdateHealthStatus(true)
	if !backend.IsHealthy() {
		t.Errorf("Expected health status to be true after update")
	}

	backend.UpdateHealthStatus(false)
	if backend.IsHealthy() {
		t.Errorf("Expected health status to be false after update")
	}
}

func TestNewGroup(t *testing.T) {
	targetURL, _ := url.Parse("http://example.com")
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	backend1 := backends.NewBackend("backend-1", proxy, targetURL)
	backend2 := backends.NewBackend("backend-2", proxy, targetURL)

	group := backends.NewGroup([]*backends.Backend{backend1, backend2})

	allBackends := group.GetAllBackends()
	if len(allBackends) != 2 {
		t.Errorf("Expected 2 allBackends, got %d", len(allBackends))
	}
	if allBackends[0].ID != "backend-1" || allBackends[1].ID != "backend-2" {
		t.Errorf("Backend IDs mismatch")
	}
}
