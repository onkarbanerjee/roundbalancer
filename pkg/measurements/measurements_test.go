package measurements_test

import (
	"errors"
	"testing"

	"github.com/onkarbanerjee/roundbalancer/pkg/measurements"
	"github.com/stretchr/testify/assert"
)

func TestUpdateCount(t *testing.T) {
	assert.NotPanics(t, func() {
		measurements.UpdateCount("operation1", nil)
	})
	assert.NotPanics(t, func() {
		measurements.UpdateCount("operation2", errors.New("some error"))
	})
}
