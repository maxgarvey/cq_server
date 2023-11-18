package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/maxgarvey/cq_server/config"
	"github.com/maxgarvey/cq_server/data"
)

type Redis struct {
	Client redis.Client
}

// Initializes a new client
func Init(config config.Redis) *Redis {
	// Connect to redis based off of config.
	redisClient := redis.NewClient(
		&redis.Options{
			Addr: fmt.Sprintf(
				"%s:%d",
				config.Host,
				config.Port,
			),
			Password: "", // no password set
			DB:       0,  // use default DB
		},
	)

	return &Redis{
		Client: *redisClient,
	}
}

// Close connection to redis.
func (r Redis) Close() {
	r.Client.Conn().Close()
}

// Retrieve an item from redis datastore.
func (r Redis) Get(ctx context.Context, key string) (data.Record, error) {
	var response data.Record
	err := r.Client.Get(
		ctx,
		key,
	// Parse response.
	).Scan(&response)
	if err != nil {
		return response, err
	}
	return response, nil
}

// Sets a value in the redis datastore.
func (r Redis) Set(
	ctx context.Context,
	key string,
	value []byte,
) error {
	return r.Client.Set(
		ctx,
		key,
		value,
		0,
	).Err()
}
