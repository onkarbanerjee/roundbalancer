package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/onkarbanerjee/roundbalancer/internal/echo"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Printf("could not create logger, got error: %s", err.Error())

		return
	}

	id := flag.String("id", "", "server id")
	port := flag.Int("port", 0, "port on which echo server will serve")
	flag.Parse()
	if *id == "" || *port == 0 {
		logger.Error("id or port is required")

		return
	}
	server := echo.NewServer(*id, logger)
	http.HandleFunc("/echo", server.Echo)
	http.HandleFunc("/livez", server.Liveness)

	fmt.Printf("running echo server id %s on port %d\n", *id, *port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		logger.Fatal(fmt.Sprintf("could not start server, got error: %s", err.Error()))
	}
}
