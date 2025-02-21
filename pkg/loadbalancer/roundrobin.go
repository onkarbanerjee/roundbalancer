package loadbalancer

import (
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
	next := (r.current + 1) % int32(r.pool.GetCount())
	j, nextBackEnd, err := r.pool.GetHealthyBackendAt(int(next))
	if err != nil {
		return nil, err
	}
	atomic.StoreInt32(&r.current, int32(j))
	return nextBackEnd, nil
}
