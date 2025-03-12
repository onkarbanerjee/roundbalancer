package dispatcher

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/onkarbanerjee/roundbalancer/pkg/measurements"

	"github.com/onkarbanerjee/roundbalancer/pkg/livenesschecker"
	"github.com/onkarbanerjee/roundbalancer/pkg/loadbalancer"
	"go.uber.org/zap"
)

type Dispatcher struct {
	loadBalancer          loadbalancer.LoadBalancer
	livenessChecker       livenesschecker.LivenessChecker
	livenessCheckInterval time.Duration
	logger                *zap.Logger
}

func New(loadBalancer loadbalancer.LoadBalancer, checker livenesschecker.LivenessChecker, livesnessCheckInterval time.Duration, logger *zap.Logger) *Dispatcher {
	return &Dispatcher{
		loadBalancer:          loadBalancer,
		livenessChecker:       checker,
		livenessCheckInterval: livesnessCheckInterval,
		logger:                logger}
}
func (r *Dispatcher) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	const operation = "dispatcher.ServeHTTP"
	if request.Method != http.MethodPost {
		r.logger.Error("method not allowed", zap.String("method", request.Method))
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		measurements.UpdateCount(operation, errors.New("method not allowed"))

		return
	}
	backendServer, err := r.loadBalancer.Next()
	if err != nil || backendServer == nil {
		r.logger.Error("could not get next backends to route this request to", zap.Error(err))
		http.Error(writer, "failed to route this request", http.StatusInternalServerError)
		measurements.UpdateCount(operation, errors.New("failed to route this request"))

		return
	}
	r.logger.Info(fmt.Sprintf("dispatching to backend ID: %s", backendServer.ID))
	measurements.UpdateCount(operation, nil)

	backendServer.Service.ServeHTTP(writer, request)
}

func (r *Dispatcher) StartCheckingLiveness() {
	t := time.NewTicker(r.livenessCheckInterval)
	defer t.Stop()

	for range t.C {
		r.livenessChecker.CheckLiveness()
	}
}
