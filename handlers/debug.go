package handlers

import (
	"log/slog"

	"github.com/maxgarvey/cq_server/data"
)

// HandleDEBUG processes a DEBUG request.
func HandleDEBUG(record *data.Record, logger *slog.Logger) error {
	logger.Info("DEBUG handler called", "record", record)
	return nil
}
