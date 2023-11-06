package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	localConfig := GetConfig(
		"localhost",
	)

	// Redis
	assert.Equal(
		t,
		localConfig.Redis.ConnectionType,
		"tcp")
	assert.Equal(
		t,
		localConfig.Redis.Host,
		"127.0.0.1",
	)
	assert.Equal(
		t,
		localConfig.Redis.Port,
		6379,
	)

	// RabbitMQ
	assert.Equal(
		t,
		localConfig.Rabbitmq.Host,
		"127.0.0.1",
	)
	assert.Equal(
		t,
		localConfig.Rabbitmq.Port,
		5672,
	)

	// This Server
	assert.Equal(
		t,
		localConfig.Server.Port,
		6666,
	)
}
