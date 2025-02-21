package loadbalancer

import (
	"fmt"
	"net/http"
	"time"

	"github.com/onkarbanerjee/roundbalancer/pkg/backends"
	"github.com/onkarbanerjee/roundbalancer/pkg/dispatcher"
	"github.com/onkarbanerjee/roundbalancer/pkg/livenesschecker"
	"github.com/onkarbanerjee/roundbalancer/pkg/loadbalancer"
	"go.uber.org/zap"
)

func Start(backendServers []*backends.Backend, logger *zap.Logger, livenessCheckInterval time.Duration, port int) error {
	backendGroup := backends.NewGroup(backendServers)

	d := dispatcher.New(
		loadbalancer.NewRoundRobin(backendGroup),
		livenesschecker.New(backendGroup),
		livenessCheckInterval,
		logger)
	logger.Info("created a dispatcher with configured backends that will check for their liveness at configured intervals",
		zap.Any("backends", backendGroup),
		zap.Duration("liveness_check_interval", livenessCheckInterval))

	go d.StartCheckingLiveness()

	http.HandleFunc("/echo", d.ServeHTTP)
	logger.Info("load balancer is starting", zap.Int("port", port))

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
