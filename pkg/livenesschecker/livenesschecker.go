package livenesschecker

import (
	"net/http"

	"github.com/onkarbanerjee/roundbalancer/pkg/backends"
)

type LivenessChecker interface {
	CheckLiveness()
}

type livenessChecker struct {
	backendGroup backends.GroupOfBackends
}

func (l *livenessChecker) CheckLiveness() {
	allBackends := l.backendGroup.GetAllBackends()
	for _, each := range allBackends {
		resp, err := http.Get(each.LivenessURL.String())
		latestHealthStatus := err == nil && resp.StatusCode == http.StatusOK
		if latestHealthStatus != each.IsHealthy() {
			each.UpdateHealthStatus(latestHealthStatus)
		}
	}
}

func New(backendGroup backends.GroupOfBackends) LivenessChecker {
	return &livenessChecker{backendGroup: backendGroup}
}
