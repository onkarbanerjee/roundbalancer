package healthupdater

import (
	"github.com/onkarbanerjee/roundbalancer/pkg/backend"
	"time"
)

type HealthUpdater struct {
	pool        *backend.Pool
	checkHealth func(url string) bool
}

func New(pool *backend.Pool) *HealthUpdater {
	return &HealthUpdater{pool: pool}
}

func (h *HealthUpdater) Start() {
	t := time.NewTicker(1 * time.Second)

	for range t.C {
		for _, each := range h.pool.Backends {
			each.UpdateHealthStatus(h.checkHealth(each.URL))
		}
	}
}
