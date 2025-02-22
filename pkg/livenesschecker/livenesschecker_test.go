package livenesschecker_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/onkarbanerjee/roundbalancer/mocks"
	"github.com/onkarbanerjee/roundbalancer/pkg/backends"
	"github.com/onkarbanerjee/roundbalancer/pkg/livenesschecker"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestLivenessChecker_CheckLiveness(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockGroupOfBackends(ctrl)
	l := livenesschecker.New(mockPool)

	shouldMockServer1BeHealthy := true
	shouldMockServer2BeHealthy := true
	mockBackendServer1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !shouldMockServer1BeHealthy {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}
	}))
	defer mockBackendServer1.Close()

	mockBackendServer1URL, err := url.Parse(mockBackendServer1.URL)
	assert.NoError(t, err)

	mockBackendServer2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !shouldMockServer2BeHealthy {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}
	}))
	defer mockBackendServer2.Close()

	mockBackendServer2URL, err := url.Parse(mockBackendServer2.URL)
	assert.NoError(t, err)

	b1 := backends.NewBackend("1", nil, mockBackendServer1URL)
	b2 := backends.NewBackend("2", nil, mockBackendServer2URL)
	mockPool.EXPECT().GetAllBackends().Return([]*backends.Backend{b1, b2}).AnyTimes()

	// all servers responding 200OK, hence they should be healthy
	l.CheckLiveness()
	assert.True(t, b1.IsHealthy())
	assert.True(t, b2.IsHealthy())

	// mockBackendServer2 is responding 500, so it should be marked as unhealthy
	shouldMockServer2BeHealthy = false
	l.CheckLiveness()
	assert.True(t, b1.IsHealthy())
	assert.False(t, b2.IsHealthy())

	// mockBackendServer1 is also responding 500, so it should be marked as unhealthy
	shouldMockServer1BeHealthy = false
	l.CheckLiveness()
	assert.False(t, b1.IsHealthy())
	assert.False(t, b2.IsHealthy())

	// both servers responding 200OK, so they should be marked as healthy
	shouldMockServer1BeHealthy = true
	shouldMockServer2BeHealthy = true
	l.CheckLiveness()
	assert.True(t, b1.IsHealthy())
	assert.True(t, b2.IsHealthy())
}
