package healthupdater

import (
	"fmt"
	"net/http"
	"time"

	"github.com/onkarbanerjee/roundbalancer/pkg/backend"
)

type HealthUpdater struct {
	pool *backend.Pool
}

func New(pool *backend.Pool) *HealthUpdater {
	return &HealthUpdater{pool: pool}
}

func (h *HealthUpdater) Start() {
	t := time.NewTicker(2 * time.Second)

	for {
		<-t.C
		for _, each := range h.pool.Backends {
			get, err := http.Get(fmt.Sprintf("%s/livez", each.URL))
			each.UpdateHealthStatus(err != nil || get.StatusCode != http.StatusOK)
		}

	}
}
