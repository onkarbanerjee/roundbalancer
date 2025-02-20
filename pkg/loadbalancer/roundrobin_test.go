package loadbalancer_test

import (
	"testing"

	"github.com/onkarbanerjee/roundbalancer/mocks"
	"github.com/onkarbanerjee/roundbalancer/pkg/backend"
	"github.com/onkarbanerjee/roundbalancer/pkg/loadbalancer"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRoundRobin_Next(t *testing.T) {
	t.Run("it should error when no healthy backends", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPool := mocks.NewMockPool(ctrl)
		mockPool.EXPECT().GetHealthyBackends().Return([]backend.Backend{})
		r := loadbalancer.NewRoundRobin(mockPool)
		b, err := r.Next()
		assert.ErrorContains(t, err, "no healthy backends")
		assert.Nil(t, b)
	})
	t.Run("it should always send next backend in order", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPool := mocks.NewMockPool(ctrl)
		expectedNexts := []string{"1", "2", "3", "1", "2", "3", "1", "2"}
		mockPool.EXPECT().GetHealthyBackends().Return([]backend.Backend{
			{ID: "1"},
			{ID: "2"},
			{ID: "3"},
		}).Times(len(expectedNexts))
		r := loadbalancer.NewRoundRobin(mockPool)
		for _, want := range expectedNexts {
			b, err := r.Next()
			assert.NoError(t, err)
			assert.Equal(t, want, b.ID)
		}
	})
}
