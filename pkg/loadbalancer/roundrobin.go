package loadbalancer

import (
	"errors"
	"sync/atomic"

	"github.com/onkarbanerjee/roundbalancer/pkg/measurements"

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
	const operation = "roundRobin.Next"
	allBackends := r.backendGroup.GetAllBackends()
	total := int32(len(allBackends))
	next := (atomic.LoadInt32(&r.current) + 1) % total
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
		err := errors.New("no healthy backends")
		measurements.UpdateCount(operation, err)

		return nil, err
	}
	nextBackend := allBackends[next]
	atomic.StoreInt32(&r.current, next)
	measurements.UpdateCount(operation, nil)

	return nextBackend, nil
}
