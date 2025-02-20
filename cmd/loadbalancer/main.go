package main

import (
	"fmt"
	"log"
	"net/http/httputil"
	"net/url"
	"path"

	"github.com/onkarbanerjee/roundbalancer/config"
	"github.com/onkarbanerjee/roundbalancer/internal/loadbalancer"
	"github.com/onkarbanerjee/roundbalancer/pkg/backend"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Printf("could not create logger, got error: %s", err.Error())

		return
	}

	cfg, err := config.Load(path.Join("config", "values", "config.json"))
	if err != nil {
		logger.Error(fmt.Sprintf("could not load config.json, got error: %s", err.Error()))

		return
	}

	var backends []*backend.Backend
	for _, backendConfig := range cfg.Backends {
		endpoint := fmt.Sprintf("http://localhost:%d", backendConfig.Port)
		parse, err := url.Parse(endpoint)
		if err != nil {
			logger.Error(fmt.Sprintf("could not parse url %s, got error: %s", endpoint, err.Error()))

			return
		}
		logger.Info(fmt.Sprintf("backend url: %s", parse.String()))
		backends = append(backends, &backend.Backend{ID: backendConfig.ID, URL: endpoint, Service: httputil.NewSingleHostReverseProxy(parse)})
	}

	if err := loadbalancer.Start(backends, logger); err != nil {
		logger.Fatal(fmt.Sprintf("could not start server, got error: %s", err.Error()))
	}
}
