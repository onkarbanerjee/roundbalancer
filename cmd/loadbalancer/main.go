package main

import (
	"fmt"
	"log"
	"net/http/httputil"
	"net/url"
	"path"
	"time"

	"github.com/onkarbanerjee/roundbalancer/config"
	"github.com/onkarbanerjee/roundbalancer/internal/loadbalancer"
	"github.com/onkarbanerjee/roundbalancer/pkg/backends"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Printf("could not create logger, got error: %s", err.Error())

		return
	}

	cfg, err := config.Load(path.Join("config", "values", "proxies_config.json"))
	if err != nil {
		logger.Error(fmt.Sprintf("could not load proxies_config.json, got error: %s", err.Error()))

		return
	}

	var bw []*backends.Backend
	for _, backendConfig := range cfg.Backends {
		endpoint := fmt.Sprintf("http://localhost:%d", backendConfig.Port)
		parse, err := url.Parse(endpoint)
		if err != nil {
			logger.Error(fmt.Sprintf("could not parse url %s, got error: %s", endpoint, err.Error()))

			return
		}
		logger.Info(fmt.Sprintf("backends url: %s", parse.String()))
		bw = append(bw, backends.NewBackend(
			backendConfig.ID, httputil.NewSingleHostReverseProxy(parse), parse))
	}

	if err := loadbalancer.Start(bw, logger, 2*time.Second, 9090); err != nil {
		logger.Fatal(fmt.Sprintf("could not start server, got error: %s", err.Error()))
	}
}
