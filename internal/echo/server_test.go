package echo_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/onkarbanerjee/roundbalancer/internal/echo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestServer_Echo(t *testing.T) {
	t.Run("error - method not allowed", func(t *testing.T) {
		s := echo.NewServer("1", zap.NewExample())
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		s.Echo(w, r)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})
	t.Run("error - invalid json", func(t *testing.T) {
		s := echo.NewServer("1", zap.NewExample())
		r := httptest.NewRequest(http.MethodPost, "/", nil)
		w := httptest.NewRecorder()
		s.Echo(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("happy - echo successfully", func(t *testing.T) {
		type MockJSON struct {
			Game    string `json:"game"`
			GamerID string `json:"gamer_id"`
			Points  int    `json:"points"`
		}
		body := MockJSON{
			Game:    "Mobile Legends",
			GamerID: "GYUTDTE",
			Points:  20,
		}
		b, err := json.Marshal(body)
		assert.NoError(t, err)
		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
		w := httptest.NewRecorder()
		s := echo.NewServer("1", zap.NewExample())
		s.Echo(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
		var result MockJSON
		assert.NoError(t, json.NewDecoder(w.Body).Decode(&result))
		assert.Equal(t, body, result)
	})
}
