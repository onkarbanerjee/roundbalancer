package strategies

import (
	"github.com/onkarbanerjee/roundbalancer/pkg/backend"
)

type RoundRobin struct {
	pool    *backend.Pool
	current int
}

func NewRoundRobin(pool *backend.Pool) *RoundRobin {
	return &RoundRobin{pool: pool, current: -1}
}

func (s *RoundRobin) Next() backend.Backend {
	healthyBackends := s.pool.GetHealthyBackends()
	next := (s.current + 1) % len(healthyBackends)
	s.current = next

	return healthyBackends[next]
}
