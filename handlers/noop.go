package handlers

import (
	"log/slog"

	"github.com/maxgarvey/cq_server/data"
)

// HandleNOOP processes a NOOP request.
func HandleNOOP(record *data.Record, logger *slog.Logger) error {
	logger.Info("NOOP handler called", "record", record)
	return nil
}
