package echo

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type Server struct {
	id     string
	logger *zap.Logger
}

func NewServer(id string, logger *zap.Logger) *Server {
	return &Server{id: id, logger: logger.With(zap.String("server_id", id))}
}
func (s *Server) Echo(w http.ResponseWriter, r *http.Request) {
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

	fmt.Println(v)
	w.Header().Set("Content-Type", "echo/json")
	err = json.NewEncoder(w).Encode(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		s.logger.Error("failed to write response", zap.Error(err))

		return
	}

	s.logger.Info("Request completed", zap.String("method", r.Method), zap.String("url", r.URL.String()))
}

func (s *Server) Liveness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
