package livenesschecker

import (
	"fmt"
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
		fmt.Println(fmt.Sprintf("%s", each.LivenessURL.String()))
		resp, err := http.Get(fmt.Sprintf("%s", each.LivenessURL.String()))
		each.UpdateHealthStatus(err == nil && resp.StatusCode == http.StatusOK)
	}
	fmt.Println("all backend health updated")
}

func New(pool backends.GroupOfBackends) LivenessChecker {
	return &livenessChecker{pool: pool}
}
