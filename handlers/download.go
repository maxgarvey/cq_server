package handlers

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/maxgarvey/cq_server/data"
)

type DownloadBody struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

// HandleDownload downloads a file from the address in record.Body
func HandleDownload(record *data.Record, logger *slog.Logger) error {
	var body DownloadBody
	if err := json.Unmarshal([]byte(record.Body), &body); err != nil {
		logger.Error("Failed to parse download body", "err", err)
		return err
	}

	logger.Info("DOWNLOAD handler called", "source", body.Source, "destination", body.Destination)
	resp, err := http.Get(body.Source)
	if err != nil {
		logger.Error("Failed to download file", "err", err)
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(body.Destination)
	if err != nil {
		logger.Error("Failed to create file", "err", err)
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		logger.Error("Failed to write file", "err", err)
		return err
	}
	logger.Info("File downloaded successfully", "destination", body.Destination)
	return nil
}
