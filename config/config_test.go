package config_test

import (
	"os"
	"testing"

	"github.com/onkarbanerjee/roundbalancer/config"
	"github.com/stretchr/testify/assert"
)

const sampleJSON = `{
	"backends": [
		{"id": "backend1", "port": 8080},
		{"id": "backend2", "port": 8081}
	]
}`

func TestLoad(t *testing.T) {
	t.Run("failed to open config file", func(t *testing.T) {
		cfg, err := config.Load("bogus.json")
		assert.ErrorContains(t, err, "failed to open config file:")
		assert.Nil(t, cfg)
	})

	t.Run("failed to decode config file", func(t *testing.T) {
		temp, err := os.CreateTemp(".", "test.json")
		assert.NoError(t, err)
		defer os.Remove(temp.Name())

		_, err = temp.WriteString("")
		assert.NoError(t, err)
		cfg, err := config.Load(temp.Name())
		assert.ErrorContains(t, err, "failed to decode config JSON:")
		assert.Nil(t, cfg)
	})

	t.Run("successfully loads the config", func(t *testing.T) {
		temp, err := os.CreateTemp(".", "test.json")
		assert.NoError(t, err)
		defer os.Remove(temp.Name())

		_, err = temp.WriteString(sampleJSON)
		assert.NoError(t, err)
		cfg, err := config.Load(temp.Name())
		assert.NoError(t, err)

		if len(cfg.Backends) != 2 {
			t.Errorf("Expected 2 backends, got %d", len(cfg.Backends))
		}

		if cfg.Backends[0].ID != "backend1" {
			t.Errorf("Expected first backends ID to be 'backend1', got %q", cfg.Backends[0].ID)
		}

		if cfg.Backends[0].Port != 8080 {
			t.Errorf("Expected first backends port to be 8080, got %d", cfg.Backends[0].Port)
		}
		if cfg.Backends[1].ID != "backend2" {
			t.Errorf("Expected first backends ID to be 'backend1', got %q", cfg.Backends[0].ID)
		}

		if cfg.Backends[1].Port != 8081 {
			t.Errorf("Expected first backends port to be 8080, got %d", cfg.Backends[0].Port)
		}
	})
}
