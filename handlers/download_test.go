package handlers

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maxgarvey/cq_server/data"
)

func TestHandleDownload(t *testing.T) {
	// Use httptest server to simulate a file download
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test file content"))
	}))
	defer ts.Close()

	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))
	body := DownloadBody{
		Source:      ts.URL,
		Destination: "testfile.txt",
	}
	bodyBytes, _ := json.Marshal(body)
	record := &data.Record{Body: string(bodyBytes)}
	err := HandleDownload(record, logger)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	logOutput := buf.String()
	if logOutput == "" || !bytes.Contains([]byte(logOutput), []byte("DOWNLOAD handler called")) {
		t.Errorf("expected log output to contain 'DOWNLOAD handler called', got: %s", logOutput)
	}
}
