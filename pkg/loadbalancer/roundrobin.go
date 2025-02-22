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
	pool    backends.GroupOfBackends
	current int32
}

func NewRoundRobin(pool backends.GroupOfBackends) *RoundRobin {
	return &RoundRobin{pool: pool, current: -1}
}

func (r *RoundRobin) Next() (*backends.Backend, error) {
	allBackends := r.pool.GetAllBackends()
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
