package loadbalancer

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"

	"github.com/onkarbanerjee/roundbalancer/config"
	"github.com/onkarbanerjee/roundbalancer/pkg/backend"
	"github.com/onkarbanerjee/roundbalancer/pkg/healthupdater"
	"github.com/onkarbanerjee/roundbalancer/pkg/strategies"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func Start(cmd *cobra.Command, args []string) error {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Printf("could not create logger, got error: %s", err.Error())

		return err
	}

	cfg, err := config.Load(path.Join("config", "values", "config.json"))
	if err != nil {
		logger.Error(fmt.Sprintf("could not load config.json, got error: %s", err.Error()))

		return err
	}

	var backends []*backend.Backend
	for _, backendConfig := range cfg.Backends {
		endpoint := fmt.Sprintf("http://localhost:%d", backendConfig.Port)
		parse, err := url.Parse(endpoint)
		if err != nil {
			logger.Error(fmt.Sprintf("could not parse url %s, got error: %s", endpoint, err.Error()))

			return err
		}
		logger.Info(fmt.Sprintf("backend url: %s", parse.String()))
		backends = append(backends, backend.NewBackend(backendConfig.ID, endpoint, httputil.NewSingleHostReverseProxy(parse)))
	}

	pool := backend.NewPool(backends)

	l := &LoadBalancer{
		healthUpdater: healthupdater.New(pool),
		roundRobin:    strategies.NewRoundRobin(pool),
		logger:        logger,
	}

	go l.healthUpdater.Start()

	http.HandleFunc("/echo", l.ServeHTTP)

	port := 9090
	fmt.Printf("load balancer is running on port %d\n", port)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

type LoadBalancer struct {
	healthUpdater *healthupdater.HealthUpdater
	roundRobin    *strategies.RoundRobin
	logger        *zap.Logger
}

func (l *LoadBalancer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	backendServer := l.roundRobin.Next()
	l.logger.Info(fmt.Sprintf("dispatching to backend ID: %s", backendServer.ID))
	backendServer.ServeHTTP(writer, request)
}
