package loadbalancer

import (
	"fmt"
	"net/http"
	"time"

	"github.com/onkarbanerjee/roundbalancer/pkg/backend"
	"github.com/onkarbanerjee/roundbalancer/pkg/dispatcher"
	"github.com/onkarbanerjee/roundbalancer/pkg/livenesschecker"
	"github.com/onkarbanerjee/roundbalancer/pkg/loadbalancer"
	"go.uber.org/zap"
)

func Start(backends []*backend.Backend, logger *zap.Logger) error {
	pool := backend.NewPool(backends)
	d := dispatcher.New(loadbalancer.NewRoundRobin(pool), livenesschecker.New(pool, 2*time.Second), logger)

	go d.StartCheckingLiveness()

	http.HandleFunc("/echo", d.ServeHTTP)

	port := 9090

	logger.Info("load balancer is starting", zap.Int("port", port))

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
