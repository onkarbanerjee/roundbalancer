package loadbalancer_test

import (
	"testing"

	"github.com/onkarbanerjee/roundbalancer/mocks"
	"github.com/onkarbanerjee/roundbalancer/pkg/backends"
	"github.com/onkarbanerjee/roundbalancer/pkg/loadbalancer"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRoundRobin_Next(t *testing.T) {
	t.Run("it should error when no healthy backends", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPool := mocks.NewMockGroupOfBackends(ctrl)
		b1 := backends.NewBackend("1", nil, nil)
		b2 := backends.NewBackend("2", nil, nil)
		b3 := backends.NewBackend("3", nil, nil)
		mockBackends := []*backends.Backend{b1, b2, b3}
		mockPool.EXPECT().GetAllBackends().Return(mockBackends)

		r := loadbalancer.NewRoundRobin(mockPool)
		b, err := r.Next()
		assert.ErrorContains(t, err, "no healthy backends")
		assert.Nil(t, b)
	})
	t.Run("it should always send next backends in order", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPool := mocks.NewMockGroupOfBackends(ctrl)
		r := loadbalancer.NewRoundRobin(mockPool)

		b1 := backends.NewBackend("1", nil, nil)
		b2 := backends.NewBackend("2", nil, nil)
		b3 := backends.NewBackend("3", nil, nil)
		mockBackends := []*backends.Backend{b1, b2, b3}
		b1.UpdateHealthStatus(true)
		b2.UpdateHealthStatus(true)
		b3.UpdateHealthStatus(true)
		mockPool.EXPECT().GetAllBackends().Return(mockBackends).AnyTimes()

		b, err := r.Next()
		assert.NoError(t, err)
		assert.Equal(t, "1", b.ID)

		b, err = r.Next()
		assert.NoError(t, err)
		assert.Equal(t, "2", b.ID)

		b, err = r.Next()
		assert.NoError(t, err)
		assert.Equal(t, "3", b.ID)

		b, err = r.Next()
		assert.NoError(t, err)
		assert.Equal(t, "1", b.ID)

		b2.UpdateHealthStatus(false)
		b, err = r.Next()
		assert.NoError(t, err)
		assert.Equal(t, "3", b.ID)

		b, err = r.Next()
		assert.NoError(t, err)
		assert.Equal(t, "1", b.ID)
	})
}
