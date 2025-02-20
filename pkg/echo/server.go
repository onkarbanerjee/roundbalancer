package echo

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type Server struct {
	id     string
	logger *zap.Logger
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		s.logger.Error("Invalid method used", zap.String("method", r.Method))

		return
	}

	var v map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		s.logger.Error("Invalid JSON", zap.Error(err))

		return
	}

	defer func() {
		err := r.Body.Close()
		if err != nil {
			s.logger.Error("failed to close request body", zap.Error(err))
		}
	}()

	w.Header().Set("Content-Type", "echo/json")
	err = json.NewEncoder(w).Encode(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		s.logger.Error("failed to write response", zap.Error(err))

		return
	}

	s.logger.Info("Request completed", zap.String("method", r.Method), zap.String("url", r.URL.String()))
}

func NewServer(id string, logger *zap.Logger) *Server {
	return &Server{id: id, logger: logger.With(zap.String("server_id", id))}
}

func Start(cmd *cobra.Command, args []string) error {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Printf("could not create logger, got error: %s", err.Error())

		return err
	}

	id, err := cmd.Flags().GetString("id")
	if err != nil {
		logger.Error("failed to parse id", zap.Error(err))

		return fmt.Errorf("failed to get id flag: %w", err)
	}
	port, err := cmd.Flags().GetInt("port")
	if err != nil {
		logger.Error("failed to parse port", zap.Error(err))

		return fmt.Errorf("failed to get port flag: %w", err)
	}
	echoServer := NewServer(id, logger)
	http.HandleFunc("/echo", echoServer.ServeHTTP)

	fmt.Printf("running application server id %s on port %d\n", id, port)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
