package loadbalancer

import (
	"errors"
	"log"
	"sync/atomic"

	"github.com/onkarbanerjee/roundbalancer/pkg/backend"
)

type LoadBalancer interface {
	Next() (*backend.Backend, error)
}
type RoundRobin struct {
	pool    backend.Pool
	current int32
}

func NewRoundRobin(pool backend.Pool) *RoundRobin {
	return &RoundRobin{pool: pool, current: -1}
}

func (r *RoundRobin) Next() (*backend.Backend, error) {
	healthyBackends := r.pool.GetHealthyBackends()
	if len(healthyBackends) == 0 {
		log.Println("No healthy backends")

		return nil, errors.New("no healthy backends")
	}
	next := (r.current + 1) % int32(len(healthyBackends))
	atomic.StoreInt32(&r.current, next)

	return &healthyBackends[next], nil
}
