package dispatcher

import (
	"fmt"
	"net/http"
	"time"

	"github.com/onkarbanerjee/roundbalancer/pkg/livenesschecker"
	"github.com/onkarbanerjee/roundbalancer/pkg/loadbalancer"
	"go.uber.org/zap"
)

type Dispatcher struct {
	loadBalancer    loadbalancer.LoadBalancer
	livenessChecker livenesschecker.LivenessChecker
	logger          *zap.Logger
}

func New(loadBalancer loadbalancer.LoadBalancer, checker livenesschecker.LivenessChecker, logger *zap.Logger) *Dispatcher {
	return &Dispatcher{loadBalancer: loadBalancer, livenessChecker: checker, logger: logger}
}
func (r *Dispatcher) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	backendServer, err := r.loadBalancer.Next()
	if err != nil || backendServer == nil {
		r.logger.Error("could not get next backend to route this request to", zap.Error(err))
		writer.WriteHeader(http.StatusServiceUnavailable)
		writer.Write([]byte("failed to route this request"))

		return
	}
	r.logger.Info(fmt.Sprintf("dispatching to backend ID: %s", backendServer.ID))
	backendServer.Service.ServeHTTP(writer, request)
}

func (r *Dispatcher) StartCheckingLiveness() {
	t := time.NewTicker(2 * time.Second)
	defer t.Stop()

	for range t.C {
		r.livenessChecker.CheckLiveness()
	}
}
