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
	backendServer, err := r.loadBalancer.Next()
	if err != nil || backendServer == nil {
		r.logger.Error("could not get next backends to route this request to", zap.Error(err))
		writer.WriteHeader(http.StatusServiceUnavailable)
		writer.Write([]byte("failed to route this request"))

		return
	}
	r.logger.Info(fmt.Sprintf("dispatching to backends ID: %s", backendServer.ID))
	fmt.Println("request is", request.URL.String())
	(backendServer.Service).ServeHTTP(writer, request)
}

func (r *Dispatcher) StartCheckingLiveness() {
	t := time.NewTicker(r.livenessCheckInterval)
	defer t.Stop()

	for range t.C {
		fmt.Println("checking liveness")
		r.livenessChecker.CheckLiveness()
	}
}
