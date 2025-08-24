package handlers

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/maxgarvey/cq_server/data"
)

func TestHandleNOOP(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))
	record := &data.Record{ID: "noop-test"}
	err := HandleNOOP(record, logger)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	logOutput := buf.String()
	if logOutput == "" || !bytes.Contains([]byte(logOutput), []byte("NOOP handler called")) {
		t.Errorf("expected log output to contain 'NOOP handler called', got: %s", logOutput)
	}
}
