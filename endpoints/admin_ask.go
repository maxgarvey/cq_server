package endpoints

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/benbjohnson/clock"
	"github.com/maxgarvey/cq_server/admin"
	"github.com/maxgarvey/cq_server/rabbitmq"
	"github.com/maxgarvey/cq_server/redis"
)

func AdminAsk(
	admin admin.Adminer, clock clock.Clock, rabbitmq rabbitmq.Rabbit, redisClient *redis.Redis, token func() string, logger *slog.Logger,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		valid, err := ValidateAndExtendSession(r, admin)
		if !valid || err != nil {
			logger.Error(
				fmt.Sprintf("invalid token or error accessing: %s\n",
					fmt.Errorf("%w", err),
				),
			)
			return
		}
		PerformAsk(
			w,
			r,
			clock,
			rabbitmq,
			redisClient,
			token,
			logger,
		)
	}
}
