package livenesschecker

import (
	"net/http"

	"github.com/onkarbanerjee/roundbalancer/pkg/backends"
)

type LivenessChecker interface {
	CheckLiveness()
}

type livenessChecker struct {
	pool backends.GroupOfBackends
}

func (l *livenessChecker) CheckLiveness() {
	allBackends := l.pool.GetAllBackends()
	for _, each := range allBackends {
		resp, err := http.Get(each.LivenessURL.String())
		each.UpdateHealthStatus(err == nil && resp.StatusCode == http.StatusOK)
	}
}

func New(pool backends.GroupOfBackends) LivenessChecker {
	return &livenessChecker{pool: pool}
}
