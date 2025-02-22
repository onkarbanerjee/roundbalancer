package dispatcher_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"
	"time"

	"github.com/onkarbanerjee/roundbalancer/mocks"
	"github.com/onkarbanerjee/roundbalancer/pkg/backends"
	"github.com/onkarbanerjee/roundbalancer/pkg/dispatcher"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestDispatcher_ServeHTTP(t *testing.T) {
	t.Run("error - could not get a backend from loadbalancer to send request to it", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockLB := mocks.NewMockLoadBalancer(ctrl)
		mockLB.EXPECT().Next().Return(nil, errors.New("no backend"))

		d := dispatcher.New(mockLB, nil, time.Duration(0), zap.NewExample())
		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		rec := httptest.NewRecorder()

		d.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "failed to route this request")
	})
	t.Run("error - could not get a backend from loadbalancer to send request to it", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockLB := mocks.NewMockLoadBalancer(ctrl)

		d := dispatcher.New(mockLB, nil, time.Duration(0), zap.NewExample())
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		d.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
		assert.Contains(t, rec.Body.String(), "method not allowed")
	})
	t.Run("happy - it is able to get a backend from loadbalancer and send request to it", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockLB := mocks.NewMockLoadBalancer(ctrl)
		d := dispatcher.New(mockLB, nil, time.Duration(0), zap.NewExample())

		backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`Backend Response`))
		}))
		defer backendServer.Close()

		backendURL, err := url.Parse(backendServer.URL)
		assert.NoError(t, err)
		mockLB.EXPECT().Next().Return(&backends.Backend{
			Service: httputil.NewSingleHostReverseProxy(backendURL),
		}, nil)
		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		rec := httptest.NewRecorder()

		d.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Backend Response", rec.Body.String())
	})
}
