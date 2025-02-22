package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
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
	port := flag.Int("port", 0, "port on which loadbalancer will serve")
	timeout := flag.Int("timeout", 0, "timeout in seconds")
	flag.Parse()
	if *port == 0 {
		log.Println("port is required")

		return
	}

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

	var backendServers []*backends.Backend
	for _, cfg := range cfg.Backends {
		baseURL := fmt.Sprintf("http://localhost:%d", cfg.Port)
		parsedBaseURL, err := url.Parse(baseURL)
		if err != nil {
			logger.Error(fmt.Sprintf("could not parse url %s, got error: %s", baseURL, err.Error()))

			return
		}

		liveZEndpoint := fmt.Sprintf("%s/livez", baseURL)
		parsedLivezURL, err := url.Parse(liveZEndpoint)
		if err != nil {
			logger.Error(fmt.Sprintf("could not parse url %s, got error: %s", liveZEndpoint, err.Error()))

			return
		}

		reversProxy := httputil.NewSingleHostReverseProxy(parsedBaseURL)
		reversProxy.Transport = &http.Transport{
			ResponseHeaderTimeout: time.Duration(*timeout) * time.Second,
		}
		reversProxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				http.Error(w, "Gateway Timeout", http.StatusGatewayTimeout)

				return
			}
			http.Error(w, "Network issues", http.StatusBadGateway)
		}
		backendServers = append(backendServers, backends.NewBackend(
			cfg.ID,
			reversProxy,
			parsedLivezURL))
	}

	if err := loadbalancer.Start(backendServers, logger, 2*time.Second, *port); err != nil {
		logger.Fatal(fmt.Sprintf("could not start server, got error: %s", err.Error()))
	}
}
