package livenesschecker

import (
	"fmt"
	"net/http"
	"time"

	"github.com/onkarbanerjee/roundbalancer/pkg/backend"
)

type LivenessChecker interface {
	CheckLiveness()
}

type livenessChecker struct {
	pool     backend.Pool
	interval time.Duration
}

func (l *livenessChecker) CheckLiveness() {
	allBackends := l.pool.GetAllBackends()
	for _, each := range allBackends {
		resp, err := http.Get(fmt.Sprintf("%s/livez", each.URL))
		each.UpdateHealthStatus(err == nil && resp.StatusCode == http.StatusOK)
	}
}

func New(pool backend.Pool, duration time.Duration) LivenessChecker {
	return &livenessChecker{pool: pool, interval: duration}
}
