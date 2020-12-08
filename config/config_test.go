package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	localConfig := GetConfig("localhost")

	assert.Equal(t, localConfig.Redis.ConnectionType, "tcp")
	assert.Equal(t, localConfig.Redis.Host, "127.0.0.1")
	assert.Equal(t, localConfig.Redis.Port, 6379)

	assert.Equal(t, localConfig.Server.Port, 6666)
}
