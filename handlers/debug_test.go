package handlers

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/maxgarvey/cq_server/data"
)

func TestHandleDEBUG(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))
	record := &data.Record{ID: "debug-test"}
	err := HandleDEBUG(record, logger)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	logOutput := buf.String()
	if logOutput == "" || !bytes.Contains([]byte(logOutput), []byte("DEBUG handler called")) {
		t.Errorf("expected log output to contain 'DEBUG handler called', got: %s", logOutput)
	}
}
