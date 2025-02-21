package loadbalancer_test

import (
	"errors"
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
		mockPool.EXPECT().GetCount().Return(3)
		mockPool.EXPECT().GetHealthyBackendAt(0).Return(nil, errors.New("no healthy backends"))
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

		mockPool.EXPECT().GetCount().Return(3).AnyTimes()
		mockPool.EXPECT().GetHealthyBackendAt(0).Return(&backends.Backend{
			ID: "1",
		}, nil)
		b, err := r.Next()
		assert.NoError(t, err)
		assert.Equal(t, "1", b.ID)

		mockPool.EXPECT().GetHealthyBackendAt(1).Return(&backends.Backend{
			ID: "2",
		}, nil)
		b, err = r.Next()
		assert.NoError(t, err)
		assert.Equal(t, "2", b.ID)

		mockPool.EXPECT().GetHealthyBackendAt(2).Return(&backends.Backend{
			ID: "3",
		}, nil)
		b, err = r.Next()
		assert.NoError(t, err)
		assert.Equal(t, "3", b.ID)

		mockPool.EXPECT().GetHealthyBackendAt(0).Return(&backends.Backend{
			ID: "1",
		}, nil)
		b, err = r.Next()
		assert.NoError(t, err)
		assert.Equal(t, "1", b.ID)

		mockPool.EXPECT().GetHealthyBackendAt(1).Return(&backends.Backend{
			ID: "3",
		}, nil)
		b, err = r.Next()
		assert.NoError(t, err)
		assert.Equal(t, "3", b.ID)

		mockPool.EXPECT().GetHealthyBackendAt(2).Return(&backends.Backend{
			ID: "1",
		}, nil)
		b, err = r.Next()
		assert.NoError(t, err)
		assert.Equal(t, "1", b.ID)
	})
}
