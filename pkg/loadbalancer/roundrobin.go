package loadbalancer

import (
	"errors"
	"sync/atomic"

	"github.com/onkarbanerjee/roundbalancer/pkg/backends"
)

type LoadBalancer interface {
	Next() (*backends.Backend, error)
}
type RoundRobin struct {
	backendGroup backends.GroupOfBackends
	current      int32
}

func NewRoundRobin(backendGroup backends.GroupOfBackends) *RoundRobin {
	return &RoundRobin{backendGroup: backendGroup, current: -1}
}

func (r *RoundRobin) Next() (*backends.Backend, error) {
	allBackends := r.backendGroup.GetAllBackends()
	total := int32(len(allBackends))
	next := (r.current + 1) % total
	var found bool
	for i := next; i < next+total; i++ {
		j := i % total
		if allBackends[j].IsHealthy() {
			next = j
			found = true

			break
		}
	}
	if !found {
		return nil, errors.New("no healthy backends")
	}
	nextBackend := allBackends[next]
	atomic.StoreInt32(&r.current, next)
	return nextBackend, nil
}
