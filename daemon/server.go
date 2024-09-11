package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (a *App) sendStatsToAPI(fileStats string, systemMonitor string) error {
	payload := map[string]string{
		"file_stats":     fileStats,
		"system_monitor": systemMonitor,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal stats to JSON: %w", err)
	}

	// was not in the requirements to make it a config
	apiEndpoint := a.config.APIEndpoint

	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send stats to API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API responded with status: %s", resp.Status)
	}

	a.logger.Println("Successfully sent stats to API")
	return nil
}

func (a *App) startHTTPServer() {
	http.HandleFunc("/health", a.healthCheckHandler)
	http.HandleFunc("/logs", a.logsHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		a.logger.Fatalf("Failed to start HTTP server: %v", err)
	}
}

func (a *App) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Service is healthy")
}
